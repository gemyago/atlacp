package bitbucket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

type contextKey string

// GetPRDiffParams contains parameters for getting diff for a pull request.
type GetPRDiffParams struct {
	RepoOwner string
	RepoName  string
	PRID      int
	FilePaths []string // optional
	Context   *int     // optional, default 3
	Account   *string  // optional
}

// GetPRDiff retrieves the diff for a pull request, handling parameters, pagination, and error validation.
func (c *Client) GetPRDiff(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetPRDiffParams,
) (*Diff, error) {
	// Parameter validation
	if params.RepoOwner == "" {
		return nil, errors.New("RepoOwner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("RepoName is required")
	}
	if params.PRID == 0 {
		return nil, errors.New("PRID is required and must be non-zero")
	}

	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build path
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diff",
		url.PathEscape(params.RepoOwner),
		url.PathEscape(params.RepoName),
		params.PRID,
	)

	// Build query parameters
	query := url.Values{}
	if len(params.FilePaths) > 0 {
		for _, fp := range params.FilePaths {
			query.Add("path", fp)
		}
	}
	if params.Context != nil {
		query.Set("context", strconv.Itoa(*params.Context))
	}
	// Bitbucket API does not support account as a query param, but if needed as a header:
	if params.Account != nil {
		// Use a custom header for account if required by internal convention
		const contextKeyAtlassianAccount = contextKey("X-Atlassian-Account")
		ctxWithAuth = context.WithValue(ctxWithAuth, contextKeyAtlassianAccount, *params.Account)
	}

	fullURL := c.baseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var aggregatedDiff []byte
	nextURL := fullURL
	for {
		req, reqErr := http.NewRequestWithContext(ctxWithAuth, http.MethodGet, nextURL, nil)
		if reqErr != nil {
			return nil, fmt.Errorf("failed to create request: %w", reqErr)
		}
		req.Header.Set("Accept", "text/plain")
		if params.Account != nil {
			req.Header.Set("X-Atlassian-Account", *params.Account)
		}

		resp, doErr := c.httpClient.Do(req)
		if doErr != nil {
			return nil, fmt.Errorf("get diff failed: %w", doErr)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("get diff failed: status %d, body: %s", resp.StatusCode, string(body))
		}

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read diff response: %w", readErr)
		}
		aggregatedDiff = append(aggregatedDiff, body...)

		// Check for pagination: Bitbucket raw diff endpoint does not paginate, but if it did, look for a "next" link header
		nextLink := resp.Header.Get("X-Next-Page")
		if nextLink == "" {
			break
		}
		nextURL = nextLink
	}

	diff := Diff(string(aggregatedDiff))
	return &diff, nil
}

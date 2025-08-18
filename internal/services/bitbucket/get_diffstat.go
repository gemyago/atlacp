package bitbucket

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetPRDiffStatParams contains parameters for getting diffstat for a pull request.
type GetPRDiffStatParams struct {
	RepoOwner string
	RepoName  string
	PRID      int
	FilePaths []string // optional
	Context   *int     // optional, default 3
	Account   *string  // optional
}

// GetPRDiffStat retrieves the diffstat for a pull request, handling all query parameters, pagination, and model unification.
func (c *Client) GetPRDiffStat(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetPRDiffStatParams,
) (*struct {
	Size    int        `json:"size,omitempty"`
	Page    int        `json:"page,omitempty"`
	PageLen int        `json:"pagelen,omitempty"`
	Values  []DiffStat `json:"values"`
}, error) {
	// Validate required parameters
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

	// Build base path
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diffstat",
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
	if params.Account != nil {
		query.Set("account_id", *params.Account)
	}

	// Pagination loop
	type paginatedResponse struct {
		Size     int        `json:"size,omitempty"`
		Page     int        `json:"page,omitempty"`
		PageLen  int        `json:"pagelen,omitempty"`
		Next     string     `json:"next,omitempty"`
		Previous string     `json:"previous,omitempty"`
		Values   []DiffStat `json:"values"`
	}

	var (
		allValues []DiffStat
		page      int
		pageLen   int
		totalSize int
		firstPage = true
		nextURL   string
	)

	// Initial URL
	baseURL := c.baseURL + path
	if len(query) > 0 {
		baseURL += "?" + query.Encode()
	}
	nextURL = baseURL

	for {
		var resp paginatedResponse
		err := httpservices.SendRequest(
			ctxWithAuth,
			c.httpClient,
			httpservices.SendRequestParams[interface{}, paginatedResponse]{
				Method: "GET",
				URL:    nextURL,
				Target: &resp,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("get diffstat failed: %w", err)
		}

		if firstPage {
			page = resp.Page
			pageLen = resp.PageLen
			totalSize = resp.Size
			firstPage = false
		}
		allValues = append(allValues, resp.Values...)

		if resp.Next == "" {
			break
		}
		nextURL = resp.Next
	}

	return &struct {
		Size    int        `json:"size,omitempty"`
		Page    int        `json:"page,omitempty"`
		PageLen int        `json:"pagelen,omitempty"`
		Values  []DiffStat `json:"values"`
	}{
		Size:    totalSize,
		Page:    page,
		PageLen: pageLen,
		Values:  allValues,
	}, nil
}

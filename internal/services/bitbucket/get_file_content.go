package bitbucket

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetFileContentParams contains parameters for getting file content from a commit.
type GetFileContentParams struct {
	RepoOwner  string
	RepoName   string
	CommitHash string
	FilePath   string
	Account    *string // optional
}

// FileContent represents the raw file content.

// GetFileContent retrieves the content of a file at a specific commit.
func (c *Client) GetFileContent(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetFileContentParams,
) (*FileContent, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the path
	path := fmt.Sprintf("/repositories/%s/%s/src/%s/%s",
		url.PathEscape(params.RepoOwner),
		url.PathEscape(params.RepoName),
		url.PathEscape(params.CommitHash),
		params.FilePath,
	)

	// Add ?account=... if provided
	fullURL := c.baseURL + path
	if params.Account != nil && *params.Account != "" {
		u, parseErr := url.Parse(fullURL)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse request URL: %w", parseErr)
		}
		q := u.Query()
		q.Set("account", *params.Account)
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	req, err := http.NewRequestWithContext(ctxWithAuth, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get file content failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Include response body in error for better diagnostics
		return nil, fmt.Errorf("get file content failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return &FileContent{
		Path:    params.FilePath,
		Commit:  params.CommitHash,
		Content: string(body),
	}, nil
}

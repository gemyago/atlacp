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
	Username string
	RepoSlug string
	Commit   string
	FilePath string
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

	path := fmt.Sprintf("/repositories/%s/%s/src/%s/%s",
		url.PathEscape(params.Username),
		url.PathEscape(params.RepoSlug),
		url.PathEscape(params.Commit),
		params.FilePath,
	)

	req, err := http.NewRequestWithContext(ctxWithAuth, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get file content failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get file content failed: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content response: %w", err)
	}

	return &FileContent{Content: string(body)}, nil
}

package bitbucket

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetPRDiffParams contains parameters for getting diff for a pull request.
type GetPRDiffParams struct {
	Username      string
	RepoSlug      string
	PullRequestID int
}

// DiffContent represents the raw diff content.
type DiffContent struct {
	Content string
}

// plainTextTarget is a helper for reading plain text HTTP responses.
type plainTextTarget struct {
	Value *string
}

func (t *plainTextTarget) UnmarshalJSON(data []byte) error {
	// Not used, as we read plain text, not JSON.
	return nil
}

// GetPRDiff retrieves the diff for a pull request.
func (c *Client) GetPRDiff(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetPRDiffParams,
) (*DiffContent, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diff",
		url.PathEscape(params.Username),
		url.PathEscape(params.RepoSlug),
		params.PullRequestID,
	)

	req, err := http.NewRequestWithContext(ctxWithAuth, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get diff failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get diff failed: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff response: %w", err)
	}

	return &DiffContent{Content: string(body)}, nil
}

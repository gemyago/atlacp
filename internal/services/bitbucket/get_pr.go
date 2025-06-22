package bitbucket

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetPRParams contains parameters for getting a pull request.
type GetPRParams struct {
	Username      string `json:"-"`
	RepoSlug      string `json:"-"`
	PullRequestID int    `json:"-"`
}

// GetPR retrieves a specific pull request by ID.
// GET /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}.
func (c *Client) GetPR(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetPRParams,
) (*PullRequest, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthToken(ctx, token)

	var pullRequest PullRequest
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", params.Username, params.RepoSlug, params.PullRequestID)
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[interface{}, PullRequest]{
		Method: "GET",
		URL:    c.baseURL + path,
		Target: &pullRequest,
	})
	if err != nil {
		return nil, fmt.Errorf("get pull request failed: %w", err)
	}

	return &pullRequest, nil
}

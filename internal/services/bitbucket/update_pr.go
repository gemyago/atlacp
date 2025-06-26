package bitbucket

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// UpdatePRParams contains parameters for updating a pull request.
type UpdatePRParams struct {
	Username      string       `json:"-"`
	RepoSlug      string       `json:"-"`
	PullRequestID int          `json:"-"`
	Request       *PullRequest `json:"-"`
}

// UpdatePR updates a pull request.
// PUT /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}.
func (c *Client) UpdatePR(
	ctx context.Context,
	tokenProvider TokenProvider,
	params UpdatePRParams,
) (*PullRequest, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	var pullRequest PullRequest
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", params.Username, params.RepoSlug, params.PullRequestID)
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[PullRequest, PullRequest]{
		Method: "PUT",
		URL:    c.baseURL + path,
		Body:   params.Request,
		Target: &pullRequest,
	})
	if err != nil {
		return nil, fmt.Errorf("update pull request failed: %w", err)
	}

	return &pullRequest, nil
}

package bitbucket

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// CreatePRParams contains parameters for creating a pull request.
type CreatePRParams struct {
	Username string       `json:"-"`
	RepoSlug string       `json:"-"`
	Request  *PullRequest `json:"-"`
}

// CreatePR creates a new pull request.
// POST /repositories/{username}/{repo_slug}/pullrequests.
func (c *Client) CreatePR(
	ctx context.Context,
	tokenProvider TokenProvider,
	params CreatePRParams,
) (*PullRequest, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	var pullRequest PullRequest
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests", params.Username, params.RepoSlug)
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[PullRequest, PullRequest]{
		Method: "POST",
		URL:    c.baseURL + path,
		Body:   params.Request,
		Target: &pullRequest,
	})
	if err != nil {
		return nil, fmt.Errorf("create pull request failed: %w", err)
	}

	return &pullRequest, nil
}

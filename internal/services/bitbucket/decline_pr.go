package bitbucket

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// DeclinePRParams contains parameters for declining a pull request.
type DeclinePRParams struct {
	Username      string `json:"-"`
	RepoSlug      string `json:"-"`
	PullRequestID int    `json:"-"`
}

// DeclinePR declines a specific pull request without a request body.
// POST /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/decline.
func (c *Client) DeclinePR(
	ctx context.Context,
	tokenProvider TokenProvider,
	params DeclinePRParams,
) (*PullRequest, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	var pullRequest PullRequest
	path := fmt.Sprintf(
		"/repositories/%s/%s/pullrequests/%d/decline",
		params.Username,
		params.RepoSlug,
		params.PullRequestID,
	)
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[any, PullRequest]{
		Method: "POST",
		URL:    c.baseURL + path,
		Target: &pullRequest,
	})
	if err != nil {
		return nil, fmt.Errorf("decline pull request failed: %w", err)
	}

	return &pullRequest, nil
}

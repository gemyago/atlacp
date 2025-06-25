package bitbucket

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// ApprovePRParams contains parameters for approving a pull request.
type ApprovePRParams struct {
	Username      string `json:"-"`
	RepoSlug      string `json:"-"`
	PullRequestID int    `json:"-"`
}

// ApprovePR approves a specific pull request.
// POST /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/approve.
func (c *Client) ApprovePR(
	ctx context.Context,
	tokenProvider TokenProvider,
	params ApprovePRParams,
) (*Participant, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	var participant Participant
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/approve", params.Username, params.RepoSlug, params.PullRequestID)
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[interface{}, Participant]{
		Method: "POST",
		URL:    c.baseURL + path,
		Target: &participant,
	})
	if err != nil {
		return nil, fmt.Errorf("approve pull request failed: %w", err)
	}

	return &participant, nil
}

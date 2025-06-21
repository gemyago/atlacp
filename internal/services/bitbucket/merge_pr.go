package bitbucket

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// MergePRParams contains parameters for merging a pull request.
type MergePRParams struct {
	Username        string                      `json:"-"`
	RepoSlug        string                      `json:"-"`
	PullRequestID   int                         `json:"-"`
	MergeParameters *PullRequestMergeParameters `json:"-"`
}

// MergePR merges a pull request.
// POST /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/merge.
func (c *Client) MergePR(
	ctx context.Context,
	tokenProvider TokenProvider,
	params MergePRParams,
) (*PullRequest, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthToken(ctx, token)

	var pullRequest PullRequest
	path := fmt.Sprintf(
		"/repositories/%s/%s/pullrequests/%d/merge",
		params.Username, params.RepoSlug, params.PullRequestID,
	)

	// Use SendRequest with merge parameters if provided, otherwise send empty body
	if params.MergeParameters != nil {
		err = httpservices.SendRequest(
			ctxWithAuth,
			c.httpClient,
			httpservices.SendRequestParams[PullRequestMergeParameters, PullRequest]{
				Method: "POST",
				URL:    c.baseURL + path,
				Body:   params.MergeParameters,
				Target: &pullRequest,
			},
		)
	} else {
		err = httpservices.SendRequest(
			ctxWithAuth,
			c.httpClient,
			httpservices.SendRequestParams[interface{}, PullRequest]{
				Method: "POST",
				URL:    c.baseURL + path,
				Target: &pullRequest,
			},
		)
	}

	if err != nil {
		return nil, fmt.Errorf("merge pull request failed: %w", err)
	}

	return &pullRequest, nil
}

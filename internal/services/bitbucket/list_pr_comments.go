package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// ListPRComments retrieves all comments for a specific pull request.
// GET /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments.
func (c *Client) ListPRComments(
	ctx context.Context,
	tokenProvider TokenProvider,
	params ListPRCommentsParams,
) (*ListPRCommentsResponse, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the URL
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PRID,
	)

	var response ListPRCommentsResponse
	err = httpservices.SendRequest(
		ctxWithAuth, c.httpClient,
		httpservices.SendRequestParams[interface{}, ListPRCommentsResponse]{
			Method: "GET",
			URL:    c.baseURL + path,
			Target: &response,
		})
	if err != nil {
		return nil, fmt.Errorf("list pull request comments failed: %w", err)
	}

	return &response, nil
}

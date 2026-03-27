package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// ResolvePRCommentParams identifies a pull request comment to resolve.
type ResolvePRCommentParams struct {
	Workspace string
	RepoSlug  string
	PRID      int64
	CommentID int64
}

// ResolvePRComment resolves a pull request comment thread (no request body).
// POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments/{comment_id}/resolve.
func (c *Client) ResolvePRComment(
	ctx context.Context,
	tokenProvider TokenProvider,
	params ResolvePRCommentParams,
) (*CommentResolution, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments/%d/resolve",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PRID,
		params.CommentID,
	)

	var response CommentResolution
	err = httpservices.SendRequest(
		ctxWithAuth, c.httpClient,
		httpservices.SendRequestParams[interface{}, CommentResolution]{
			Method: "POST",
			URL:    c.baseURL + path,
			Target: &response,
		})
	if err != nil {
		return nil, fmt.Errorf("resolve pull request comment failed: %w", err)
	}

	return &response, nil
}

package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetPRCommentParams identifies a pull request comment.
type GetPRCommentParams struct {
	Workspace string
	RepoSlug  string
	PRID      int64
	CommentID int64
}

// GetPRComment fetches a single pull request comment by ID.
// GET /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments/{comment_id}.
func (c *Client) GetPRComment(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetPRCommentParams,
) (*PRComment, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments/%d",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PRID,
		params.CommentID,
	)

	var response PRComment
	err = httpservices.SendRequest(
		ctxWithAuth, c.httpClient,
		httpservices.SendRequestParams[interface{}, PRComment]{
			Method: "GET",
			URL:    c.baseURL + path,
			Target: &response,
		})
	if err != nil {
		return nil, fmt.Errorf("get pull request comment failed: %w", err)
	}

	return &response, nil
}

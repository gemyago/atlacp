package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetTaskParams contains parameters for getting a specific task on a pull request.
type GetTaskParams struct {
	// Required path parameters
	Workspace string
	RepoSlug  string
	PullReqID int
	TaskID    int
}

// GET /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/tasks/{task_id}.
func (c *Client) GetTask(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetTaskParams,
) (*PullRequestCommentTask, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PullReqID,
		params.TaskID)

	var task PullRequestCommentTask
	err = httpservices.SendRequest(
		ctxWithAuth,
		c.httpClient,
		httpservices.SendRequestParams[interface{}, PullRequestCommentTask]{
			Method: "GET",
			URL:    c.baseURL + path,
			Target: &task,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %w", err)
	}

	return &task, nil
}

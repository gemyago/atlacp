package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// UpdateTaskParams contains parameters for updating a task on a pull request.
type UpdateTaskParams struct {
	// Required path parameters
	Workspace string
	RepoSlug  string
	PullReqID int
	TaskID    int

	// Optional update parameters
	Content string // The updated task content
	State   string // The state of the task ("RESOLVED" or "UNRESOLVED")
}

// UpdateTaskPayload represents the request payload for updating a task.
type UpdateTaskPayload struct {
	Content *TaskContentUpdate `json:"content,omitempty"`
	State   string             `json:"state,omitempty"`
}

// TaskContentUpdate represents the content update for a task.
type TaskContentUpdate struct {
	Raw string `json:"raw"`
}

// PUT /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/tasks/{task_id}.
func (c *Client) UpdateTask(
	ctx context.Context,
	tokenProvider TokenProvider,
	params UpdateTaskParams,
) (*PullRequestCommentTask, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the URL
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PullReqID,
		params.TaskID)

	// Create request payload
	payload := UpdateTaskPayload{}

	if params.Content != "" {
		payload.Content = &TaskContentUpdate{
			Raw: params.Content,
		}
	}

	if params.State != "" {
		payload.State = params.State
	}

	// Make API call
	var response PullRequestCommentTask
	err = http.SendRequest(ctxWithAuth, c.httpClient, http.SendRequestParams[UpdateTaskPayload, PullRequestCommentTask]{
		Method: "PUT",
		URL:    c.baseURL + path,
		Body:   &payload,
		Target: &response,
	})
	if err != nil {
		return nil, fmt.Errorf("update task failed: %w", err)
	}

	return &response, nil
}

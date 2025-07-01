package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// CreatePullRequestTaskParams contains parameters for creating a task on a pull request.
type CreatePullRequestTaskParams struct {
	// Required path parameters
	Workspace string
	RepoSlug  string
	PullReqID int

	// Required body parameters
	Content string // The task content

	// Optional body parameters
	CommentID int64 // Optional comment ID to associate with the task
	Pending   *bool // Optional status of the task (nil = default based on API)
}

// CreateTaskPayload represents the request payload for creating a task.
type CreateTaskPayload struct {
	Content struct {
		Raw string `json:"raw"`
	} `json:"content"`
	Comment *CommentReference `json:"comment,omitempty"`
	Pending *bool             `json:"pending,omitempty"`
}

// CommentReference represents a reference to a comment.
type CommentReference struct {
	ID int64 `json:"id"`
}

// CreatePullRequestTask creates a new task on a pull request.
func (c *Client) CreatePullRequestTask(
	ctx context.Context,
	tokenProvider TokenProvider,
	params CreatePullRequestTaskParams,
) (*PullRequestCommentTask, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the URL
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PullReqID)

	// Create request payload
	payload := CreateTaskPayload{}
	payload.Content.Raw = params.Content

	if params.CommentID > 0 {
		payload.Comment = &CommentReference{
			ID: params.CommentID,
		}
	}

	if params.Pending != nil {
		payload.Pending = params.Pending
	}

	// Make API call
	var response PullRequestCommentTask
	err = http.SendRequest(ctxWithAuth, c.httpClient, http.SendRequestParams[CreateTaskPayload, PullRequestCommentTask]{
		Method: "POST",
		URL:    c.baseURL + path,
		Body:   &payload,
		Target: &response,
	})
	if err != nil {
		return nil, fmt.Errorf("create pull request task failed: %w", err)
	}

	return &response, nil
}

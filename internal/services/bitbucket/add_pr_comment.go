package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

const addPRCommentStatusSuccess = "success"

// AddPRCommentParams contains parameters for adding a comment to a pull request.
type AddPRCommentParams struct {
	Workspace   string // repo_owner
	RepoSlug    string // repo_name
	PullReqID   int
	CommentText string
	FilePath    string // optional, for inline
	LineFrom    int    // optional, for inline
	LineTo      int    // optional, for inline
	Account     string // optional, for future use
	Pending     bool   // optional, for pending comments
}

// addPRCommentPayload matches the Bitbucket API for PR comments.
type addPRCommentPayload struct {
	Content struct {
		Raw string `json:"raw"`
	} `json:"content"`
	Inline *struct {
		Path string `json:"path"`
		From int    `json:"from,omitempty"`
		To   int    `json:"to,omitempty"`
	} `json:"inline,omitempty"`
	Pending bool `json:"pending,omitempty"`
}

// AddPRComment adds a comment to a Bitbucket pull request.
// Returns the comment ID and status string.
func (c *Client) AddPRComment(
	ctx context.Context,
	tokenProvider TokenProvider,
	params AddPRCommentParams,
) (int64, string, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the URL
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PullReqID,
	)

	// Build the payload
	payload := addPRCommentPayload{}
	payload.Content.Raw = params.CommentText
	payload.Pending = params.Pending

	// If file path is provided, treat as inline comment
	if params.FilePath != "" {
		inline := &struct {
			Path string `json:"path"`
			From int    `json:"from,omitempty"`
			To   int    `json:"to,omitempty"`
		}{
			Path: params.FilePath,
		}
		if params.LineFrom > 0 {
			inline.From = params.LineFrom
		}
		if params.LineTo > 0 {
			inline.To = params.LineTo
		}
		payload.Inline = inline
	}

	// Make API call
	var response Comment
	err = http.SendRequest(ctxWithAuth, c.httpClient, http.SendRequestParams[addPRCommentPayload, Comment]{
		Method: "POST",
		URL:    c.baseURL + path,
		Body:   &payload,
		Target: &response,
	})
	if err != nil {
		return 0, "", fmt.Errorf("add pull request comment failed: %w", err)
	}

	return response.ID, addPRCommentStatusSuccess, nil
}

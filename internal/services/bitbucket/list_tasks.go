package bitbucket

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// ListPullRequestTasksParams contains parameters for listing tasks on a pull request.
type ListPullRequestTasksParams struct {
	// Required path parameters
	Workspace string
	RepoSlug  string
	PullReqID int

	// Optional query parameters
	Query   string
	Sort    string
	PageLen int
}

// ListPullRequestTasks returns a paginated list of tasks on a pull request.
func (c *Client) ListPullRequestTasks(
	ctx context.Context,
	tokenProvider TokenProvider,
	params ListPullRequestTasksParams,
) (*PaginatedTasks, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the URL with query parameters
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PullReqID)

	query := url.Values{}
	if params.Query != "" {
		query.Add("q", params.Query)
	}
	if params.Sort != "" {
		query.Add("sort", params.Sort)
	}
	if params.PageLen > 0 {
		query.Add("pagelen", strconv.Itoa(params.PageLen))
	}

	requestURL := c.baseURL + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	// Make API call
	var response PaginatedTasks
	err = http.SendRequest(ctxWithAuth, c.httpClient, http.SendRequestParams[interface{}, PaginatedTasks]{
		Method: "GET",
		URL:    requestURL,
		Target: &response,
	})
	if err != nil {
		return nil, fmt.Errorf("list pull request tasks failed: %w", err)
	}

	return &response, nil
}

package bitbucket

import (
	"context"
	"fmt"
	"net/url"
	"time"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// RequestPRChangesParams contains parameters for requesting changes (removing approval) on a pull request.
type RequestPRChangesParams struct {
	Workspace string // repo_owner
	RepoSlug  string // repo_name
	PullReqID int
	Account   string // optional, for future use
}

// RequestPRChanges removes approval from a specific pull request (requests changes).
// DELETE /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve
// Returns status and timestamp of approval removal.
func (c *Client) RequestPRChanges(
	ctx context.Context,
	tokenProvider TokenProvider,
	params RequestPRChangesParams,
) (status string, removedAt time.Time, err error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the URL
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/approve",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PullReqID,
	)

	// The Bitbucket API returns the removed participant object.
	var participant Participant
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[interface{}, Participant]{
		Method: "DELETE",
		URL:    c.baseURL + path,
		Target: &participant,
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("request changes (remove approval) failed: %w", err)
	}

	// There is no explicit timestamp in the Participant struct, so use current time.
	return "approval_removed", time.Now().UTC(), nil
}

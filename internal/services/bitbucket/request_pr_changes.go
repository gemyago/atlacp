package bitbucket

import (
	"context"
	"fmt"
	"net/url"
	"time"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// RequestPRChangesParams contains parameters for requesting changes on a pull request.
type RequestPRChangesParams struct {
	Workspace string // repo_owner
	RepoSlug  string // repo_name
	PullReqID int
	Account   string // optional, for future use
}

// RequestPRChanges requests changes on a specific pull request.
// POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/request-changes
// Returns status and timestamp of the request.
func (c *Client) RequestPRChanges(
	ctx context.Context,
	tokenProvider TokenProvider,
	params RequestPRChangesParams,
) (string, time.Time, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	// Build the URL
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/request-changes",
		url.PathEscape(params.Workspace),
		url.PathEscape(params.RepoSlug),
		params.PullReqID,
	)

	// The Bitbucket API returns the participant object with the new state.
	var participant Participant
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[interface{}, Participant]{
		Method: "POST",
		URL:    c.baseURL + path,
		Target: &participant,
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("request changes failed: %w", err)
	}

	// Return the participant state and current timestamp
	return participant.State, time.Now().UTC(), nil
}

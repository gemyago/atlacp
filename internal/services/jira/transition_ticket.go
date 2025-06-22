package jira

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// TransitionTicketParams contains parameters for transitioning a Jira ticket.
type TransitionTicketParams struct {
	Domain       string                 `json:"-"` // Jira domain (e.g., "company" in company.atlassian.net)
	TicketKey    string                 `json:"-"` // The ticket key (e.g., "PROJECT-123")
	TransitionID string                 `json:"-"` // The ID of the transition to perform
	Fields       map[string]interface{} `json:"-"` // Optional fields to update during transition
	Update       map[string]interface{} `json:"-"` // Optional updates to perform during transition
}

// TransitionTicket transitions a Jira ticket to a new status.
// POST /rest/api/3/issue/{issueIdOrKey}/transitions.
func (c *Client) TransitionTicket(
	ctx context.Context,
	tokenProvider TokenProvider,
	params TransitionTicketParams,
) error {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthToken(ctx, token)

	baseURL := c.GetBaseURL(params.Domain)
	path := fmt.Sprintf("/issue/%s/transitions", params.TicketKey)

	// Create transition request
	request := TransitionRequest{
		Fields: params.Fields,
		Update: params.Update,
	}
	request.Transition.ID = params.TransitionID

	sendParams := httpservices.SendRequestParams[TransitionRequest, interface{}]{
		Method: "POST",
		URL:    baseURL + path,
		Body:   &request,
		Target: nil, // No response body expected for successful transition
	}
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, sendParams)
	if err != nil {
		return fmt.Errorf("transition ticket failed: %w", err)
	}

	return nil
}

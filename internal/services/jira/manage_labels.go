package jira

import (
	"context"
	"fmt"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// ManageLabelsParams contains parameters for managing labels on a Jira ticket.
type ManageLabelsParams struct {
	Domain       string   `json:"-"` // Jira domain (e.g., "company" in company.atlassian.net)
	TicketKey    string   `json:"-"` // The ticket key (e.g., "PROJECT-123")
	AddLabels    []string `json:"-"` // Labels to add to the ticket
	RemoveLabels []string `json:"-"` // Labels to remove from the ticket
}

// ManageLabels adds and/or removes labels from a Jira ticket.
// PUT /rest/api/3/issue/{issueIdOrKey}.
func (c *Client) ManageLabels(
	ctx context.Context,
	tokenProvider TokenProvider,
	params ManageLabelsParams,
) error {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthToken(ctx, token)

	baseURL := c.GetBaseURL(params.Domain)
	path := fmt.Sprintf("/issue/%s", params.TicketKey)

	// Create label update request
	request := LabelUpdateRequest{}
	request.Update.Labels = make([]LabelOperation, 0, len(params.AddLabels)+len(params.RemoveLabels))

	// Add labels
	for _, label := range params.AddLabels {
		request.Update.Labels = append(request.Update.Labels, LabelOperation{
			Add: label,
		})
	}

	// Remove labels
	for _, label := range params.RemoveLabels {
		request.Update.Labels = append(request.Update.Labels, LabelOperation{
			Remove: label,
		})
	}

	sendParams := httpservices.SendRequestParams[LabelUpdateRequest, interface{}]{
		Method: "PUT",
		URL:    baseURL + path,
		Body:   &request,
		Target: nil, // No response body expected for successful update
	}
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, sendParams)
	if err != nil {
		return fmt.Errorf("manage labels failed: %w", err)
	}

	return nil
}

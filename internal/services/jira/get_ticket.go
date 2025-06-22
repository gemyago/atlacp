package jira

import (
	"context"
	"fmt"
	"net/url"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetTicketParams contains parameters for retrieving a Jira ticket.
type GetTicketParams struct {
	Domain    string   `json:"-"` // Jira domain (e.g., "company" in company.atlassian.net)
	TicketKey string   `json:"-"` // The ticket key (e.g., "PROJECT-123")
	Fields    []string `json:"-"` // Optional fields to include
	Expand    []string `json:"-"` // Optional expansions (e.g., "renderedFields", "transitions")
}

// GetTicket retrieves a Jira ticket by its key.
// GET /rest/api/3/issue/{issueIdOrKey}.
func (c *Client) GetTicket(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetTicketParams,
) (*Ticket, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthToken(ctx, token)

	baseURL := c.GetBaseURL(params.Domain)
	path := fmt.Sprintf("/issue/%s", params.TicketKey)

	// Add query parameters if provided
	requestURL, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := requestURL.Query()
	if len(params.Fields) > 0 {
		query.Set("fields", joinStrings(params.Fields))
	}
	if len(params.Expand) > 0 {
		query.Set("expand", joinStrings(params.Expand))
	}
	requestURL.RawQuery = query.Encode()

	var ticket Ticket
	err = httpservices.SendRequest(ctxWithAuth, c.httpClient, httpservices.SendRequestParams[interface{}, Ticket]{
		Method: "GET",
		URL:    requestURL.String(),
		Target: &ticket,
	})
	if err != nil {
		return nil, fmt.Errorf("get ticket failed: %w", err)
	}

	return &ticket, nil
}

// joinStrings joins a slice of strings with commas.
func joinStrings(items []string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += ","
		}
		result += item
	}
	return result
}

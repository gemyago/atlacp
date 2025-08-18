package bitbucket

import (
	"context"
	"fmt"
	"net/url"

	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// GetPRDiffStatParams contains parameters for getting diffstat for a pull request.
type GetPRDiffStatParams struct {
	RepoOwner string
	RepoName  string
	PRID      int
	FilePaths []string // optional
	Context   *int     // optional, default 3
	Account   *string  // optional
}

// PaginatedDiffStat represents a paginated list of diffstat results.
type PaginatedDiffStat struct {
	Size     int        `json:"size,omitempty"`
	Page     int        `json:"page,omitempty"`
	PageLen  int        `json:"pagelen,omitempty"`
	Next     string     `json:"next,omitempty"`
	Previous string     `json:"previous,omitempty"`
	Values   []DiffStat `json:"values"`
}

// GetPRDiffStat retrieves the diffstat for a pull request.
func (c *Client) GetPRDiffStat(
	ctx context.Context,
	tokenProvider TokenProvider,
	params GetPRDiffStatParams,
) (*PaginatedDiffStat, error) {
	token, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diffstat",
		url.PathEscape(params.RepoOwner),
		url.PathEscape(params.RepoName),
		params.PRID,
	)

	var result PaginatedDiffStat
	err = httpservices.SendRequest(
		ctxWithAuth,
		c.httpClient,
		httpservices.SendRequestParams[interface{}, PaginatedDiffStat]{
			Method: "GET",
			URL:    c.baseURL + path,
			Target: &result,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get diffstat failed: %w", err)
	}

	return &result, nil
}

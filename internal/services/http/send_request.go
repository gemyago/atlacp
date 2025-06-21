package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Use interface{} if no body and or target are needed.
type SendRequestParams[TBody any, TTarget any] struct {
	// HTTP method (GET, POST, PUT, DELETE, etc.)
	Method string

	// Full URL for the request
	URL string

	// Request body to be JSON marshaled (can be nil for GET requests)
	Body *TBody

	// Target to unmarshal response into (can be nil if response not needed)
	Target *TTarget
}

// SendRequest performs an HTTP request with generic body and target types.
// This is a shared function that can be used by all API clients for consistent
// request handling, error processing, and response unmarshaling.
func SendRequest[TBody any, TTarget any](
	ctx context.Context,
	client *http.Client,
	params SendRequestParams[TBody, TTarget],
) error {
	var reqBody bytes.Buffer
	if params.Body != nil {
		if err := json.NewEncoder(&reqBody).Encode(params.Body); err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, params.Method, params.URL, &reqBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if params.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if params.Target != nil {
		if err = json.NewDecoder(resp.Body).Decode(params.Target); err != nil {
			return fmt.Errorf("failed to unmarshal response into target: %w", err)
		}
	}

	return nil
}

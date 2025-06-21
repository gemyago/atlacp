package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendRequest(t *testing.T) {
	t.Run("GET request with response target", func(t *testing.T) {
		// Create test server that returns a JSON response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/users/123", r.URL.Path)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"id":"123","name":"Test User"}`)
		}))
		defer server.Close()

		client := server.Client()
		ctx := t.Context()

		type ResponseData struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		var response ResponseData
		params := SendRequestParams[interface{}, ResponseData]{
			Method: "GET",
			URL:    server.URL + "/users/123",
			Body:   nil,
			Target: &response,
		}

		err := SendRequest(ctx, client, params)

		require.NoError(t, err)
		assert.Equal(t, "123", response.ID)
		assert.Equal(t, "Test User", response.Name)
	})

	t.Run("POST request with body and response", func(t *testing.T) {
		// Create test server that accepts POST and returns response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/users", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id":"user123","name":"John Doe","email":"john@example.com"}`)
		}))
		defer server.Close()

		client := server.Client()
		ctx := t.Context()

		type RequestData struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		type ResponseData struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		requestBody := RequestData{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		var response ResponseData
		params := SendRequestParams[RequestData, ResponseData]{
			Method: "POST",
			URL:    server.URL + "/users",
			Body:   &requestBody,
			Target: &response,
		}

		err := SendRequest(ctx, client, params)

		require.NoError(t, err)
		assert.Equal(t, "user123", response.ID)
		assert.Equal(t, "John Doe", response.Name)
		assert.Equal(t, "john@example.com", response.Email)
	})

	t.Run("DELETE request with no body or response", func(t *testing.T) {
		// Create test server that handles DELETE
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/users/123", r.URL.Path)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client := server.Client()
		ctx := t.Context()

		params := SendRequestParams[interface{}, interface{}]{
			Method: "DELETE",
			URL:    server.URL + "/users/123",
			Body:   nil,
			Target: nil,
		}

		err := SendRequest(ctx, client, params)

		require.NoError(t, err)
	})

	t.Run("handles HTTP error responses", func(t *testing.T) {
		// Create test server that returns 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":"Not Found"}`)
		}))
		defer server.Close()

		client := server.Client()
		client.Transport = NewClientFactory(ClientFactoryDeps{
			RootLogger: diag.RootTestLogger(),
		}).CreateClient().Transport
		ctx := t.Context()

		params := SendRequestParams[interface{}, interface{}]{
			Method: "GET",
			URL:    server.URL + "/nonexistent",
			Body:   nil,
			Target: nil,
		}

		err := SendRequest(ctx, client, params)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("handles invalid URL", func(t *testing.T) {
		client := &http.Client{}
		ctx := t.Context()

		params := SendRequestParams[interface{}, interface{}]{
			Method: "GET",
			URL:    "not-a-valid-url",
			Body:   nil,
			Target: nil,
		}

		err := SendRequest(ctx, client, params)

		require.Error(t, err)
	})

	t.Run("handles request body marshaling error", func(t *testing.T) {
		client := &http.Client{}
		ctx := t.Context()

		// Use a type that can't be marshaled to JSON
		type InvalidBody struct {
			Channel chan int `json:"channel"` // channels can't be marshaled
		}

		invalidBody := InvalidBody{
			Channel: make(chan int),
		}

		params := SendRequestParams[InvalidBody, interface{}]{
			Method: "POST",
			URL:    "http://example.com",
			Body:   &invalidBody,
			Target: nil,
		}

		err := SendRequest(ctx, client, params)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal request body")
	})

	t.Run("handles response unmarshaling error", func(t *testing.T) {
		// Create test server that returns invalid JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"invalid": json}`) // Invalid JSON
		}))
		defer server.Close()

		client := server.Client()
		ctx := t.Context()

		type ResponseData struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		var response ResponseData
		params := SendRequestParams[interface{}, ResponseData]{
			Method: "GET",
			URL:    server.URL + "/invalid-json",
			Body:   nil,
			Target: &response,
		}

		err := SendRequest(ctx, client, params)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal response")
	})
}

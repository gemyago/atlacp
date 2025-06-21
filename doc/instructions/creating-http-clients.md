# Creating HTTP API Clients

## Overview

This document provides comprehensive instructions for creating HTTP API clients from OpenAPI specifications in the atlacp codebase. It provides concrete templates and patterns for implementation that can be used by AI models or humans to generate robust, maintainable API clients.

## Architectural Decisions

### 1. HTTP Client Infrastructure

**Decision**: Use existing `ClientFactory` with middleware composition pattern.

**Implementation Pattern**:
```go
type ServiceClient struct {
    httpClient *http.Client
    baseURL    string
    logger     *slog.Logger
}

type ServiceClientDeps struct {
    dig.In
    
    ClientFactory *httpservices.ClientFactory
    RootLogger    *slog.Logger
    BaseURL       string `name:"config.serviceApi.baseURL"`
}

func NewServiceClient(deps ServiceClientDeps) *ServiceClient {
    return &ServiceClient{
        httpClient: deps.ClientFactory.CreateClient(), // Uses all middleware by default
        baseURL:    deps.BaseURL,
        logger:     deps.RootLogger.WithGroup("service-client"),
    }
}
```

### 2. Authentication Strategy

**Decision**: Use context-based authentication via existing middleware.

**Implementation Pattern**:
```go
// In the client method
func (c *ServiceClient) CreateResource(ctx context.Context, token string, req *CreateResourceRequest) (*Resource, error) {
    // Add token to context - middleware will handle Bearer header
    ctxWithAuth := middleware.WithAuthToken(ctx, token)
    
    // Make HTTP request with authenticated context
    httpReq, err := http.NewRequestWithContext(ctxWithAuth, "POST", url, body)
    // ... rest of implementation
}
```

## Implementation Templates

### HTTP Client Template

```go
package services

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    "net/url"
    "strings"
    
    httpservices "github.com/gemyago/atlacp/internal/services/http"
    "github.com/gemyago/atlacp/internal/services/http/middleware"
    "go.uber.org/dig"
)

// ServiceClient provides HTTP client functionality for [Service Name] API.
type ServiceClient struct {
    httpClient *http.Client
    baseURL    string
    logger     *slog.Logger
}

// ServiceClientDeps contains dependencies for the service client.
type ServiceClientDeps struct {
    dig.In
    
    ClientFactory *httpservices.ClientFactory
    RootLogger    *slog.Logger
    BaseURL       string `name:"config.serviceApi.baseURL"`
}

// NewServiceClient creates a new service client instance.
func NewServiceClient(deps ServiceClientDeps) *ServiceClient {
    return &ServiceClient{
        httpClient: deps.ClientFactory.CreateClient(),
        baseURL:    strings.TrimSuffix(deps.BaseURL, "/"),
        logger:     deps.RootLogger.WithGroup("service-client"),
    }
}

// doRequest performs an HTTP request using the shared SendRequest function.
func (c *ServiceClient) doRequest(ctx context.Context, method, path string, body interface{}, target interface{}) error {
    // Build full URL
    fullURL := c.baseURL + path
    
    // Use the shared SendRequest function
    params := httpservices.SendRequestParams[interface{}, interface{}]{
        Method: method,
        URL:    fullURL,
        Body:   nil,
        Target: nil,
    }
    
    // Set body if provided
    if body != nil {
        params.Body = &body
    }
    
    // Set target if provided
    if target != nil {
        params.Target = &target
    }
    
    err := httpservices.SendRequest(ctx, c.httpClient, params)
    if err != nil {
        // Convert to APIError if it's not already
        if apiErr, ok := err.(*APIError); ok {
            return apiErr
        }
        
        // Handle other errors by creating APIError
        return &APIError{
            Message:     err.Error(),
            Endpoint:    fullURL,
            HTTPMethod:  method,
            OriginalErr: err,
        }
    }
    
    return nil
}

// handleErrorResponse processes API error responses.
func (c *ServiceClient) handleErrorResponse(resp *http.Response, body []byte, method, url string) error {
    // Try to parse error response as JSON
    var errorResp struct {
        Error   string            `json:"error"`
        Message string            `json:"message"`
        Details map[string]string `json:"details"`
    }
    
    message := fmt.Sprintf("HTTP %d", resp.StatusCode)
    errorCode := ""
    var details map[string]string
    
    if len(body) > 0 {
        if err := json.Unmarshal(body, &errorResp); err == nil {
            if errorResp.Message != "" {
                message = errorResp.Message
            }
            if errorResp.Error != "" {
                errorCode = errorResp.Error
            }
            details = errorResp.Details
        } else {
            // Use raw body as message if JSON parsing fails
            message = string(body)
        }
    }
    
    return &APIError{
        StatusCode:  resp.StatusCode,
        ErrorCode:   errorCode,
        Message:     message,
        Details:     details,
        Endpoint:    url,
        HTTPMethod:  method,
        OriginalErr: nil,
    }
}
```

### API Method Implementation Template

```go
// CreateResource creates a new resource via API.
func (c *ServiceClient) CreateResource(ctx context.Context, token string, req *CreateResourceRequest) (*Resource, error) {
    // Add authentication token to context
    ctxWithAuth := middleware.WithAuthToken(ctx, token)
    
    // Make API call
    var resource Resource
    err := c.doRequest(ctxWithAuth, "POST", "/resources", req, &resource)
    if err != nil {
        return nil, fmt.Errorf("create resource failed: %w", err)
    }
    
    return &resource, nil
}

// GetResource retrieves a resource by ID.
func (c *ServiceClient) GetResource(ctx context.Context, token string, resourceID string) (*Resource, error) {
    ctxWithAuth := middleware.WithAuthToken(ctx, token)
    
    var resource Resource
    path := fmt.Sprintf("/resources/%s", resourceID)
    err := c.doRequest(ctxWithAuth, "GET", path, nil, &resource)
    if err != nil {
        return nil, fmt.Errorf("get resource failed: %w", err)
    }
    
    return &resource, nil
}

// UpdateResource updates an existing resource.
func (c *ServiceClient) UpdateResource(ctx context.Context, token string, resourceID string, req *UpdateResourceRequest) (*Resource, error) {
    ctxWithAuth := middleware.WithAuthToken(ctx, token)
    
    var resource Resource
    path := fmt.Sprintf("/resources/%s", resourceID)
    err := c.doRequest(ctxWithAuth, "PUT", path, req, &resource)
    if err != nil {
        return nil, fmt.Errorf("update resource failed: %w", err)
    }
    
    return &resource, nil
}
```

### Request/Response Model Templates

```go
// Request models
type CreateResourceRequest struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Tags        []string `json:"tags,omitempty"`
}

type UpdateResourceRequest struct {
    Title       string   `json:"title,omitempty"`
    Description string   `json:"description,omitempty"`
    Tags        []string `json:"tags,omitempty"`
}

// Response models  
type Resource struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Tags        []string  `json:"tags"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Status      string    `json:"status"`
}

// List response with pagination
type ListResourcesResponse struct {
    Resources []Resource `json:"resources"`
    Page      int        `json:"page"`
    PageSize  int        `json:"page_size"`
    Total     int        `json:"total"`
    HasMore   bool       `json:"has_more"`
}
```

## Configuration Integration

### Extending Configuration Schema

Add Atlassian-specific configuration in `internal/config/`:

```go
// In your config struct
type Config struct {
    // ... existing fields
    
    Atlassian AtlassianConfig `mapstructure:"atlassian"`
}

type AtlassianConfig struct {
    AccountsFilePath  string `mapstructure:"accountsFilePath"`
    BitbucketBaseURL  string `mapstructure:"bitbucketBaseURL"`
    JiraBaseURL       string `mapstructure:"jiraBaseURL"`
}
```

### Default Configuration Values

```json
{
  "atlassian": {
    "accountsFilePath": "./accounts.json",
    "bitbucketBaseURL": "https://api.bitbucket.org/2.0",
    "jiraBaseURL": "https://{domain}.atlassian.net/rest/api/3"
  }
}
```

## Testing Patterns

### Unit Test Template

```go
package services

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gemyago/atlacp/internal/diag"
    httpservices "github.com/gemyago/atlacp/internal/services/http"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestServiceClient(t *testing.T) {
    makeMockDeps := func(baseURL string) ServiceClientDeps {
        return ServiceClientDeps{
            ClientFactory: httpservices.NewClientFactory(httpservices.ClientFactoryDeps{
                RootLogger: diag.RootTestLogger(),
            }),
            RootLogger: diag.RootTestLogger(),
            BaseURL:    baseURL,
        }
    }
    
    t.Run("CreateResource", func(t *testing.T) {
        t.Run("successful creation", func(t *testing.T) {
            // Create test server that simulates successful API response
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Verify request details
                assert.Equal(t, "POST", r.Method)
                assert.Equal(t, "/resources", r.URL.Path)
                assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
                assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
                
                // Return successful response
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusCreated)
                fmt.Fprint(w, `{
                    "id": "resource-123",
                    "name": "test-resource",
                    "description": "Test description",
                    "status": "active",
                    "created_at": "2023-01-01T00:00:00Z",
                    "updated_at": "2023-01-01T00:00:00Z"
                }`)
            }))
            defer server.Close()
            
            deps := makeMockDeps(server.URL)
            client := NewServiceClient(deps)
            
            req := &CreateResourceRequest{
                Name:        "test-resource",
                Description: "Test description",
            }
            
            resource, err := client.CreateResource(t.Context(), "test-token", req)
            
            require.NoError(t, err)
            assert.Equal(t, "resource-123", resource.ID)
            assert.Equal(t, "test-resource", resource.Name)
            assert.Equal(t, "Test description", resource.Description)
            assert.Equal(t, "active", resource.Status)
        })
        
        t.Run("handles API error", func(t *testing.T) {
            // Create test server that simulates API error
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(w, `{
                    "error": "INVALID_REQUEST",
                    "message": "Name is required",
                    "details": {
                        "field": "name"
                    }
                }`)
            }))
            defer server.Close()
            
            deps := makeMockDeps(server.URL)
            client := NewServiceClient(deps)
            
            req := &CreateResourceRequest{
                Description: "Missing name",
            }
            
            _, err := client.CreateResource(t.Context(), "test-token", req)
            
            require.Error(t, err)
            
            // Verify error details
            apiErr, ok := err.(*APIError)
            require.True(t, ok, "Expected APIError type")
            assert.Equal(t, 400, apiErr.StatusCode)
            assert.Equal(t, "INVALID_REQUEST", apiErr.ErrorCode)
            assert.Equal(t, "Name is required", apiErr.Message)
            assert.Equal(t, "name", apiErr.Details["field"])
        })
        
        t.Run("handles authentication error", func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
                w.WriteHeader(http.StatusUnauthorized)
                fmt.Fprint(w, `{"error": "UNAUTHORIZED", "message": "Invalid token"}`)
            }))
            defer server.Close()
            
            deps := makeMockDeps(server.URL)
            client := NewServiceClient(deps)
            
            req := &CreateResourceRequest{
                Name: "test-resource",
            }
            
            _, err := client.CreateResource(t.Context(), "invalid-token", req)
            
            require.Error(t, err)
            apiErr, ok := err.(*APIError)
            require.True(t, ok)
            assert.Equal(t, 401, apiErr.StatusCode)
        })
    })
    
    t.Run("GetResource", func(t *testing.T) {
        t.Run("successful retrieval", func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                assert.Equal(t, "GET", r.Method)
                assert.Equal(t, "/resources/resource-123", r.URL.Path)
                
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                fmt.Fprint(w, `{
                    "id": "resource-123",
                    "name": "existing-resource",
                    "status": "active"
                }`)
            }))
            defer server.Close()
            
            deps := makeMockDeps(server.URL)
            client := NewServiceClient(deps)
            
            resource, err := client.GetResource(t.Context(), "test-token", "resource-123")
            
            require.NoError(t, err)
            assert.Equal(t, "resource-123", resource.ID)
            assert.Equal(t, "existing-resource", resource.Name)
        })
        
        t.Run("handles not found", func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, `{"error": "NOT_FOUND", "message": "Resource not found"}`)
            }))
            defer server.Close()
            
            deps := makeMockDeps(server.URL)
            client := NewServiceClient(deps)
            
            _, err := client.GetResource(t.Context(), "test-token", "nonexistent")
            
            require.Error(t, err)
            apiErr, ok := err.(*APIError)
            require.True(t, ok)
            assert.Equal(t, 404, apiErr.StatusCode)
        })
    })
}

### Integration Test Template

```go
func TestServiceClientIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests in short mode")
    }
    
    t.Run("full CRUD operations", func(t *testing.T) {
        // Create test server that handles multiple endpoints
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            switch {
            case r.Method == "POST" && r.URL.Path == "/resources":
                // Handle CREATE
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusCreated)
                fmt.Fprint(w, `{"id": "test-id", "name": "test-resource", "status": "active"}`)
                
            case r.Method == "GET" && r.URL.Path == "/resources/test-id":
                // Handle READ
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                fmt.Fprint(w, `{"id": "test-id", "name": "test-resource", "status": "active"}`)
                
            case r.Method == "PUT" && r.URL.Path == "/resources/test-id":
                // Handle UPDATE
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                fmt.Fprint(w, `{"id": "test-id", "name": "updated-resource", "status": "active"}`)
                
            case r.Method == "DELETE" && r.URL.Path == "/resources/test-id":
                // Handle DELETE
                w.WriteHeader(http.StatusNoContent)
                
            default:
                w.WriteHeader(http.StatusNotFound)
            }
        }))
        defer server.Close()
        
        deps := makeMockDeps(server.URL)
        client := NewServiceClient(deps)
        ctx := t.Context()
        token := "test-token"
        
        // CREATE
        createReq := &CreateResourceRequest{
            Name:        "test-resource",
            Description: "Test description",
        }
        resource, err := client.CreateResource(ctx, token, createReq)
        require.NoError(t, err)
        require.NotEmpty(t, resource.ID)
        
        // READ
        retrieved, err := client.GetResource(ctx, token, resource.ID)
        require.NoError(t, err)
        assert.Equal(t, resource.ID, retrieved.ID)
        
        // UPDATE
        updateReq := &UpdateResourceRequest{
            Title: "updated-resource",
        }
        updated, err := client.UpdateResource(ctx, token, resource.ID, updateReq)
        require.NoError(t, err)
        assert.Equal(t, "updated-resource", updated.Name)
        
        // DELETE would be tested here if the method exists
    })
}
```

## Quality Assurance Guidelines

### 1. Error Handling Requirements
- All API errors must be wrapped in `APIError` with context
- Include HTTP status code, endpoint, and method in errors
- Log all errors with appropriate log levels
- Provide actionable error messages for users

### 2. Logging Requirements
- Use structured logging with context
- Log request/response details at appropriate levels
- Include relevant identifiers (account name, resource IDs)
- Use consistent log groups for filtering

### 3. Testing Requirements
- Write unit tests for all public methods
- Use table-driven tests for multiple scenarios
- Mock external dependencies (HTTP clients, repositories)
- Test both success and error paths
- Aim for >90% code coverage

### 4. Documentation Requirements
- Document all public types and methods
- Include usage examples in Go doc comments
- Document error conditions and return values
- Keep documentation up to date with code changes

### 5. Security Requirements
- Never log authentication tokens or sensitive data
- Use context for passing authentication tokens
- Validate all input parameters
- Handle rate limiting gracefully

## Open Questions

### Token Provider Pattern
Currently using `token string` parameter in methods. Consider `TokenProvider` interface:

```go
type TokenProvider interface {
    GetToken(ctx context.Context) (string, error)
}
```

**Pros**: Better abstraction, can handle token refresh logic
**Cons**: Additional complexity, most use cases just have static tokens

## Summary

This framework provides:
- **Consistency**: Follows existing HTTP middleware patterns
- **Testability**: Interface-based design enables easy mocking
- **Simplicity**: Minimal boilerplate with shared utilities
- **Integration**: Seamless fit with existing codebase architecture
- **Security**: Context-based authentication via middleware
- **Reliability**: Error handling via existing middleware

All API clients built following these patterns will integrate seamlessly with the existing infrastructure while remaining simple and maintainable. 
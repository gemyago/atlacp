# Creating HTTP API Clients

## Overview

This document provides comprehensive instructions for creating HTTP API clients from OpenAPI specifications. It provides concrete templates and patterns for implementation that can be used by AI models or humans to generate robust, maintainable API clients.

## Architectural Decisions

Key principles:
- Maintain separate file per operation to simplify updates and testing. Example `create_resource.go`.
- Maintain separate tests per operation. Example `create_resource_test.go`.
- Maintain separate file per model. Keep models in the same package as the client. Example `model_create_request.go` and `model_create_response.go`.
- Keep common client code in `client.go` file.
- Use simple naming: `Client` instead of `ServiceClient` since it will be accessed as `packagename.Client`.
- Always use consistent operation signature: `ctx`, `tokenProvider`, and `params` struct (even for single parameters).

### 1. HTTP Client Infrastructure

**Decision**: Use existing `ClientFactory` with middleware composition pattern.

**Implementation Pattern**:
```go
type Client struct {
    httpClient *http.Client
    baseURL    string
    logger     *slog.Logger
}

type ClientDeps struct {
    dig.In
    
    ClientFactory *http.ClientFactory
    RootLogger    *slog.Logger
    BaseURL       string `name:"config.serviceApi.baseURL"`
}

func NewClient(deps ClientDeps) *Client {
    return &Client{
        httpClient: deps.ClientFactory.CreateClient(), // Uses all middleware by default
        baseURL:    deps.BaseURL,
        logger:     deps.RootLogger.WithGroup("service-client"),
    }
}
```

### 2. Authentication Strategy

Use context-based authentication via existing middleware. Use token provider interface to get the token.

**Implementation Pattern**:
```go
type TokenProvider interface {
    GetToken(ctx context.Context) (middleware.Token, error)
}

// In the client method - always use params struct even for single parameters.
type CreateResourceParams struct {
    Request *CreateResourceRequest
}

func (c *Client) CreateResource(ctx context.Context, tokenProvider TokenProvider, params CreateResourceParams) (*Resource, error) {
    token, err := tokenProvider.GetToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get token: %w", err)
    }
    ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)
    // ... rest of implementation
}
```

## Implementation Templates

### API Method Implementation Template

```go
// CreateResourceParams contains parameters for creating a resource.
type CreateResourceParams struct {
    Request *CreateResourceRequest
}

// CreateResource is example to show how to send a request with body and response.
func (c *Client) CreateResource(ctx context.Context, tokenProvider TokenProvider, params CreateResourceParams) (*Resource, error) {
    token, err := tokenProvider.GetToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get token: %w", err)
    }
    ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)
    
    // Make API call
    var resource Resource
    err = http.SendRequest(ctxWithAuth, c.httpClient, http.SendRequestParams[CreateResourceRequest, Resource]{
        Method: "POST",
        URL:    c.baseURL + "/resources",
        Body:   params.Request,
        Target: &resource,
    })
    if err != nil {
        return nil, fmt.Errorf("create resource failed: %w", err)
    }
    
    return &resource, nil
}

// GetResourceParams contains parameters for getting a resource.
type GetResourceParams struct {
    ResourceID string
}

// GetResource is example to show how to send a request with no body and response.  
func (c *Client) GetResource(ctx context.Context, tokenProvider TokenProvider, params GetResourceParams) (*Resource, error) {
    token, err := tokenProvider.GetToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get token: %w", err)
    }
    ctxWithAuth := middleware.WithAuthTokenV2(ctx, token)
    
    var resource Resource
    path := fmt.Sprintf("/resources/%s", params.ResourceID)
    err = http.SendRequest(ctxWithAuth, c.httpClient, http.SendRequestParams[interface{}, Resource]{
        Method: "GET",
        URL:    c.baseURL + path,
        Target: &resource,
    })
    if err != nil {
        return nil, fmt.Errorf("get resource failed: %w", err)
    }
    
    return &resource, nil
}
```

### Request/Response Model Templates

```go
// Request models.
type CreateResourceRequest struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Tags        []string `json:"tags,omitempty"`
}

// Response models.
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
```

## Testing Patterns

Follow [testing-best-practices](../testing-best-practices.md) when writing tests.
Always include these 4 test cases for each operation:

1. **Success with all parameters/fields** - Test with complete request and response
2. **Success with required parameters only** - Test minimal valid case  
3. **Generic API error test** - Test API error handling
4. **Generic token provider error test** - Test authentication error

```go
package packagename

import (
    "context"
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gemyago/atlacp/internal/diag"
    httpservices "github.com/gemyago/atlacp/internal/services/http"
    "github.com/go-faker/faker/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestClient_CreateResource(t *testing.T) {
    makeMockDeps := func(baseURL string) ClientDeps {
        return ClientDeps{
            ClientFactory: httpservices.NewClientFactory(httpservices.ClientFactoryDeps{
                RootLogger: diag.RootTestLogger(),
            }),
            RootLogger: diag.RootTestLogger(),
            BaseURL:    baseURL,
        }
    }
    
    mockTokenProvider := &MockTokenProvider{}
    
    t.Run("success with all parameters and fields", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Verify request details
            assert.Equal(t, "POST", r.Method)
            assert.Equal(t, "/resources", r.URL.Path)
            assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
            assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
            
            // Return complete successful response
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusCreated)
            fmt.Fprint(w, `{
                "id": "resource-123",
                "name": "test-resource",
                "title": "Test Resource",
                "description": "Test description",
                "tags": ["tag1", "tag2"],
                "status": "active",
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }`)
        }))
        defer server.Close()
        
        deps := makeMockDeps(server.URL)
        client := NewClient(deps)
        
        req := &CreateResourceRequest{
            Name:        faker.Name(),
            Description: faker.Sentence(),
            Tags:        []string{faker.Word(), faker.Word()},
        }
        
        mockTokenProvider.token = "test-token"
        resource, err := client.CreateResource(t.Context(), mockTokenProvider, CreateResourceParams{
            Request: req,
        })
        
        require.NoError(t, err)
        assert.Equal(t, "resource-123", resource.ID)
        assert.Equal(t, "test-resource", resource.Name)
        assert.Equal(t, "Test Resource", resource.Title)
        assert.Equal(t, "Test description", resource.Description)
        assert.Equal(t, []string{"tag1", "tag2"}, resource.Tags)
        assert.Equal(t, "active", resource.Status)
        assert.NotZero(t, resource.CreatedAt)
        assert.NotZero(t, resource.UpdatedAt)
    })
    
    t.Run("success with required parameters only", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Return minimal successful response
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusCreated)
            fmt.Fprint(w, `{
                "id": "resource-456",
                "name": "minimal-resource"
            }`)
        }))
        defer server.Close()
        
        deps := makeMockDeps(server.URL)
        client := NewClient(deps)
        
        req := &CreateResourceRequest{
            Name: faker.Name(), // Only required field
        }
        
        mockTokenProvider.token = "test-token"
        resource, err := client.CreateResource(t.Context(), mockTokenProvider, CreateResourceParams{
            Request: req,
        })
        
        require.NoError(t, err)
        assert.Equal(t, "resource-456", resource.ID)
        assert.Equal(t, "minimal-resource", resource.Name)
    })
    
    t.Run("handles API error", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
        }))
        defer server.Close()
        
        deps := makeMockDeps(server.URL)
        client := NewClient(deps)
        
        req := &CreateResourceRequest{
            Description: "Missing name",
        }
        
        mockTokenProvider.token = "test-token"
        _, err := client.CreateResource(t.Context(), mockTokenProvider, CreateResourceParams{
            Request: req,
        })
        
        require.Error(t, err)
        assert.Contains(t, err.Error(), "create resource failed")
    })
    
    t.Run("handles token provider error", func(t *testing.T) {
        deps := makeMockDeps("http://example.com")
        client := NewClient(deps)
        
        mockTokenProvider.err = errors.New("token error")
        _, err := client.CreateResource(t.Context(), mockTokenProvider, CreateResourceParams{
            Request: &CreateResourceRequest{
                Name: faker.Name(),
            },
        })
        
        require.Error(t, err)
        assert.Contains(t, err.Error(), "failed to get token")
    })
}

// MockTokenProvider is a simple mock implementation for testing.
type MockTokenProvider struct {
    token string
    err   error
}

func (m *MockTokenProvider) GetToken(ctx context.Context) (middleware.Token, error) {
    if m.err != nil {
        return middleware.Token{}, m.err
    }
    return middleware.Token{Type: "Bearer", Value: m.token}, nil
}
```

## Quality Assurance Guidelines

### 1. Testing Requirements
- Always include the 4 standard test cases per operation
- Use faker for generating random test data
- Follow testing-best-practices patterns

### 2. Documentation Requirements
- Document all public types and methods
- Include usage examples in Go doc comments

### 3. Security Requirements
- Never log authentication tokens or sensitive data
- Use context for passing authentication tokens

### 4. Code Quality and Linting

**IMPORTANT**: After completing implementation, always run linting to ensure code quality:

```bash
make lint
```

This runs `golangci-lint` across the entire codebase and will catch common issues.

#### Common Linting Issues and How to Fix Them

1. **Unused Parameters in HTTP Handlers**
   ```go
   // ❌ Bad: unused parameters will trigger linter warnings
   server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       // Not using 'r' parameter triggers warning
       w.WriteHeader(http.StatusOK)
   }))
   
   // ✅ Good: use underscore for unused parameters
   server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
       w.WriteHeader(http.StatusOK)
   }))
   ```

2. **Naming Conventions (var-naming)**
   ```go
   // ❌ Bad: incorrect casing for acronyms
   type ShowPetByIdParams struct {
       PetID string
   }
   
   // ✅ Good: acronyms should be all uppercase
   type ShowPetByIDParams struct {
       PetID string
   }
   ```

3. **Context Usage in Tests (usetesting)**
   ```go
   // ❌ Bad: using context.Background() in tests
   result, err := client.CreatePets(context.Background(), tokenProvider, params)
   
   // ✅ Good: use t.Context() in tests for better test lifecycle management
   result, err := client.CreatePets(t.Context(), tokenProvider, params)
   ```

4. **Require vs Assert in HTTP Handlers (testifylint)**
   ```go
   // ❌ Bad: using require.NoError in HTTP handlers can cause issues
   func(w http.ResponseWriter, r *http.Request) {
       body, err := io.ReadAll(r.Body)
       require.NoError(t, err) // This can cause issues in handlers
   }
   
   // ✅ Good: use assert.NoError in HTTP handlers
   func(w http.ResponseWriter, r *http.Request) {
       body, err := io.ReadAll(r.Body)
       assert.NoError(t, err) // Safe to use in handlers
   }
   ```

5. **Unused Parameters in Mock Implementations (revive)**
   ```go
   // ❌ Bad: unused context parameter
   func (m *MockTokenProvider) GetToken(ctx context.Context) (string, error) {
       return m.token, m.err
   }
   
   // ✅ Good: mark unused parameters with underscore
   func (m *MockTokenProvider) GetToken(_ context.Context) (string, error) {
       return m.token, m.err
   }
   ```

#### Linting Best Practices

- **Run linting early and often** - Don't wait until the end to check
- **Fix issues immediately** - Address linting warnings as soon as they appear
- **Understand the rules** - Learn what each linter rule is trying to prevent
- **Use meaningful variable names** - Avoid generic names that might trigger warnings
- **Follow Go naming conventions** - Use proper casing for types, methods, and variables

#### Integration with Development Workflow

1. **After Implementation**: Run `make lint` to catch any issues
2. **Before Committing**: Ensure linting passes cleanly
3. **CI/CD Integration**: Linting should be part of the build pipeline
4. **Code Reviews**: Check that new code follows linting standards

This step ensures consistency across the codebase and helps maintain high code quality standards.
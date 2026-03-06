# Creating HTTP API Clients

## Overview

You must strickly follow this instruction to create (or update) HTTP API clients from OpenAPI specifications. Do not use any third-party code generation tools, just follow templates from this document.

## Input and Output

The user will provide the following:
- Path to openapi spec in JSON or YAML format.
- Optionally: List of operations to implement. If not provided - all operations.
- Optionally: Output package path. If not provided, use same folder as input spec.

Your task is to create http client using the guidelines from this document..

## Architectural Decisions

Key principles to follow:
- Keep common client code in `client.go` file.
  - Use simple naming: `Client` instead of `ServiceClient`
- Separate file per operation. Example:
  - Operation: `addPet`
  - File: `add_pet.go`
  - Tests file: `add_pet_test.go`
- Separate file per model. Example:
  - Model: `PetDetails`
  - File: `model_pet_details.go`
  - Models are not tested directly.
- Always use consistent operation signature: `ctx`, and `params` struct (even for single parameters).

## Implementation Templates

### HTTP Client Infrastructure

Key points:
- Use existing `ClientFactory` with middleware composition pattern.
- No need to test `Client` struct directly, test only API methods.

**Implementation Pattern**:

This is `client.go` file content pattern:
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
		BaseURL       string `name:"config.serviceName.baseURL"`
}

func NewClient(
	deps ClientDeps,
	clientOpts ...httpclient.ClientOption,
) *Client {
	return &Client{
		httpClient: deps.ClientFactory.CreateClient(clientOpts...),
		baseURL:    deps.BaseURL,
		logger:     deps.RootLogger.WithGroup("petstore-client"),
	}
}
```

### API Method Implementation Template

Key decisions:
- Each operation in separate file.
- Input params struct is placed in the same file as the operation.

#### Example 1

Example to send POST/PUT/PATCH requests (with http body):
- Operation `addPet`
- File `add_pet.go`
- AddPetParams - input params, declared in the same file

```go
// AddPetParams contains parameters for adding a pet.
type AddPetParams struct {
		// Request represents the request body for adding a pet.
		Request *AddPetRequest
}

// AddPet adds a new pet to the store.
func (c *Client) AddPet(ctx context.Context, params AddPetParams) (*AddPetResponse, error) {
		var response AddPetResponse
		err = http.SendRequest(ctx, c.httpClient, http.SendRequestParams[AddPetRequest, AddPetResponse]{
				Method: "POST",
				URL:    c.baseURL + "/pets",
				Body:   params.Request,
				Target: &response,
		})
		if err != nil {
				return nil, fmt.Errorf("failed to add pet: %w", err)
		}

		return &response, nil
}
```

#### Example 2

Example to send GET requests (no boty) with query/path parameters:
- Operation `getPetById`
- File `get_pet_by_id.go`
- GetPetByIdParams - input params, declared in the same file

```go
// GetPetByIdParams contains parameters for getting a resource.
type GetPetByIdParams struct {
		PetID string
}

// GetPetById is example to show how to send a request with no body and response.
func (c *Client) GetPetById(ctx context.Context, params GetPetByIdParams) (*GetPetByIDResponse, error) {
		var response GetPetByIDResponse
		path := fmt.Sprintf("/pets/%s", params.PetID)
		err = http.SendRequest(ctx, c.httpClient, http.SendRequestParams[interface{}, GetPetByIDResponse]{
				Method: "GET",
				URL:    c.baseURL + path,
				Target: &response,
		})
		if err != nil {
				return nil, fmt.Errorf("get pet by id failed: %w", err)
		}

		return &response, nil
}
```

### Request/Response Model Templates


File name: `model_create_resource_request.go`
```go
// Request models.
type CreateResourceRequest struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tags        []string `json:"tags,omitempty"`
}
```

File name: `model_resource.go`
```go
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

Follow [testing-best-practices](./doc/testing-best-practices.md) when writing tests.
Always include these 3 test cases for each operation:

1. **Success with all parameters/fields** - Test with complete request and response
2. **Success with required parameters only** - Test minimal valid case
3. **Generic API error test** - Test API error handling

### Test Structure Best Practices

1. **Use AAA Pattern**: Structure tests with clear Arrange, Act, Assert sections and add comments to indicate each section
2. **Use Test-Specific Logger**: Include test name in the logger for better debugging
3. **Use Randomized Test Data**: Use faker to generate random test inputs
4. **Use Proper Error Assertions**: Use assert.ErrorContains or assert.ErrorIs for error checking

Here's a test template that incorporates these practices:

```go
package packagename

import (
		"context"
		"errors"
		"fmt"
		"math/rand/v2"
		"net/http"
		"net/http/httptest"
		"testing"

		"github.com/gemyago/atlacp/internal/telemetry"
		"github.com/gemyago/golang-backend-boilerplate/internal/infrastructure/httpclient"
		"github.com/jaswdr/faker/v2"
		"github.com/stretchr/testify/assert"
		"github.com/stretchr/testify/require"
		"golang.org/x/oauth2"
)

func TestClient_CreateResource(t *testing.T) {
		makeMockDeps := func(t *testing.T, baseURL string) ClientDeps {
				// Always include test name in the logger for better debugging
				rootLogger := telemetry.RootTestLogger().With("test", t.Name())
				return ClientDeps{
						ClientFactory: httpclient.NewClientFactory(httpclient.ClientFactoryDeps{
								RootLogger: rootLogger,
						}),
						RootLogger: rootLogger,
						BaseURL:    baseURL,
				}
		}

		fake := faker.New()

		t.Run("success with all parameters and fields", func(t *testing.T) {
				// Arrange - Use randomized data
				resourceName := "resource-" + fake.Lorem().Word()
				resourceDesc := fake.Lorem().Sentence(10)
				resourceAmount := 100 + rand.IntN(10000)

				// Prepare expected response with randomized data
				responseID := "resource-" + fake.UUID().V4()
				responseName := "response-" + fake.Person().Name()
				responseTitle := fake.Lorem().Word()
				responseDesc := fake.Lorem().Sentence(10)
				responseTags := []string{fake.Lorem().Word(), fake.Lorem().Word()}
				responseStatus := "active"
				createdAt := fake.Time().ISO8601()
				updatedAt := fake.Time().ISO8601()

				expectedResponse := &Resource{
						ID:          responseID,
						Name:        responseName,
						Title:       responseTitle,
						Description: responseDesc,
						Tags:        responseTags,
						Status:      responseStatus,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
				}

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// Verify request details
						assert.Equal(t, "POST", r.Method)
						assert.Equal(t, "/resources", r.URL.Path)
						assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

						// Return complete successful response
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)

						fmt.Fprintf(w, `{
								"id": "%s",
								"name": "%s",
								"title": "%s",
								"description": "%s",
								"tags": ["%s", "%s"],
								"status": "%s",
								"created_at": "%s",
								"updated_at": "%s"
						}`, responseID, responseName, responseTitle, responseDesc, responseTags[0], responseTags[1], responseStatus, createdAt, updatedAt)
				}))
				defer server.Close()

				deps := makeMockDeps(t, server.URL)
				client := NewClient(deps)

				req := &CreateResourceRequest{
						Name:        resourceName,
						Description: resourceDesc,
						Amount:      resourceAmount,
						Tags:        []string{fake.Lorem().Word(), fake.Lorem().Word()},
				}

				// Act
				resource, err := client.CreateResource(t.Context(), CreateResourceParams{
						Request: req,
				})

				// Assert
				require.NoError(t, err)
				// Compare entire structs to avoid field-by-field assertions
				assert.Equal(t, expectedResponse, resource)
		})

		t.Run("success with required parameters only", func(t *testing.T) {
				// Arrange - Use randomized data
				resourceName := "resource-" + fake.Lorem().Word()

				// Prepare expected minimal response with randomized data
				expectedID := "resource-" + fake.UUID().V4()
				expectedName := "minimal-" + fake.Lorem().Word()

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// Return minimal successful response with randomized data
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusCreated)
						fmt.Fprintf(w, `{
								"id": "%s",
								"name": "%s"
						}`, expectedID, expectedName)
				}))
				defer server.Close()

				deps := makeMockDeps(t, server.URL)
				client := NewClient(deps)

				req := &CreateResourceRequest{
						Name: resourceName, // Only required field
				}

				// Act
				resource, err := client.CreateResource(t.Context(), CreateResourceParams{
						Request: req,
				})

				// Assert
				require.NoError(t, err)
				// Check only the fields that should be present in minimal response
				assert.Equal(t, expectedID, resource.ID)
				assert.Equal(t, expectedName, resource.Name)
		})

		t.Run("handles API error", func(t *testing.T) {
				// Arrange
				resourceName := "resource-" + fake.Lorem().Word()

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusBadRequest)
				}))
				defer server.Close()

				deps := makeMockDeps(t, server.URL)
				client := NewClient(deps)

				req := &CreateResourceRequest{
						Name: resourceName,
				}

				// Act
				result, err := client.CreateResource(t.Context(), CreateResourceParams{
						Request: req,
				})

				// Assert
				require.Error(t, err)
				assert.Nil(t, result)
				assert.ErrorContains(t, err, "create resource failed")
		})
}
```

## Quality Assurance Guidelines

### 1. Testing Requirements
- Always include the 3 standard test cases per operation
- Use faker for generating random test data
- Follow testing-best-practices patterns
- Use AAA (Arrange-Act-Assert) pattern with clear comments
- Include test name in logger for better debugging
- Use proper error assertions (assert.ErrorContains or assert.ErrorIs)

### 2. Documentation Requirements
- Document all public types and methods
- Include usage examples in Go doc comments

### 3. Security Requirements
- Never log authentication tokens or sensitive data
- Use context for passing authentication tokens

### 4. Endpoints Documentation

Create an `ENDPOINTS.md` file in the same directory to document using instruction below.

#### ENDPOINTS.md File

Purpose of ENDPOINTS.md:

- **Quick Reference**: Easily see what endpoints are implemented
- **Maintenance**: Track which API operations are available
- **Updates**: Identify what needs to be added or modified
- **Onboarding**: Help new developers understand the client scope

The file must follow the below format:

```markdown
# [Service Name] API Client Endpoints

POST /path/to/resource
Client method: CreateResource(ctx, CreateResourceParams)

GET /path/to/resource/{id}
Client method: GetResource(ctx, GetResourceParams)

.....e.t.c.....
```

#### When to Update

- After adding new endpoints
- After removing deprecated endpoints
- When refactoring endpoint signatures
- Before major releases

### 5. Code Quality and Linting

**IMPORTANT**: After completing implementation, always run linting and tests to ensure code quality:

```bash
make lint
make test
```

Failing lint or test is an indication of task **FAILURE**.

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
       PetId string
   }

   // ✅ Good: acronyms should be all uppercase
   type ShowPetByIDParams struct {
       PetID string
   }
   ```

3. **Context Usage in Tests (usetesting)**
   ```go
   // ❌ Bad: using context.Background() in tests
   ctx := context.Background()

   // ✅ Good: use t.Context() in tests for better test lifecycle management
   ctx := t.Context()
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

#### Task Completion Protocol

Follow established project guidelines including the following:

1. Run `make lint` and fix remaining issues
2. Run `make test` to ensure all tests pass

Report the result to the user as per project guidelines.

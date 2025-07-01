# Engineering Requirements Document: Removing Template Tools

## Introduction/Overview

This document outlines the requirements for removing unnecessary template-generated tools from the codebase. Specifically, the Math and Time tools that were included in the boilerplate template are not required for the project and should be removed to clean up the codebase.

## Business Logic

The Math and Time tools are example tools that were included in the template to demonstrate how to implement MCP tools. They are not part of the core business functionality of the application and are therefore candidates for removal.

## High Level Architecture

The removal process will involve:

1. Removing the controller files for Math and Time tools
2. Removing the service files for Math and Time functionality
3. Updating registration code to remove references to these components
4. Ensuring that the removal doesn't break any existing functionality

## Detailed Architecture

### Components to Remove

#### Controller Layer
1. `internal/api/mcp/controllers/math.go`
2. `internal/api/mcp/controllers/math_test.go`
3. `internal/api/mcp/controllers/time.go`
4. `internal/api/mcp/controllers/time_test.go`

#### Application Layer
1. `internal/app/math.go`
2. `internal/app/math_test.go`
3. `internal/app/time.go`
4. `internal/app/time_test.go`

### Files to Update

1. `internal/api/mcp/controllers/register.go`
   - Remove imports for Math and Time controllers
   - Remove registration of Math and Time controllers in the `Register` function

2. `internal/app/register.go`
   - Remove registration of Math and Time services in the `Register` function

### Dependency Considerations

The Time service might be used by other components for time-related operations. However, the codebase has a separate `TimeProvider` in the `internal/services` package, which will be retained for actual time-related functionality instead of the template's `TimeService`.

## Key Architectural Decisions

1. **Complete Removal**: We will completely remove the Math and Time tools rather than just disabling them, as they are not needed for the project's functionality.

2. **Preserving Core Structure**: The removal will maintain the overall MCP controller structure and registration pattern for remaining tools (like Bitbucket).

3. **Handling Time Dependencies**: We will retain the `TimeProvider` in the `internal/services` package as it may be required by other components in the system.

4. **Scope Limitation**: Only the Math and Time tools will be removed; all other template components will be kept intact.

## Testing Strategy

1. **Unit Tests**: Since we're removing code, we'll also remove the associated test files.

2. **Integration Tests**: After removal, run integration tests to ensure the remaining MCP tools still function correctly.

3. **Manual Testing**: Verify that the MCP server starts correctly and that the removed tools are no longer available in the API.

4. **Regression Testing**: Ensure that removing these components doesn't affect other parts of the system.

## Implementation Steps

1. Delete the controller files:
   - `internal/api/mcp/controllers/math.go`
   - `internal/api/mcp/controllers/math_test.go`
   - `internal/api/mcp/controllers/time.go`
   - `internal/api/mcp/controllers/time_test.go`

2. Delete the service files:
   - `internal/app/math.go`
   - `internal/app/math_test.go`
   - `internal/app/time.go`
   - `internal/app/time_test.go`

3. Update the registration files:
   - Remove Math and Time controller registrations from `internal/api/mcp/controllers/register.go`
   - Remove Math and Time service registrations from `internal/app/register.go`

4. Run tests to verify that the system still works correctly after the removals. 
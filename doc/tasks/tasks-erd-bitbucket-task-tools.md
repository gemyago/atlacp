## Requirements

Please read referenced files to understand the problem:
- `doc/erd-bitbucket-task-tools.md`

## Relevant Files

### Source Files
- `internal/app/ports.go` - Update the bitbucketClient interface
- `internal/app/bitbucket.go` - Add new parameter structs
- `internal/app/bitbucket.go` - Implement service methods for task operations
- `internal/services/bitbucket/list_tasks.go` - Client implementation (already exists)
- `internal/services/bitbucket/update_task.go` - Client implementation (already exists)
- `internal/api/mcp/controllers/bitbucket.go` - Add new MCP tools for task operations

### Test Files
- `internal/app/bitbucket_test.go` - Unit tests for BitbucketService task methods
- `internal/api/mcp/controllers/bitbucket_test.go` - Unit tests for bitbucket controller tools

### Notes

- **Testing Framework:** Use testify for assertions and table-driven tests
- **Architecture:** Follow clean architecture with layers: Controller → App → Services
- **Mock Generation:** Use mockery to generate updated mocks if needed
- **Follow TDD process:**
  1. Define interfaces and struct definitions first
  2. Write test for the specific functionality
  3. Create minimal stub implementation that compiles but fails the test
  4. Implement the full functionality to pass the test
  5. Refactor if needed
- **Existing Architecture:** Follow the established pattern of controller → app service → client

## Tasks

- [x] 1.0 Update App Layer Interfaces
  - [x] 1.1 Extend bitbucketClient interface with ListPullRequestTasks and UpdateTask methods
  - [x] 1.2 Regenerate mocks if needed

- [ ] 2.0 Implement App Layer Services
  - [ ] 2.1 Create BitbucketListTasksParams struct in bitbucket.go
  - [ ] 2.2 Create minimal stub implementation for ListTasks that compiles but fails the test
  - [ ] 2.3 Write test for BitbucketService.ListTasks method
  - [ ] 2.4 Implement BitbucketService.ListTasks method to pass the test
  - [ ] 2.5 Create BitbucketUpdateTaskParams struct in bitbucket.go
  - [ ] 2.6 Create minimal stub implementation for UpdateTask that compiles but fails the test
  - [ ] 2.7 Write test for BitbucketService.UpdateTask method
  - [ ] 2.8 Implement BitbucketService.UpdateTask method to pass the test
  - [ ] 2.9 Review tests if they follow [testing-best-practices](../testing-best-practices.md)
  - [ ] 2.10 Run `lint-and-test`

- [ ] 3.0 Implement Controller Layer for Task Tools
  - [ ] 3.1 Update bitbucketService interface in the controller to include new methods
  - [ ] 3.2 Regenerate mocks if needed
  - [ ] 3.3 Create minimal stub implementation for list_pr_tasks tool that compiles but fails the test
  - [ ] 3.4 Write test for bitbucket_list_pr_tasks tool
  - [ ] 3.5 Implement bitbucket_list_pr_tasks tool to pass the test
  - [ ] 3.6 Create minimal stub implementation for update_pr_task tool that compiles but fails the test
  - [ ] 3.7 Write test for bitbucket_update_pr_task tool
  - [ ] 3.8 Implement bitbucket_update_pr_task tool to pass the test
  - [ ] 3.9 Register new tools in NewTools method 
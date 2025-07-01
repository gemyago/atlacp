## Requirements

Please read referenced files to understand the problem:
- `doc/erd-remove-template-tools.md`

## Relevant Files

### Source Files
- `internal/api/mcp/controllers/math.go` - Math controller to be removed
- `internal/api/mcp/controllers/time.go` - Time controller to be removed
- `internal/api/mcp/controllers/register.go` - Registration file to update
- `internal/app/math.go` - Math service to be removed
- `internal/app/time.go` - Time service to be removed
- `internal/app/register.go` - Application registration file to update

### Test Files
- `internal/api/mcp/controllers/math_test.go` - Math controller tests to be removed
- `internal/api/mcp/controllers/time_test.go` - Time controller tests to be removed
- `internal/app/math_test.go` - Math service tests to be removed
- `internal/app/time_test.go` - Time service tests to be removed

### Notes

- **Testing Strategy:** After removal, run existing tests to ensure remaining MCP tools still function correctly
- **Dependency Considerations:** Retain the `TimeProvider` in `internal/services` package as it may be required by other components

## Tasks

- [ ] 1.0 Remove Math and Time Controllers and Update Registration
  - [ ] 1.1 Update `internal/api/mcp/controllers/register.go` to remove Math and Time controller imports
  - [ ] 1.2 Update `internal/api/mcp/controllers/register.go` to remove Math and Time controller registrations
  - [ ] 1.3 Delete `internal/api/mcp/controllers/math.go`
  - [ ] 1.4 Delete `internal/api/mcp/controllers/math_test.go`
  - [ ] 1.5 Delete `internal/api/mcp/controllers/time.go`
  - [ ] 1.6 Delete `internal/api/mcp/controllers/time_test.go`
  - [ ] 1.7 Run `make lint test` to verify controller layer still works correctly
  - [ ] 1.8 Review and ensure no other files depend on these controllers

- [ ] 2.0 Remove Math and Time Services and Update Registration
  - [ ] 2.1 Update `internal/app/register.go` to remove Math and Time service registrations
  - [ ] 2.2 Delete `internal/app/math.go`
  - [ ] 2.3 Delete `internal/app/math_test.go`
  - [ ] 2.4 Delete `internal/app/time.go`
  - [ ] 2.5 Delete `internal/app/time_test.go`
  - [ ] 2.6 Run `make lint test` to verify service layer still works correctly
  - [ ] 2.7 Verify that removing these services doesn't affect other components

- [ ] 3.0 Final Verification
  - [ ] 3.1 Run `make lint test` to verify successful compilation of the entire project
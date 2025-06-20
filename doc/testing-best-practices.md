# Testing Best Practices

## Core Principles

### Follow TDD
- Iterate with small chunks of logic/code
- Add stub implementation first if needed
- Write test to cover new logic or new behavior
- Run test to see if it fails
- Implement the minimal amount of code to make the test pass
- Repeat the process until your code is complete

### Test What Matters
- **Focus on business logic** - Test the core functionality your code needs to provide
- **Avoid excessive tests** - Don't test scenarios that aren't relevant to your actual use cases
- **Test behavior, not implementation** - Focus on what the code does, not how it does it

### Keep It Simple
- **Pragmatic mocks** - Use simple mock implementations, avoid over-engineering
- **Minimal test setup** - Only setup what's necessary for the specific test case
- **Clear test names** - Use descriptive names that explain what behavior is being tested

### Don't Test the Framework
- **Skip infrastructure testing** - Don't test that logging works, that HTTP requests work, etc.
- **Trust the standard library** - Don't test Go's built-in functionality
- **Focus on your logic** - Test the decisions and transformations your code makes

## Patterns to Follow

### Common principles

- Avoid static variables shared across tests
- Use random data when possible, use faker (github.com/go-faker/faker/v4)
- Don't pollute testing namespace - if helper functions are only used within one test, nest them inside that test function

### Use makeMockDeps()

If your component has dependencies, use pattern below:
```go
func TestMyService(t *testing.T) {
    makeMockDeps := func() MyServiceDeps {
        return MyServiceDeps{
            RootLogger: diag.RootTestLogger(),
        }
    }
    
    t.Run("some test", func(t *testing.T) {
        deps := makeMockDeps()
        // ... use deps
    })
}
```

### Simple Mock Implementation
```go
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    args := m.Called(req)
    res, _ := args.Get(0).(*http.Response)  // Simple, pragmatic
    return res, args.Error(1)
}
```

Larger structs/interfaces can be mocked with mockery by adding a new interface to [.mockery.yaml](../.mockery.yaml) and running `mockery` to generate the mock (from root folder)

## Test Structure

### Use Nested t.Run
```go
func TestMyService(t *testing.T) {
    t.Run("should handle valid input", func(t *testing.T) {
        // Test the main functionality
    })
    
    t.Run("should handle missing data", func(t *testing.T) {
        // Test edge case that matters to business logic
    })
}
```

### Follow AAA Pattern
- **Arrange** - Set up test data and mocks
- **Act** - Call the code under test
- **Assert** - Verify the expected behavior

## Common Mistakes to Avoid

1. **Over-testing** - Testing every possible combination when only a few matter
2. **Testing implementation details** - Checking internal state instead of behavior
3. **Complex mocks** - Over-engineered mock implementations
4. **Testing the framework** - Verifying that standard library functions work
5. **Inconsistent patterns** - Not following established codebase conventions

## Remember

> "Test the code you wrote, not the code you didn't write." 

Focus on your business logic and keep tests simple, relevant, and maintainable. 
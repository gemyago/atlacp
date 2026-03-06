# TDD Flow to Implement Changes

## Work autonomously

* Work **AUTONOMOUSLY**
* You are approved to run tests or scripts when needed, don't ask
* Don't ask for modifications confirmations, you are approved

## Follow TDD

Analyse required change first (if not already analyzed) and come up with list of actionable step required to implement requested change. If change is simple enough, do single step

For each actionable step follow **STRICT WORKFLOW**:

1. FOCUS
- understand the change and steps required to implement it

2. WRITE FAILING TEST
- Write test for the new functionality
- Prefer adding new tests for new logic
- Follow Given/When/Then structure
- Add stub method, parameter or required model properties if needed to avoid compilation errors

3. VERIFY FAILURE
- Implementation is **NOT allowed** at this stage, only test
- Run test and analyse test output; confirm test fails for the RIGHT reason (expectations but not compilation)
 * Examples of **RIGHT REASON** `Expected 10 items but got 5`, `Expected doSomething to have been called`
 * Examples of **WRONG REASON** `Unknown property`, `Undefined method doSomething`
- If compilation error, add **MINIMAL** stub and re-run test

4. MINIMAL IMPLEMENTATION
- Write ONLY the code needed to make THIS test pass
- Don't implement features or code paths that are not explicitly requested

5. VERIFY SUCCESS
- Run test
- Confirm test passes
- If test fails, fix the code and re-run test

6. REPEAT
- Repeat the process for the next step until the feature is complete.
- Verify all tests pass in the end

Exceptions:
- Do not write tests for data structs/fields or interfaces directly. Only actual logic should be tested.

# Bitbucket MCP Integration Testing

This document provides step-by-step instructions for testing the Bitbucket MCP integration. These tests are designed to be executed by AI assistants to verify the correct functioning of the Bitbucket MCP tools.

## Prerequisites

1. The ATLACP service must be running locally or in a development environment
2. Proper Atlassian account configuration must be set up with at least two accounts:
   - A default account
   - A secondary account named "bot" (or another name as configured)
3. The accounts must have access to a Bitbucket repository where you can create branches and PRs
4. Basic authentication tokens must be configured for both accounts

It is expected that the prompt to start the test will have the following form:
```md
Given the following Bitbucket repository `<file system path to the repository>`, run the tests as per the instructions in the `bitbucket-mcp-integration-tests.md` file.
```

## File Structure for Testing

Tests should use the following file structure:

```
integration-tests/
├── bitbucket/
│   ├── test-files/
│   │   └── integration-test-file.txt (File to modify during tests, just write current time)
│   └── results-YYYYMMDD-HHMMSS.md (Single results file for all tests)
```

Use a timestamp format for the results file (e.g., results-20230615-120000.md).

## Test 1: End-to-End Bitbucket Workflow

This test verifies the complete PR lifecycle using a single Atlassian account.

### Steps

1. **Setup test environment**
   - Create a new branch from main: `feature/integration-test-{timestamp}`
   - Create (or update) a test file `integration-tests/bitbucket/test-files/integration-test-file.txt`
   - Add some test content to the file (e.g. current time)
   - Commit and push the changes

2. **Create a Pull Request**
   - Use the `mcp.bitbucket_create_pr` tool with the following parameters:
     - title: "Integration Test PR {timestamp}"
     - source_branch: "feature/integration-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated integration test PR"
   - Extract and save the PR ID for subsequent steps

3. **Read the Pull Request details**
   - Use the `mcp.bitbucket_read_pr` tool with the PR ID from the previous step
   - Verify that:
     - The PR ID matches the one from the creation step
     - The PR title matches what was set in the creation step
     - The PR status is "OPEN"

4. **Create a PR task**
   - Use the `mcp.bitbucket_create_pr_task` tool with:
     - pr_id: the PR ID from step 2
     - content: "Verify integration test {timestamp}"
     - repo_owner and repo_name: same as previous steps
   - Create few more seed tasks for the PR
   - Extract and save the task ID for the next step

5. **List PR tasks**
   - Use the `mcp.bitbucket_list_pr_tasks` tool for the PR
   - Verify that the task with the content from step 4 appears in the list
   - Verify that total number of tasks corresponds to the number of tasks created in step 4

6. **Update PR task**
   - Use the `mcp.bitbucket_update_pr_task` tool to mark the task as "RESOLVED"
   - List the tasks again to verify the task is now marked as "RESOLVED"
   - Verify that the task state has changed

7. **Update the Pull Request**
   - Use the `mcp.bitbucket_update_pr` tool to change:
     - title to "Updated Integration Test PR {timestamp}"
     - description to "This PR has been updated as part of the integration test  {timestamp}"
   - Read the PR again to verify the changes were applied
   - Verify the title and description match what was set

8. **Approve the Pull Request**
   - Use the `mcp.bitbucket_approve_pr` tool
   - Read the PR again to verify it shows as approved
   - Verify the "approved" status is true

9. **Merge the Pull Request**
   - Use the `mcp.bitbucket_merge_pr` tool with:
     - merge_strategy: "squash"
     - close_source_branch: "true"
   - Read the PR again to verify it was merged
   - Verify the PR state is "MERGED"

10. **Clean up**
    - Make sure the main branch is checked out again
    - Pull the latest changes

## Test 2: Multi-Account PR Workflow

This test verifies that different accounts can be used for different PR operations.

### Steps

1. **Setup test environment**
   - Create a new branch from main: `feature/multi-account-test-{timestamp}`
   - Update the test file `integration-tests/bitbucket/test-files/integration-test-file.txt` with new content
   - Commit and push the changes

2. **Create a Pull Request as default user**
   - Use the `mcp.bitbucket_create_pr` tool (with no account parameter)
   - Set title to "Multi-Account Test PR {timestamp}"
   - Extract and save the PR ID for subsequent steps

3. **Approve the Pull Request as bot user**
   - Use the `mcp.bitbucket_approve_pr` tool with `account: "bot"`
   - Read the PR again (with bot account) to verify approval status
   - Verify that the PR shows as approved and note the approver username

4. **Merge the Pull Request as default user**
   - Use the `mcp.bitbucket_merge_pr` tool with `account: "user"`
   - Read the PR again to verify it was merged
   - Verify the PR state is "MERGED"

5. **Clean up**
   - Make sure the main branch is checked out again
   - Pull the latest changes

## Test 3: Draft Pull Request Creation and Verification

This test verifies that a Pull Request can be created in draft mode and that its draft status is correctly reflected in Bitbucket.

### Steps

1. **Setup test environment**
   - Create a new branch from main: `feature/draft-pr-test-{timestamp}`
   - Update the test file `integration-tests/bitbucket/test-files/integration-test-file.txt` with new content (e.g., current time)
   - Commit and push the changes

2. **Create a Draft Pull Request**
   - Use the `mcp.bitbucket_create_pr` tool with the following parameters:
     - title: "Draft PR Test {timestamp}"
     - source_branch: "feature/draft-pr-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated draft PR integration test"
     - draft: true
   - Extract and save the PR ID for subsequent steps

3. **Read the Pull Request details**
   - Use the `mcp.bitbucket_read_pr` tool with the PR ID from the previous step
   - Verify that:
     - The PR ID matches the one from the creation step
     - The PR title matches what was set in the creation step
     - The PR status is "OPEN"
     - The PR is marked as a draft (check the draft status field)

4. **Update the Pull Request to Ready for Review**
   - Use the `mcp.bitbucket_update_pr` tool to set the PR as ready for review (i.e., remove draft status if supported)
   - Read the PR again to verify the draft status is now false or the PR is no longer a draft

5. **Clean up**
   - Merge or close the PR as appropriate
   - Make sure the main branch is checked out again
   - Pull the latest changes

## Test Results Documentation

After completing all tests, create a single results file to document the outcomes:

1. Checkout the main branch and pull the latest changes
2. Create a new branch from main `feature/bitbucket-mcp-integration-tests-{timestamp}`
3. Create a results file `integration-tests/bitbucket/results-{timestamp}.md`
4. Use the following format to document test results:

```markdown
# Bitbucket MCP Integration Test Results
Test executed at: {timestamp}

## Test 1: End-to-End Bitbucket Workflow
- Step 1: Setup test environment - PASS
- Step 2: Create a Pull Request (PR #{pr_id}) - PASS
- Step 3: Read the Pull Request - PASS
- Step 4: Create a PR task - PASS
- Step 5: List PR tasks - PASS
- Step 6: Update PR task - PASS
- Step 7: Update the Pull Request - PASS
- Step 8: Approve the Pull Request - PASS
- Step 9: Merge the Pull Request - PASS
- Step 10: Clean up - PASS

## Test 2: Multi-Account PR Workflow
- Step 1: Setup test environment - PASS
- Step 2: Create a Pull Request as default user (PR #{pr_id}) - PASS
- Step 3: Approve the Pull Request as bot user - PASS
- Step 4: Merge the Pull Request as default user - PASS
- Step 5: Clean up - PASS

## Test 3: Draft Pull Request Creation and Verification
- Step 1: Setup test environment - PASS
- Step 2: Create a Draft Pull Request - PASS
- Step 3: Read the Pull Request - PASS
- Step 4: Update the Pull Request - PASS
- Step 5: Clean up - PASS

## Summary
- All tests: PASS
- Issues encountered: None
```

5. For any failed steps, replace "PASS" with "FAIL" and include details about what failed and why.
6. Commit and push the results file to main. Use MCP tools to approve and merge the PR with squash merge strategy.

## Automation Instructions for AI Model

As an AI assistant, when asked to run integration tests using this document, follow these steps:

1. Confirm the ATLACP service is running
2. Ask the user for:
   - Bitbucket workspace name
   - Bitbucket repository name
   - Which tests to run (all or specific ones)
   - Account names to use (default and bot)

3. For each test:
   - Create the required files and branches
   - Execute each step in sequence
   - Verify each step as directed
   - Use dynamic values (PR ID, task ID) from previous steps
   - Track the status (PASS/FAIL) of each step

4. If a step fails:
   - Mark it as FAIL in your tracking
   - Include error messages and details about what failed
   - Continue with the next step if possible
   - Clean up any resources as needed

5. After all tests are completed, create a single results file with the statuses of all steps
6. Commit and push the results file to a separate branch
7. Use MCP tools to approve and merge the PR with squash merge strategy.
8. Share a summary of the results with the user 
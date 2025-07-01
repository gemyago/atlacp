# Bitbucket MCP Integration Testing

This document provides step-by-step instructions for testing the Bitbucket MCP integration. These tests are designed to be executed by AI assistants to verify the correct functioning of the Bitbucket MCP tools.

**Important**: If otherwise mentioned - figure out all the details of the repository and files **yourself** before you start the tests. Feel free to do it by running various git commands to get repository details. The user will provide the file system path to the repository and that's it.

**Important**: If not mentioned, run all tests from this file.

**Very Important**: The user **MUST** provide the file system path to the repository. Current workspace is NOT the bitbucket repository. If not provided - ask the user for the path.

It is expected that the prompt to start the test will have the following form:
```markdown
Given the following Bitbucket repository `<file system path to the repository>`, run the tests as per the instructions in the `bitbucket-mcp-integration-tests.md` file.
```

## Prerequisites done by the user

It should be assumed that below is already prepared by the user:
1. The ATLACP service must be running and MCP tools are registered and available to AI assistant.
2. Proper Atlassian account configuration is set up with at least two accounts:
   - A default account - assume it is named "user" if not otherwise mentioned
   - A secondary account named "bot"

## Working with the repository

You need to use shell commands to work with the repository and and all files. Some of the commands are:

```bash
cd "<file system path to the repository>"

git pull

# Create a new branch
git checkout -b "feature/integration-test-{timestamp}"

# Add a new file
touch "integration-tests/bitbucket/test-files/integration-test-file.txt"

# Use echo and cat to write and read the file
echo "Current time: $(date)" > "integration-tests/bitbucket/test-files/integration-test-file.txt"

# Use any other techniques to write to the file as needed

# Commit the changes
git add "integration-tests/bitbucket/test-files/integration-test-file.txt"

git commit -m "Add integration test file"

# Push and set upstream
git push origin "feature/integration-test-{timestamp}" --set-upstream
```

**Important**: You need to use the shell commands to work with the repository and and all files. Do not use the file system MCP tools to work with the repository and files. You will not be able to since they're in a different folder (other than current workspace folder).

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

## Test 1: PR Creation and Updates

This test verifies the PR creation, reading, updating, approval, and merging using a single Atlassian account.

### Steps

1. **Setup test environment**
   - Create a new working branch from main: `feature/pr-lifecycle-test-{timestamp}`
   - Create (or update) a test file `integration-tests/bitbucket/test-files/integration-test-file.txt`
   - Add some test content to the file (e.g. current time)
   - Commit the changes
   - Update and commit the changes again to have two commits in the branch
   - Note the commit hashes for both commits
   - Push the changes to the branch

2. **Create a Pull Request**
   - Use the `mcp.bitbucket_create_pr` tool with the following parameters:
     - title: "PR Lifecycle Test {timestamp}"
     - source_branch: "feature/pr-lifecycle-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated PR lifecycle test (Test 1)"
   - Extract and save the PR ID for subsequent steps

3. **Read the Pull Request details**
   - Use the `mcp.bitbucket_read_pr` tool with the PR ID from the previous step
   - Verify that:
     - The PR ID matches the one from the creation step
     - The PR title matches what was set in the creation step
     - The PR status is "OPEN"

4. **Update the Pull Request**
   - Use the `mcp.bitbucket_update_pr` tool to change:
     - title to "Updated PR Lifecycle Test {timestamp}"
     - description to "This PR has been updated as part of the PR lifecycle test (Test 1) {timestamp}"
   - Read the PR again to verify the changes were applied
   - Verify the title and description match what was set

5. **Approve the Pull Request**
   - Use the `mcp.bitbucket_approve_pr` tool
   - Read the PR again to verify it shows as approved
   - Verify the "approved" status is true

6. **Merge the Pull Request**
   - Use the `mcp.bitbucket_merge_pr` tool with:
     - merge_strategy: "squash"
     - merge_message: "Squash merged PR Lifecycle Test (Test 1) {timestamp}"
     - close_source_branch: "true"
   - Read the PR again to verify it was merged
   - Verify the PR state is "MERGED"

7. **Clean up**
   - Make sure the main branch is checked out again
   - Pull the latest changes
   - Delete the working branch
   - Review the commit history to verify that there is a single commit (Squash merge)
   - Ensure commits from step 1 are not present in the main branch

Prepare a report as per instruction further down in this instruction.

## Test 2: PR Tasks Management

This test verifies the PR tasks creation, listing, and updating functionality.

### Steps

1. **Setup test environment**
   - Create a new working branch from main: `feature/pr-tasks-test-{timestamp}`
   - Create (or update) a test file `integration-tests/bitbucket/test-files/integration-test-file.txt`
   - Add some test content to the file (e.g. current time)
   - Commit and push the changes

2. **Create a Pull Request**
   - Use the `mcp.bitbucket_create_pr` tool with the following parameters:
     - title: "PR Tasks Test {timestamp}"
     - source_branch: "feature/pr-tasks-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated PR tasks test (Test 2)"
   - Extract and save the PR ID for subsequent steps

3. **Create multiple PR tasks**
   - Use the `mcp.bitbucket_create_pr_task` tool with:
     - pr_id: the PR ID from step 2
     - content: "Task 1: Verify integration test {timestamp}"
     - repo_owner and repo_name: same as previous steps
   - Create two more tasks with different content:
     - "Task 2: Review code changes {timestamp}"
     - "Task 3: Test functionality {timestamp}"
   - Extract and save at least one task ID for the next steps

4. **List PR tasks**
   - Use the `mcp.bitbucket_list_pr_tasks` tool for the PR
   - Verify that all three tasks created in step 3 appear in the list
   - Verify that total number of tasks corresponds to the number of tasks created in step 3

5. **Update PR tasks**
   - Use the `mcp.bitbucket_update_pr_task` tool to mark "Task 1" as "RESOLVED"
   - Use the `mcp.bitbucket_update_pr_task` tool to update the content of "Task 2" to "Task 2: Code review completed {timestamp}"
   - List the tasks again to verify:
     - "Task 1" is now marked as "RESOLVED"
     - "Task 2" content has been updated
     - "Task 3" remains unchanged

6. **Clean up**
   - Approve the PR using the `mcp.bitbucket_approve_pr` tool
   - Merge the PR using the `mcp.bitbucket_merge_pr` tool with:
     - merge_strategy: "squash"
     - close_source_branch: "true"
   - Make sure the main branch is checked out again
   - Pull the latest changes
   - Delete the working branch

Prepare a report as per instruction further down in this instruction.

## Test 3: Multi-Account PR Workflow

This test verifies that different accounts can be used for different PR operations.

### Steps

1. **Setup test environment**
   - Create a new working branch from main: `feature/multi-account-test-{timestamp}`
   - Update the test file `integration-tests/bitbucket/test-files/integration-test-file.txt` with new content
   - Commit and push the changes
   - Note the commit hash

2. **Create a Pull Request as default user**
   - Use the `mcp.bitbucket_create_pr` tool (with no account parameter)
   - Set title to "Multi-Account Test PR (Test 2) {timestamp}"
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
   - Delete the working branch
   - Ensure the commit from step 1 is present in addition to the merge commit

Prepare a report as per instruction further down in this instruction.

## Test 4: Draft Pull Request Creation and Verification

This test verifies that a Pull Request can be created in draft mode and that its draft status is correctly reflected in Bitbucket.

### Steps

1. **Setup test environment**
   - Create a new working branch from main: `feature/draft-pr-test-{timestamp}`
   - Update the test file `integration-tests/bitbucket/test-files/integration-test-file.txt` with new content (e.g., current time)
   - Commit and push the changes. Note the commit hash.

2. **Create a Draft Pull Request**
   - Use the `mcp.bitbucket_create_pr` tool with the following parameters:
     - title: "Draft PR Test (Test 3) {timestamp}"
     - source_branch: "feature/draft-pr-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated draft PR integration test (Test 3)"
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

3. **Merge the Pull Request**
   - Use the `mcp.bitbucket_merge_pr` tool with:
     - merge_strategy: "fast_forward"
   - Read the PR again to verify it was merged
   - Verify the PR state is "MERGED"

5. **Clean up**
   - Merge or close the PR as appropriate
   - Make sure the main branch is checked out again
   - Pull the latest changes
   - Delete the working branch
   - Ensure the commit from step 1 is present and no merge commit is present after the merge

Prepare a report as per instruction further down in this instruction.

## Test Results Reporting

Follow the protocol below when performing the test:
* In a **current workspace** create a file `tmp/integration-tests-{timestamp}-results.md`
* With each step - update the file (see format below) mentioning the step number, description, status (PASS/FAIL)
* If failed - comment what failed and why
* **Do not stop** if any step fails, document and continue

**Format of the results file**
```markdown
# Bitbucket MCP Integration Test Results
Test executed at: {timestamp}

## Test 1: PR Creation and Updates
- Step 1: <Step description> - PASS
- Step 2: <Step 2 description> (Pull Request (PR #{pr_id}) - PASS
- Step 3: <Step 3 description> - PASS
......

## Test 2: PR Tasks Management
- Step 1: <Step description> - PASS
- Step 2: <Step 2 description> (Pull Request (PR #{pr_id}) - PASS
- Step 3: <Step 3 description> - PASS
......

## Test 3: Multi-Account PR Workflow
- Step 1: <Step description> - PASS
- Step 2: <Step 2 description> (Pull Request (PR #{pr_id}) - PASS
- Step 3: <Step 3 description> - PASS
......

## Test 4: Draft Pull Request Creation and Verification
- Step 1: <Step description> - PASS
- Step 2: <Step 2 description> (Pull Request (PR #{pr_id}) - PASS
- Step 3: <Step 3 description> - PASS

## Summary
- All tests: PASS/FAIL
- Issues encountered: None/List issues
```

When completed all tests, copy the results file from a **current workspace** to the integration tests repository. Do steps below:
1. Checkout the main branch and pull the latest changes
2. Create a new branch from main `feature/bitbucket-mcp-integration-tests-results-{timestamp}`
3. Copy the results file from the **current workspace** to integration tests repository file `bitbucket/results-{timestamp}.md`
4. Commit and push the branch
5. Use MCP tool to create, approve and merge the PR:
    - title: "Bitbucket MCP Integration Tests Results {timestamp}"
    - description: "This is an automated integration tests results (Test 1, Test 2, ....., Test X) {timestamp}"
    - source_branch: "feature/bitbucket-mcp-integration-tests-results-{timestamp}"
    - target_branch: "main"
    - repo_owner: your workspace name
    - repo_name: your repository name
    - merge_strategy: "squash"
    - close_source_branch: "true"
6. Share a summary of the results with the user 

## Automation Instructions for AI Model

As an AI assistant, when asked to run integration tests using this document, follow these steps:

1. Confirm the ATLACP tools are available to you and you can use them.
2. For each tests follow the exact instruction
3. Document the results as you progress
4. Report the results as per instruction in the end of the process
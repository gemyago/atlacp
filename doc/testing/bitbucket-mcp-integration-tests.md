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
     - The PR is not a draft

4. **Update the Pull Request**
   - Use the `mcp.bitbucket_update_pr` tool to change:
     - title to "Updated PR Lifecycle Test {timestamp}"
     - description to "This PR has been updated as part of the PR lifecycle test (Test 1) {timestamp}"
   - Read the PR again to verify the changes were applied
   - Verify the title and description match what was set

5. **Approve the Pull Request**
   - Use the `mcp.bitbucket_approve_pr` tool
   - Read the PR again to verify it shows as approved
   - Verify the "approved" status is true in participants list

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

Update a report as per [instruction](#test-results-reporting).

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

Update a report as per [instruction](#test-results-reporting).

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

Update a report as per [instruction](#test-results-reporting).

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

4. **Update the Pull Request Draft Status**
   - Use the `mcp.bitbucket_update_pr` with just draft parameter set to false
   - Read the PR again to verify the draft status is now false
   - Use the `mcp.bitbucket_update_pr` with just draft parameter set to true again
   - Read the PR again to verify the draft status is now true
   - Use the `mcp.bitbucket_update_pr` with just draft parameter set to false
   - Read the PR again to verify the draft status is now false

5. **Merge the Pull Request**
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

Update a report as per [instruction](#test-results-reporting).

## Test 5: Bitbucket PR Review Tools End-to-End

This test verifies the end-to-end functionality of the Bitbucket PR review tools, including retrieving diffstat, fetching diffs, accessing file content, adding comments (general, inline, and pending), verifying pending status, and requesting changes on a pull request. This test now uses multiple files and multiline content to ensure robust coverage of review scenarios.

### Steps

1. **Setup test environment**
   - Create a new working branch from main: `feature/pr-review-tools-test-{timestamp}`
   - Copy the two TypeScript example files (`example1.ts` and `example2.ts`) from the current workspace to a unique subdirectory in `integration-tests/bitbucket/test-files/` (e.g., `integration-tests/bitbucket/test-files/ts-examples-{timestamp}/`).
   - For each file, ensure it is present and unmodified in the new location.
   - Commit and push the changes, including the two TypeScript files.
   - Note the commit hash.

2. **Create a Pull Request**
   - Use the `mcp.bitbucket_create_pr` tool with the following parameters:
     - title: "PR Review Tools Test {timestamp}"
     - source_branch: "feature/pr-review-tools-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated PR review tools test (Test 5)"
   - Extract and save the PR ID for subsequent steps
   - Ensure the PR includes the two TypeScript files and the three text files.
   - Use the PR diff to identify the exact line numbers for all `// Update marker` comments in both TypeScript files.
   - Add inline comments to each marker in both files using the diff line numbers. This will verify if comments are properly added for new files.

3. **List and verify PR comments and line numbers**
   - Use the `mcp.bitbucket_list_pr_comments` tool to list all comments for the PR.
   - For each comment on a TypeScript file:
     - Use the `mcp.bitbucket_get_file_content` tool to fetch the full file content from the PR branch.
     - Write the fetched file content to a uniquely named file in the `tmp/` directory in the current workspace, including a timestamp in the filename (e.g., `tmp/example1-<timestamp>.ts`).
     - Use the following one-liner script to determine the actual line numbers for all marker comments in the file:
       ```
       grep -n "Update marker" <filename> | cut -d: -f1,2
       ```
     - Compare the line numbers found by this script with the line numbers in the PR comments to verify mapping.
   - This approach ensures that the verification uses the actual file content as fetched from Bitbucket and provides a reproducible, timestamped record of the verification process.

4. **Create a derived branch and update markers**
   - Create a new branch from the PR branch: `feature/pr-review-tools-test-derived-{timestamp}`
   - For each `// Update marker` in both TypeScript files, update the marker by appending a timestamp or other unique text.
   - Commit and push the changes.

5. **Create a PR from the derived branch to the original PR branch**
   - Create a PR from `feature/pr-review-tools-test-derived-{timestamp}` to `feature/pr-review-tools-test-{timestamp}`

6. **Add inline comments to updated markers in the derived PR**
   - Use the diff tool to locate changes and identify the new line numbers for all updated markers.
   - Add inline comments to a selection of updated markers (at least one at the start, one in the middle, and one at the end of the files).

7. **List and verify comments and line numbers in the derived PR**
   - Use the `mcp.bitbucket_list_pr_comments` tool to list all comments for the derived PR.
   - For each comment:
     - Use the `mcp.bitbucket_get_file_content` tool to fetch the full file content from the derived PR branch.
     - Write the fetched file content to a uniquely named file in the `tmp/` directory in the current workspace, including a timestamp in the filename (e.g., `tmp/example1-derived-<timestamp>.ts`).
     - Use the following one-liner script to determine the actual line numbers for all marker comments in the file:
       ```
       grep -n "Update marker" <filename> | cut -d: -f1,2
       ```
     - Compare the line numbers found by this script with the line numbers in the PR comments to verify mapping.
   - This approach ensures that the verification uses the actual file content as fetched from Bitbucket and provides a reproducible, timestamped record of the verification process.

8. **Merge the derived PR**
   - Merge (squash) the derived PR into the original PR branch.
   - Continue with any other test steps that make sense for end-to-end review tool coverage.

6. **Add a general (non-inline) comment to the PR**
   - Use the `mcp.bitbucket_add_pr_comment` tool to add a general comment (not associated with a file or line) to the PR, such as "General comment for PR review tools test {timestamp}".
   - Verify that the general comment appears in the PR's comment list and is not associated with any file or line number.

Update a report as per [instruction](#test-results-reporting).

## Test Results Reporting

Follow the protocol below when performing the test:
1. In a **current workspace** create a file (if not yet exists) `tmp/integration-tests-{YYYYMMDD-HHMMSS}-results.md`
2. With each step - update the file (see format below) mentioning the step number, description, status (PASS/FAIL). Keep formatting.
3. If failed - comment what failed and why
4. **Do not stop** if any step fails, document and continue

### Format of the results file

```markdown
# Bitbucket MCP Integration Test Results
Test executed at: {timestamp}

## Test 1: <Title of the test>
- Step 1: <Step description> - PASS
- Step 2: <Step 2 description> (Pull Request (PR #{pr_id}) - PASS
- Step 3: <Step 3 description> - PASS
......

## Test 2: <Title of the test>
- Step 1: <Step description> - PASS
- Step 2: <Step 2 description> (Pull Request (PR #{pr_id}) - PASS
- Step 3: <Step 3 description> - PASS
......

<Other Reports in a same format>

## Summary
- All tests: PASS/FAIL
- Issues encountered: None/List issues
  - <Test 1> Short issue details
  - <Test 2> Short issue details
  ......

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
    - merge message: <same as pr title>
    - close_source_branch: "true"
6. Share a summary of the results with the user
7. Attach the results file URL to the PR in this repository as follows:
```bash
# In a current workspace, check if there is an active PR for this branch
gh pr view

# If there is no Active PR, DO NOTHING

# If there is an Active PR, attach the results file URL to the PR
# Results file url should have the following structure: 
# https://bitbucket.org/gemyago/atlacp-integration-tests/src/main/integration-tests/bitbucket/results-20250702-082610.md
gh pr comment <pr_id> --body "Integration tests results: <results file URL>"
```

## Automation Instructions for AI Model

As an AI assistant, when asked to run integration tests using this document, follow these steps:

1. Confirm the ATLACP tools are available to you and you can use them.
2. For each tests follow the exact instruction
3. Document the results as you progress
4. Report the results as per instruction in the end of the process
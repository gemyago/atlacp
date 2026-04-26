# Bitbucket MCP Integration Testing

This document provides step-by-step instructions for testing the Bitbucket MCP integration. These tests are designed to be executed by AI assistants to verify the correct functioning of the Bitbucket MCP tools.

**Important**: If not mentioned, run all tests from this file.

**Important**: The user may provide the file system path to the bitbucket repository. Current workspace is NOT the bitbucket repository. If not provided - try to see if `atlacp-integration-tests` repository exists in a parent folder of the current workspace. Once you get repository, `git remote show origin` will help you to figure-out repo owner and name. Do not attempt to recursively search for the repository in all parent folders, just check the repo as suggested above, if not found - ask the user for the path.

It is expected that the prompt to start the test will have the following form:
```markdown
Given the following Bitbucket repository `<file system path to the repository>`, run the tests as per the instructions in the `bitbucket-mcp-integration-tests.md` file.
```

### Execution model (mandatory, no exceptions)

- **Tests always run in sub-agents.** Every numbered test (Test 1, Test 2, …) **must** be executed **only** inside a **separate delegated sub-agent** (e.g. Cursor Task tool or equivalent). The orchestrator **must not** run that test’s git commands, Bitbucket MCP calls, or any step that mutates the integration repo or Bitbucket state.
- **Sequential execution across tests.** All tests share the **same** integration repository clone: they race if run **in parallel** (concurrent sub-agents or terminals) on `git checkout`, `main`, pushes, and shared paths. When running more than one test, the orchestrator **must** run them **one after another**—wait until the sub-agent for Test *N* has fully finished before spawning the sub-agent for the next test. (Parallelism is **not** allowed for the suite; only isolation per test via sub-agents.)
- **No waiver.** This rule applies if the user asks for one test, “just Test N”, “first test only”, or the full suite. Convenience, speed, or “single coherent thread” are **not** reasons to run test steps in the orchestrator.
- **Orchestrator-only work is allowed** (and encouraged where useful): planning order, spawning sub-agents, passing repo path / `repo_owner` / `repo_name` / test id, **verification that does not execute the test** (e.g. confirming sub-agents have MCP, merging `tmp/integration-tests-*.md` from outputs, summarizing pass/fail). The orchestrator does **not** substitute for a sub-agent when executing a test.
- **Why sub-agents**: isolates failures, avoids mixing shell/MCP steps for different tests in one thread, and matches review expectations. This is **not** permission to run multiple tests in parallel on the same repo.
- **Orchestrator responsibilities**: confirm Bitbucket MCP is available to sub-agents; when running multiple tests, launch sub-agents **sequentially** (wait for each to complete before the next); collect outputs and produce the final report.
- **Sub-agent responsibilities**: perform **all** shell work in the provided repo path and **all** `bitbucket_*` MCP calls for **that test only**; write/update the workspace results file for that test’s steps; return a concise pass/fail summary and PR ids to the orchestrator.

## Prerequisites done by the user

It should be assumed that below is already prepared by the user:
1. The ATLACP service must be running and MCP tools are registered and available to AI assistant.
2. Proper Atlassian account configuration is set up with at least two accounts:
   - A default account - assume it is named "user" if not otherwise mentioned
   - A secondary account named "bot"

For Git/SSH issues point the user on the [SSH troubleshooting](./README.md#ssh-troubleshooting). Don't read or do anything about it yourself, report and halt.

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
   - Use the `create Bitbucket pull request` action with the following parameters:
     - title: "PR Lifecycle Test {timestamp}"
     - source_branch: "feature/pr-lifecycle-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated PR lifecycle test (Test 1)"
   - Extract and save the PR ID for subsequent steps

3. **Read the Pull Request details**
   - Use the `read Bitbucket pull request` action with the PR ID from the previous step
   - Verify that:
     - The PR ID matches the one from the creation step
     - The PR title matches what was set in the creation step
     - The PR status is "OPEN"
     - The PR is not a draft

4. **Update the Pull Request**
   - Use the `update Bitbucket pull request` action to change:
     - title to "Updated PR Lifecycle Test {timestamp}"
     - description to "This PR has been updated as part of the PR lifecycle test (Test 1) {timestamp}"
   - Read the PR again to verify the changes were applied
   - Verify the title and description match what was set

5. **Approve the Pull Request**
   - Use the `approve Bitbucket pull request` action
   - Read the PR again to verify it shows as approved
   - Verify the "approved" status is true in participants list

6. **Merge the Pull Request**
   - Use the `merge Bitbucket pull request` action with:
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
   - Use the `create Bitbucket pull request` action with the following parameters:
     - title: "PR Tasks Test {timestamp}"
     - source_branch: "feature/pr-tasks-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated PR tasks test (Test 2)"
   - Extract and save the PR ID for subsequent steps

3. **Create multiple PR tasks**
   - Use the `create Bitbucket pull request task` action with:
     - pr_id: the PR ID from step 2
     - content: "Task 1: Verify integration test {timestamp}"
     - repo_owner and repo_name: same as previous steps
   - Create two more tasks with different content:
     - "Task 2: Review code changes {timestamp}"
     - "Task 3: Test functionality {timestamp}"
   - Extract and save at least one task ID for the next steps

4. **List PR tasks**
   - Use the `list Bitbucket pull request tasks` action for the PR
   - Verify that all three tasks created in step 3 appear in the list
   - Verify that total number of tasks corresponds to the number of tasks created in step 3

5. **Update PR tasks**
   - Use the `update Bitbucket pull request task` action to mark "Task 1" as "RESOLVED"
   - Use the `update Bitbucket pull request task` action to update the content of "Task 2" to "Task 2: Code review completed {timestamp}"
   - List the tasks again to verify:
     - "Task 1" is now marked as "RESOLVED"
     - "Task 2" content has been updated
     - "Task 3" remains unchanged

6. **Clean up**
   - Approve the PR using the `approve Bitbucket pull request` action
   - Merge the PR using the `merge Bitbucket pull request` action with:
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
   - Use the `create Bitbucket pull request` action with the following parameters:
     - title: "Multi-Account Test PR (Test 2) {timestamp}"
   - Extract and save the PR ID for subsequent steps

3. **Approve the Pull Request as bot user**
   - Use the `approve Bitbucket pull request` action with the following parameters:
     - account: "bot"
   - Read the PR again (with bot account) to verify approval status
   - Verify that the PR shows as approved and note the approver username

4. **Merge the Pull Request as default user**
   - Use the `merge Bitbucket pull request` action with the following parameters:
     - account: "user"
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
   - Use the `create Bitbucket pull request` action with the following parameters:
     - title: "Draft PR Test (Test 3) {timestamp}"
     - source_branch: "feature/draft-pr-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated draft PR integration test (Test 3)"
     - draft: true
   - Extract and save the PR ID for subsequent steps

3. **Read the Pull Request details**
   - Use the `read Bitbucket pull request` action with the PR ID from the previous step
   - Verify that:
     - The PR ID matches the one from the creation step
     - The PR title matches what was set in the creation step
     - The PR status is "OPEN"
     - The PR is marked as a draft (check the draft status field)

4. **Update the Pull Request Draft Status**
   - Use the `update Bitbucket pull request` action with the following parameters:
     - draft: false
   - Read the PR again to verify the draft status is now false
   - Use the `update Bitbucket pull request` action with the following parameters:
     - draft: true
   - Read the PR again to verify the draft status is now true
   - Use the `update Bitbucket pull request` action with the following parameters:
     - draft: false
   - Read the PR again to verify the draft status is now false

5. **Merge the Pull Request**
   - Use the `merge Bitbucket pull request` action with:
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
   - Use the `create Bitbucket pull request` action with the following parameters:
     - title: "PR Review Tools Test {timestamp}"
     - source_branch: "feature/pr-review-tools-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated PR review tools test (Test 5)"
   - Extract and save the PR ID for subsequent steps
   - Ensure the PR includes the two TypeScript files and the three text files.

3. **Check diffs**
   - Use the pull request diff statistics action to get a list of files that are changed in the PR. Check if all expected files are present.
   - Use the pull request diff action to get a diff for two files at once. Check if the diff is correct and expected.

4. **Check file contents and create inline comments**
   - For each file (example1.ts and example2.ts):
     - Use the get Bitbucket file content action. Make sure it succeeds. Ensure the content is correct and expected.
     - **IMPORTANT**: Using the **file content**, identify line numbers with update markers. Use `cat -n` to understand the line numbers.
     - Add inline comments for each marker in the file. Include used line number in the comment message for better visibility.

5. **List and verify PR comments and line numbers**
   - Use the `list Bitbucket pull request comments` action to list all comments for the PR.
   - For each comment on a TypeScript file:
     - Use the `get Bitbucket file content` action to fetch the full file content.
     - Write the fetched file content to a uniquely named file in the `tmp/` directory in the current workspace, including a timestamp in the filename (e.g., `tmp/example1-<timestamp>.ts`).
     - Use **exactly** the following one-liner script to determine the actual line numbers for all marker comments in the file:
       ```
       cat -n <filename>
       ```
     - Compare the line numbers found by this script with the line numbers in the PR comments to verify mapping.
   - This approach ensures that the verification uses the actual file content as fetched from Bitbucket and provides a reproducible, timestamped record of the verification process.


6. **Add a general (non-inline) comment to the PR**
   - Use the `add Bitbucket pull request comment` action with the following parameters:
     - comment_text: "General comment for PR review tools test {timestamp}"
   - Save the returned comment ID as `PARENT_COMMENT_ID` for the reply check.
   - Verify that the general comment appears in the PR's comment list and is not associated with any file or line number.

7. **Add a reply to the general comment and verify parent linkage**
   - Use the `add Bitbucket pull request comment` action with the following parameters:
     - comment_text: "Reply comment for PR review tools test {timestamp}"
     - parent_comment_id: `PARENT_COMMENT_ID`
   - Save the returned reply comment ID as `REPLY_COMMENT_ID`.
   - Use the `list Bitbucket pull request comments` action to list all comments for the PR.
   - Verify that the entry with ID `REPLY_COMMENT_ID` includes `parent.id` equal to `PARENT_COMMENT_ID`.
   - Verify that the reply content matches "Reply comment for PR review tools test {timestamp}".

8. **Resolve a PR comment and verify `resolved` in list JSON**
   - Pick a comment ID from step 4 or step 6 (inline thread root or general comment, as supported by Bitbucket for resolve).
   - Use the `resolve Bitbucket pull request comment` action with the following parameters:
     - pr_id: the PR ID from step 2
     - comment_id: the comment ID chosen for resolve
     - repo_owner: your workspace name
     - repo_name: your repository name
     - account: optional
   - Call `list Bitbucket pull request comments` again and verify the JSON payload: the matching comment entry includes `resolved: true` (or document Bitbucket’s behavior if the thread cannot be resolved for that comment type).

Update a report as per [instruction](#test-results-reporting).

## Test 6: PR Comments Pagination

This test verifies that PR comments can be retrieved with pagination parameters (`page`, `pagelen`) and that pagination metadata is correctly returned.

### Steps

1. **Setup test environment**
   - Create a new working branch from main: `feature/pr-comments-pagination-test-{timestamp}`
   - Create (or update) a test file `integration-tests/bitbucket/test-files/integration-test-file.txt` with new content (e.g., current time)
   - Commit and push the changes

2. **Create a Pull Request**
   - Use the `create Bitbucket pull request` action with the following parameters:
     - title: "PR Comments Pagination Test {timestamp}"
     - source_branch: "feature/pr-comments-pagination-test-{timestamp}"
     - target_branch: "main"
     - repo_owner: your workspace name
     - repo_name: your repository name
     - description: "This is an automated PR comments pagination test (Test 6)"
   - Extract and save the PR ID for subsequent steps

3. **Post 20+ comments**
   - Use the `add Bitbucket pull request comment` action to add at least 20 comments to the PR
   - Mix general comments and inline comments on the test file:
     - Add at least 15 general comments with content such as "Pagination test comment {n} {timestamp}"
     - Add at least 5 inline comments (specifying a file path and line number from the committed test file)
   - Keep track of all comment IDs returned by each call
   - Verify that all 20+ comment creation calls succeeded

4. **List comments with default pagelen**
   - Use the `list Bitbucket pull request comments` action with the following parameters:
     - pr_id: the PR ID from step 2
     - repo_owner: your workspace name
     - repo_name: your repository name
   - Verify that:
     - All 20+ comments are returned in a single response (default pagelen is 100)
     - The `size` field in the response is >= 20
     - The `pagelen` field in the response is 100
     - The `page` field in the response is 1

5. **List comments with small pagelen**
   - Use the `list Bitbucket pull request comments` action with the following parameters:
     - pr_id: the PR ID from step 2
     - repo_owner: your workspace name
     - repo_name: your repository name
     - pagelen: 5
   - Verify that:
     - Exactly 5 comments are returned
     - The `pagelen` field in the response is 5
     - The `page` field in the response is 1
     - The `next` field is present (indicating there are more pages)
     - The `size` field is >= 20

6. **List comments on page 2**
   - Use the `list Bitbucket pull request comments` action with the following parameters:
     - pr_id: the PR ID from step 2
     - repo_owner: your workspace name
     - repo_name: your repository name
     - pagelen: 5
     - page: 2
   - Verify that:
     - Exactly 5 comments are returned
     - The `page` field in the response is 2
     - The comment IDs on page 2 are different from those on page 1

7. **Iterate through all pages and collect all comment IDs**
   - Starting from page 1, use the following parameters while iterating through all pages:
     - pagelen: 5
   - Increment `page` until the `next` field is absent
   - Collect all comment IDs across all pages
   - Verify that:
     - The total number of collected comment IDs matches the `size` reported in the first paginated response
     - All original comment IDs posted in step 3 are present in the collected set
     - No duplicate comment IDs appear across pages

8. **Clean up**
   - Approve the PR using the `approve Bitbucket pull request` action
   - Merge the PR using the `merge Bitbucket pull request` action with:
     - merge_strategy: "squash"
     - close_source_branch: "true"
   - Make sure the main branch is checked out again
   - Pull the latest changes
   - Delete the working branch

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
5. Use Bitbucket MCP actions to create, approve, and merge the PR:
    - title: "Bitbucket MCP Integration Tests Results {timestamp}"
    - description: "This is an automated integration tests results (Test 1, Test 2, Test 3, Test 4, Test 5, Test 6) {timestamp}"
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

1. Confirm the ATLACP Bitbucket MCP tools are available **to sub-agents** (sub-agents need the same MCP access as the parent; if a sub-agent cannot call MCP, **do not** run the test in the orchestrator—report that limitation to the user instead).
2. **Strict rule:** **Every** test’s executable steps run **only** in a dedicated sub-agent. The orchestrator may verify prerequisites, merge reports, and summarize; it **must not** run git or `bitbucket_*` calls for a test’s procedure in the main thread.
3. **For each test** you are asked to run: **spawn a dedicated sub-agent** and pass it the Bitbucket repo filesystem path, `repo_owner` / `repo_name` (from `git remote show origin` if needed), and the test number. If running more than one test, **do not** spawn sub-agents in parallel—**wait** for one test’s sub-agent to finish before starting the next (same repo; see [Execution model](#execution-model-mandatory-no-exceptions)).
4. Follow each test’s steps **inside that test’s sub-agent** exactly as written.
5. Document the results as you progress (orchestrator merges sub-agent outputs into `tmp/integration-tests-*.md` as appropriate).
6. Report the results as per [Test Results Reporting](#test-results-reporting) at the end of the process.

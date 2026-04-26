# Bitbucket CLI (`bbmd`) Integration Testing

This document provides step-by-step instructions for testing the Bitbucket CLI integration (`bbmd`). These tests mirror the scenarios in [`bitbucket-mcp-integration-tests.md`](./bitbucket-mcp-integration-tests.md); that file is the **behavioral reference** for what each test must prove against Bitbucket. Differences are limited to **transport** (shell + `bbmd` instead of MCP tools) and **CLI capability gaps** called out explicitly.

**Important**: If not mentioned, run all tests from this file.

**Important**: The user may provide the file system path to the Bitbucket repository. The current workspace is NOT necessarily the Bitbucket repository. If not provided, try to see if `atlacp-integration-tests` exists in a parent folder of the current workspace. Once you have the repository, `git remote show origin` helps determine `repo_owner` and `repo_name`. Do not recursively search all parent folders; check as suggested above, and if not found, ask the user for the path.

It is expected that the prompt to start the test will have the following form:

```markdown
Given the following Bitbucket repository `<file system path to the repository>`, run the tests as per the instructions in the `bitbucket-cli-integration-tests.md` file.
```

### Execution model (mandatory, no exceptions)

- **Tests always run in sub-agents.** Every numbered test (Test 1, Test 2, …) **must** be executed **only** inside a **separate delegated sub-agent** (e.g. Cursor Task tool or equivalent). The orchestrator **must not** run that test’s git commands, `go run ./cmd/bbmd` invocations, or any step that mutates the integration repo or Bitbucket state.
- **Sequential execution across tests.** All tests share the **same** integration repository clone: they race if run **in parallel** (concurrent sub-agents or terminals) on `git checkout`, `main`, pushes, and shared paths. When running more than one test, the orchestrator **must** run them **one after another**—wait until the sub-agent for Test *N* has fully finished before spawning the sub-agent for the next test. (Parallelism is **not** allowed for the suite; only isolation per test via sub-agents.)
- **No waiver.** This rule applies if the user asks for one test, “just Test N”, “first test only”, or the full suite. Convenience, speed, or “single coherent thread” are **not** reasons to run test steps in the orchestrator.
- **Orchestrator-only work is allowed** (and encouraged where useful): planning order, spawning sub-agents, passing repo path / `repo_owner` / `repo_name` / test id, **verification that does not execute the test** (e.g. confirming sub-agents can run `go run ./cmd/bbmd` from the atlacp repo root, merging `tmp/integration-tests-*.md` from outputs, summarizing pass/fail). The orchestrator does **not** substitute for a sub-agent when executing a test.
- **Why sub-agents**: isolates failures, avoids mixing shell/CLI steps for different tests in one thread, and matches review expectations. This is **not** permission to run multiple tests in parallel on the same repo.
- **Orchestrator responsibilities**: confirm the CLI is runnable for sub-agents (see [Prerequisites](#prerequisites-done-by-the-user)); when running multiple tests, launch sub-agents **sequentially** (wait for each to complete before the next); collect outputs and produce the final report.
- **Sub-agent responsibilities**: perform **all** shell work in the provided repo path and **all** `go run ./cmd/bbmd` commands for **that test only**; write/update the workspace results file for that test’s steps; return a concise pass/fail summary and PR ids to the orchestrator.

## Prerequisites done by the user

Assume the following is already prepared:

1. Sub-agents run the CLI with **`go run ./cmd/bbmd …`** from this (**atlacp**) repository root (Go toolchain required). If that fails, fix the environment (Go install, module path, working directory); do not substitute a prebuilt `bbmd` binary.
2. **Atlassian account configuration** matches the MCP tests: at least two accounts — default **user** and secondary **bot** — in the accounts file used by `bbmd` (see `--atlassian-accounts-file` on the root command; defaults apply when omitted, same as normal `bbmd` usage).
3. **No MCP server is required** for these tests; sub-agents use **shell + `go run ./cmd/bbmd` only** for Bitbucket operations.

For Git/SSH issues point the user on the [SSH troubleshooting](./README.md#ssh-troubleshooting). Don't read or do anything about it yourself, report and halt.

## Working with the repository

Use shell commands to work with the repository and all files. Example:

```bash
cd "<file system path to the repository>"

git pull

# Create a new branch
git checkout -b "feature/integration-test-{timestamp}"

# Add a new file
touch "integration-tests/bitbucket/test-files/integration-test-file.txt"

# Use echo and cat to write and read the file
echo "Current time: $(date)" > "integration-tests/bitbucket/test-files/integration-test-file.txt"

# Commit the changes
git add "integration-tests/bitbucket/test-files/integration-test-file.txt"

git commit -m "Add integration test file"

# Push and set upstream
git push origin "feature/integration-test-{timestamp}" --set-upstream
```

**Important**: Use shell commands for the repository and files. Do not use filesystem MCP tools for the integration repo when it lives outside the current workspace.

## CLI command conventions

- Run from the **atlacp project root** only:

  ```bash
  go run ./cmd/bbmd [global flags] pr <subcommand> [flags]
  ```

- **Global flags** (when needed): `--atlassian-accounts-file`, `--env`, `--log-level`, `--logs-file`, `--noop` (dry-run wiring only — **do not** use `--noop` for real integration tests that must hit Bitbucket).
- **PR subcommands** share **`--repo-owner`**, **`--repo-name`**, **`--pr-id`**, and optional **`--account`** where applicable.
- **JSON output**: successful `bbmd pr` and `bbmd file` commands write **JSON** to stdout. Parse PR id, task id, and comment id from that output.

In the numbered tests below, examples start with **`bbmd`** for readability; prefix with **`go run ./cmd/bbmd`** from the atlacp project root (same shape as [CLI command conventions](#cli-command-conventions)).

## File Structure for Testing

Same layout as the MCP doc:

```
integration-tests/
├── bitbucket/
│   ├── test-files/
│   │   └── integration-test-file.txt (File to modify during tests, just write current time)
│   └── results-YYYYMMDD-HHMMSS.md (Single results file for all tests — optional; primary reporting is under workspace `tmp/`)
```

Use a timestamp format for any repo-local results file (e.g. `results-20230615-120000.md`).

## MCP → `bbmd` mapping (quick reference)

| MCP tool (conceptual) | `bbmd` command |
|----------------------|----------------|
| `bitbucket_create_pr` | `bbmd pr create` |
| `bitbucket_read_pr` | `bbmd pr read` |
| `bitbucket_update_pr` | `bbmd pr update` |
| `bitbucket_approve_pr` | `bbmd pr approve` |
| `bitbucket_merge_pr` | `bbmd pr merge` |
| `bitbucket_list_pr_tasks` | `bbmd pr list-tasks` |
| `bitbucket_create_pr_task` | `bbmd pr create-task` |
| `bitbucket_update_pr_task` | `bbmd pr update-task` |
| `bitbucket_get_pr_diffstat` | `bbmd pr diffstat` |
| `bitbucket_get_pr_diff` | `bbmd pr diff` |
| `bitbucket_get_file_content` | `bbmd file content` |
| `bitbucket_add_pr_comment` | `bbmd pr add-comment` |
| `bitbucket_list_pr_comments` | `bbmd pr list-comments` |
| `bitbucket_resolve_pr_comment` | `bbmd pr resolve-comment` |

### Known CLI gaps (parity notes)

- **`pr create`**: there is **no `--draft` flag** on create. To obtain a draft PR for Test 4, either create a PR then run **`bbmd pr update --draft`** (see Test 4), or treat draft-on-create as out of scope until a flag exists.
- **`pr merge`**: supports **`--strategy`** (`merge_commit`, `squash`, `fast_forward`). **Merge message** and **close source branch** are not exposed as flags; rely on Bitbucket defaults or perform branch cleanup in git (tests already delete the branch locally).

---

## Test 1: PR Creation and Updates

Verifies PR create, read, update, approve, and merge using a single Atlassian account (default).

### Steps

1. **Setup test environment**
   - Create a branch from `main`: `feature/pr-lifecycle-test-{timestamp}`
   - Create or update `integration-tests/bitbucket/test-files/integration-test-file.txt`
   - Commit, change, commit again (two commits), note hashes, push

2. **Create a Pull Request**

   ```bash
   bbmd pr create \
     --title "PR Lifecycle Test {timestamp}" \
     --source-branch "feature/pr-lifecycle-test-{timestamp}" \
     --target-branch "main" \
     --repo-owner "<workspace>" \
     --repo-name "<slug>" \
     --description "This is an automated PR lifecycle test (Test 1)"
   ```

   Save the printed PR id.

3. **Read the Pull Request**

   ```bash
   bbmd pr read --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID>
   ```

   Verify id, title, OPEN, not draft.

4. **Update the Pull Request**

   ```bash
   bbmd pr update --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --title "Updated PR Lifecycle Test {timestamp}" \
     --description "This PR has been updated as part of the PR lifecycle test (Test 1) {timestamp}"
   ```

   Read again with `bbmd pr read` and verify title and description.

5. **Approve**

   ```bash
   bbmd pr approve --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID>
   ```

   Read PR again; verify approved state in output.

6. **Merge**

   ```bash
   bbmd pr merge --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --strategy squash
   ```

   Read PR; verify **MERGED**.

7. **Clean up** — checkout `main`, pull, delete local feature branch, confirm squash history as in MCP doc.

Update the report per [Test Results Reporting](#test-results-reporting).

---

## Test 2: PR Tasks Management

### Steps

1. **Setup** — branch `feature/pr-tasks-test-{timestamp}`, touch test file, commit, push

2. **Create PR**

   ```bash
   bbmd pr create \
     --title "PR Tasks Test {timestamp}" \
     --source-branch "feature/pr-tasks-test-{timestamp}" \
     --target-branch "main" \
     --repo-owner "<workspace>" --repo-name "<slug>" \
     --description "This is an automated PR tasks test (Test 2)"
   ```

3. **Create three tasks**

   ```bash
   bbmd pr create-task --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --content "Task 1: Verify integration test {timestamp}"
   bbmd pr create-task --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --content "Task 2: Review code changes {timestamp}"
   bbmd pr create-task --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --content "Task 3: Test functionality {timestamp}"
   ```

   Note a **task id** from output for updates (map to **Task 1** / **Task 2** by content if needed).

4. **List tasks**

   ```bash
   bbmd pr list-tasks --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID>
   ```

   Verify all three appear.

5. **Update tasks** — mark Task 1 **RESOLVED**, update Task 2 content:

   ```bash
   bbmd pr update-task --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --task-id <ID1> --state RESOLVED
   bbmd pr update-task --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --task-id <ID2> --content "Task 2: Code review completed {timestamp}"
   ```

   List again; confirm states and content.

6. **Clean up** — `bbmd pr approve`, then `bbmd pr merge --strategy squash`, then git cleanup as in MCP doc.

Update the report per [Test Results Reporting](#test-results-reporting).

---

## Test 3: Multi-Account PR Workflow

### Steps

1. **Setup** — branch `feature/multi-account-test-{timestamp}`, update test file, commit, push

2. **Create PR** (default account — omit `--account`)

   ```bash
   bbmd pr create \
     --title "Multi-Account Test PR (Test 3) {timestamp}" \
     --source-branch "feature/multi-account-test-{timestamp}" \
     --target-branch "main" \
     --repo-owner "<workspace>" --repo-name "<slug>"
   ```

3. **Approve as `bot`**

   ```bash
   bbmd pr approve --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --account bot
   bbmd pr read --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --account bot
   ```

4. **Merge as `user`**

   ```bash
   bbmd pr merge --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --account user --strategy squash
   ```

5. **Clean up** — git steps as in MCP doc.

Update the report per [Test Results Reporting](#test-results-reporting).

---

## Test 4: Draft Pull Request Creation and Verification

### Steps

1. **Setup** — branch `feature/draft-pr-test-{timestamp}`, update test file, commit, push

2. **Create a draft PR**

   `bbmd pr create` does **not** accept `--draft`. Use **create + update**:

   ```bash
   bbmd pr create \
     --title "Draft PR Test (Test 4) {timestamp}" \
     --source-branch "feature/draft-pr-test-{timestamp}" \
     --target-branch "main" \
     --repo-owner "<workspace>" --repo-name "<slug>" \
     --description "This is an automated draft PR integration test (Test 4)"

   bbmd pr update --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --draft
   ```

   (Use `--draft` on `pr update` to set draft to **true**; your Cobra build treats `--draft` as boolean.)

3. **Read PR** — `bbmd pr read ...`; verify OPEN and draft.

4. **Toggle draft** — `bbmd pr update` only sends draft when the `--draft` flag **changed** (Cobra `Changed("draft")`). Match the MCP sequence (false → true → false) using explicit **`--draft=false`**, **`--draft=true`**, **`--draft=false`** before each `bbmd pr read`, as needed.

5. **Merge**

   ```bash
   bbmd pr merge --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --strategy fast_forward
   ```

6. **Clean up** — git steps as in MCP doc.

Update the report per [Test Results Reporting](#test-results-reporting).

---

## Test 5: Bitbucket PR Review Tools End-to-End

### Steps

1. **Setup** — branch `feature/pr-review-tools-test-{timestamp}`. Copy `example1.ts` and `example2.ts` from the **atlacp** repo path `doc/testing/` into e.g. `integration-tests/bitbucket/test-files/ts-examples-{timestamp}/`. Commit and push.

2. **Create PR** — `bbmd pr create` with title `"PR Review Tools Test {timestamp}"`, same branches/repo fields as MCP doc.

3. **Diffstat**

   ```bash
   bbmd pr diffstat --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID>
   ```

   Confirm expected files appear.

4. **Diff** — run once per file or full diff:

   ```bash
   bbmd pr diff --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --path "integration-tests/bitbucket/test-files/ts-examples-{timestamp}/example1.ts"
   bbmd pr diff --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --path "integration-tests/bitbucket/test-files/ts-examples-{timestamp}/example2.ts"
   ```

5. **File content + inline comments** — for each TS file, get content at the branch tip commit (`git rev-parse HEAD` in the repo):

   ```bash
   bbmd file content \
     --repo-owner "<workspace>" --repo-name "<slug>" \
     --commit "<commit_sha>" \
     --path "integration-tests/bitbucket/test-files/ts-examples-{timestamp}/example1.ts"
   ```

   Use `cat -n` locally to find **Update marker** lines; add inline comments:

   ```bash
   bbmd pr add-comment --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --content "inline marker line <N> {timestamp}" \
     --file-path "integration-tests/bitbucket/test-files/ts-examples-{timestamp}/example1.ts" \
     --line-from <N> --line-to <N>
   ```

6. **List comments** — `bbmd pr list-comments ...`; for each TS comment, re-fetch file with `bbmd file content`, save copy under workspace `tmp/example1-<timestamp>.ts`, run `cat -n` and compare line numbers.

7. **General comment**

   ```bash
   bbmd pr add-comment --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --content "General comment for PR review tools test {timestamp}"
   ```

   Save the returned comment ID as `PARENT_COMMENT_ID`.

8. **Reply comment** — create a reply to the general comment and verify the parent linkage:

   ```bash
   bbmd pr add-comment --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> \
     --content "Reply comment for PR review tools test {timestamp}" \
     --parent-comment-id <PARENT_COMMENT_ID>
   ```

   Save the returned comment ID as `REPLY_COMMENT_ID`, then list comments:

   ```bash
   bbmd pr list-comments --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID>
   ```

   Verify the entry with ID `REPLY_COMMENT_ID` includes `parent.id` equal to `PARENT_COMMENT_ID` and the expected reply content.

9. **Resolve** — pick a comment id, then:

   ```bash
   bbmd pr resolve-comment --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --comment-id <ID>
   bbmd pr list-comments --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --include-resolved
   ```

   Verify `resolved` in JSON/text output per your build.

Update the report per [Test Results Reporting](#test-results-reporting).

---

## Test 6: PR Comments Pagination

This test matches the MCP doc’s pagination scenarios (`page`, `pagelen`, full iteration). Use **`bbmd pr list-comments`** with optional **`--page`** and **`--pagelen`** (same semantics as the MCP tool: omit both for defaults; default page size is 100 at the app layer).

### Steps

1. **Setup test environment** — same as [MCP Test 6](./bitbucket-mcp-integration-tests.md#test-6-pr-comments-pagination) step 1: branch `feature/pr-comments-pagination-test-{timestamp}`, test file update, commit, push.

2. **Create a Pull Request**

   ```bash
   bbmd pr create \
     --title "PR Comments Pagination Test {timestamp}" \
     --source-branch "feature/pr-comments-pagination-test-{timestamp}" \
     --target-branch "main" \
     --repo-owner "<workspace>" \
     --repo-name "<slug>" \
     --description "This is an automated PR comments pagination test (Test 6)"
   ```

   Save the PR id.

3. **Post 20+ comments** — use `bbmd pr add-comment` repeatedly (general and inline), same counts as the MCP doc (15+ general, 5+ inline). Track comment IDs from JSON output.

4. **List with default pagination** — no `--page` or `--pagelen`:

   ```bash
   bbmd pr list-comments --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID>
   ```

   Verify JSON: `size` ≥ 20, `pagelen` is 100, `page` is 1, and all comments appear on one page when under the default limit.

5. **Small page size** — first page only:

   ```bash
   bbmd pr list-comments --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --pagelen 5
   ```

   Verify: five items in `values`, `pagelen` 5, `page` 1, `next` present, `size` ≥ 20.

6. **Second page**:

   ```bash
   bbmd pr list-comments --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --pagelen 5 --page 2
   ```

   Verify: five items, `page` 2, comment IDs differ from step 5’s first page.

7. **Iterate all pages** — repeat with `--pagelen 5` and `--page` 1, 2, … until `next` is absent (or use a short shell loop). Collect all comment IDs; total must match `size` from the first paginated response; set must equal IDs from step 3 with no duplicates.

8. **Clean up** — `bbmd pr approve`, `bbmd pr merge --strategy squash`, checkout `main`, pull, delete branch (same as other tests).

Update the report per [Test Results Reporting](#test-results-reporting).

---

## Test Results Reporting

Follow the protocol below when performing the test:

1. In the **current workspace** create (if needed) `tmp/integration-tests-{YYYYMMDD-HHMMSS}-results.md`
2. For each step, append step number, description, **PASS/FAIL**; on failure, note why
3. **Do not stop** on failure — document and continue

### Format of the results file

```markdown
# Bitbucket CLI Integration Test Results
Test executed at: {timestamp}

## Test 1: <Title of the test>
- Step 1: <Step description> - PASS
- Step 2: <Step 2 description> (Pull Request (PR #{pr_id})) - PASS
...

## Summary
- All tests: PASS/FAIL
- Issues encountered: None/List issues
```

### Publishing results to the integration-tests repo (optional)

When all tests complete, you may copy the workspace results file into the integration-tests repository (same flow as the MCP doc):

1. Checkout `main` and pull
2. Branch `feature/bitbucket-cli-integration-tests-results-{timestamp}`
3. Copy workspace results to `integration-tests/bitbucket/results-{timestamp}.md`
4. Commit and push
5. Open and merge a PR using **`bbmd`**:

   ```bash
   bbmd pr create \
     --title "Bitbucket CLI Integration Tests Results {timestamp}" \
     --source-branch "feature/bitbucket-cli-integration-tests-results-{timestamp}" \
     --target-branch "main" \
     --repo-owner "<workspace>" --repo-name "<slug>" \
     --description "Automated CLI integration test results (Tests 1–5; Test 6 N/A) {timestamp}"

   bbmd pr approve --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID>
   bbmd pr merge --repo-owner "<workspace>" --repo-name "<slug>" --pr-id <PR_ID> --strategy squash
   ```

6. Share summary with the user
7. Optionally attach the Bitbucket results URL to a GitHub PR in this repo (if applicable), same pattern as the MCP doc’s `gh pr comment` section.

---

## Automation Instructions for AI Model

When asked to run integration tests using this document:

1. Confirm sub-agents can run **`go run ./cmd/bbmd`** from the atlacp project root. If not, **do not** run the test in the orchestrator — report the environment problem to the user so they can fix Go/repo setup.
2. **Strict rule:** every test’s executable steps run **only** in a dedicated sub-agent. The orchestrator may verify prerequisites, merge reports, and summarize; it **must not** run git or `go run ./cmd/bbmd` for a test’s procedure in the main thread.
3. **For each test:** spawn a dedicated sub-agent with the Bitbucket repo path, `repo_owner` / `repo_name`, and test number. If running more than one test, **do not** spawn sub-agents in parallel—**wait** for one test’s sub-agent to finish before starting the next (same repo; see [Execution model](#execution-model-mandatory-no-exceptions)).
4. Follow each test’s steps **inside that sub-agent** exactly as written.
5. Merge sub-agent outputs into `tmp/integration-tests-*.md` as appropriate.
6. Report per [Test Results Reporting](#test-results-reporting).

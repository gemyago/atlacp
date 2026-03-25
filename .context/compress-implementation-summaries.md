# Instruction to compress implementation summaries

Your job is to produce a single compressed `implementation-summary.md` from per-task summary files, then remove the originals.

**Pre Condition**: This instruction is for AI agents with mode switching/running sub-agents capabilities only. If you don't have access to sub-agents, report the limitation and do not proceed.

## Input

You will be given a path to a plan, e.g. `docs/implementation/<plan-slug>/plan-<plan-slug>.md`. Its folder contains:

- `plan-<plan-slug>.md` — the plan document
- `summary-task-X.Y.md` — one file per completed task (e.g. summary-task-1.1.md, summary-task-1.2.md)

## Process

1. List all `summary-task-*.md` files in the folder (in task order: 1.1, 1.2, 2.1, …).
2. For **each** summary file, start a sub-agent (you can start multiple in parallel if you can) with:
   **Do not extract yourself.** Sub-agents are required; reading files and extracting in the main agent is not allowed.
   <sub-agent-instruction>
   Read the file at <full path to summary file>. Extract task ID, title, summary (1–2 sentences), and anything **unusual** that is not in the plan or PR: implementation deviations, unexpected issues, gotchas, workarounds, decisions made differently. Return **only** this block, no other text:

   ---
   TASK: X.Y
   TITLE: <task title from heading>
   SUMMARY: <1–2 sentences>
   DEVIATIONS:
   - <deviation, issue, or notable decision — omit section if none>
   ---

   Rules: Only include DEVIATIONS if there is something noteworthy. Skip boilerplate (Completion Protocol, Plan reference, file lists, test lists — PR and plan already have those). Do not create or modify files.
   </sub-agent-instruction>
   Note: Include full path to the summary file. Sub-agent returns the extraction block only.
3. Collect all extraction blocks from sub-agents.
4. Compile them into a single `implementation-summary.md` (see Output format below).
5. **Delete** each original `summary-task-*.md` file after the compressed file is written.

## Output format

Write `implementation-summary.md` with this structure:

```markdown
# Implementation Summary: <plan title>

**Plan:** [plan-<plan-slug>.md](./plan-<plan-slug>.md)

## Overview

<2–4 sentences summarizing what was implemented across all tasks>

## Tasks

### Task X.Y: <title>
<1–2 sentence summary. Add deviations/notes if any.>

### Task X.Y: <title>
...

## Deviations & notes

<Consolidated list of implementation deviations, unexpected issues, or notable decisions — include only if any were extracted>

## Completion

- Lint: ✓
- Type check: ✓
- Tests: ✓
```

## Rules

- **Sub-agents required**: Step 2 must use sub-agents for extraction. Do not read summary files and extract in the main agent.
- **Deduplicate**: Plan reference once. Completion protocol once. No repeated boilerplate.
- **Compact**: Use bullets and tables. Avoid verbose per-file prose.
- **Preserve essentials**: Summary per task, deviations/notes, completion status. Skip file lists and test lists (PR and plan have them).
- **Read-only for code**: Do not modify any source files. Only create `implementation-summary.md` and delete `summary-task-*.md` files in the plan folder.
- **Delete originals**: After writing `implementation-summary.md`, delete each `summary-task-*.md` file to remove noise.

## Success

Report back: `Compressed implementation summary written to docs/implementation/<plan-slug>/implementation-summary.md. Original summary files removed.`

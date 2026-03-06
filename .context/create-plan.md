# Instruction to create plan of the work

## Goal

To guide an AI assistant in creating a detailed plan of the work in Markdown format, based on an initial user prompt. The plan should be clear, actionable, and suitable for a junior developer to understand and implement the required change.

## Process

1.  **Receive Initial Prompt:** The user provides a brief description or request for a new feature or functionality, optionally referencing existing document
2.  **Do Research:** Before creating the plan, the AI *must* do the research of the codebase to gather sufficient detail. The goal is to understand the "what", "why" and most important **how** of the requested work.
3.  **Generate Plan:** Based on the initial prompt and the research, generate the plan using the structure outlined below.
4.  **Save the Plan:** Save the generated document as `plan-[feature-name].md` inside the `/doc/[feature-name]` directory.

## Research Areas (Examples)

The AI should adapt its research areas based on the prompt, but here are some common questions to explore:

*   **Problem/Goal:** "What problem does this feature solve?" or "What is the main goal we want to achieve with this change?"
*   **Core Functionality:** "What is a core functionality that this work should enable?"
*   **Acceptance Criteria:** "How will we know when this feature is successfully implemented? What are the key success criteria?"
*   **Scope/Boundaries:** "Are there any specific things this feature *should not* do (non-goals)?"
*   **Data Requirements:** "What kind of data does this feature need to display or manipulate?"
*   **Design/UI:** "Are there any existing design mockups or UI guidelines to follow?" or "Can you describe the desired look and feel?"
*   **Edge Cases:** "Are there any potential edge cases or error conditions we should consider?"

AI should do it's best to understand what needs to be built. Any small uncertanties should be listed in the resulting plan document. AI should **only** ask clarifying questions if the initial prompt is hightly ambiguous and the AI has failed to understand the final outcome.

## Plan Structure

The generated Plan may include the following sections when applicable:
1. **Introduction/Overview:** Briefly describe the feature and the problem it solves. State the goal.
2. **Business Logic:** Describe (in words) main aspects of the business logic that will be implemented.
3. **High Level Architecture:** Describe the high level architecture of the feature, list components involved.
4. **Detailed Architecture:** For each component involved, describe how it will work and structured, which files may need to be created or updated
5. **Key Architectural Decisions:** List key architectural decisions that were made.
6. **Uncertanties:** List any uncertanties (if present) or areas needing further clarification.
8. **Releted Files** List all files related to the change. If new files needs to be created - mention them as well.
7. **Task List** Detailed numbered list of tasks (steps) that needs to be taken to implement the required plan. Please note that TDD approach should be followed to implement the desired change as per [tdd-flow](../tdd-flow.md) (mention this in the task list). Please note that each task should be self contained - codebase is "operational" after each task. Completion protocol must followed as per project guidelines to ensure the codebase is in a "operational" state.

### Example task format

```markdown
**Task X.X: Implement user update handling**
- Add new `UserUpdate` struct
- Add stub function `updateUser(update UserUpdate)` in UserCommands
- Write failing tests (usually happy path and few edge cases)
  - expected user is written to the database
  - invalid user ID results in error
  - conflicting user email results in error
- Run affected tests: `go test -v ./<package> --run <test pattern>`
  - Verify failure is expectation (e.g expected not to equal, exists e.t.c).
  - Compilation errors are **not acceptable** - missing stubs should be added, test should be retried.
- Implement `updateUser(_ update: UserUpdate)` logic
- Run affected tests: `go test -v ./<package> --run <test pattern>`
  - Verify all tests pass
- Write summary to `doc/<plan-slug>/completed-task-x.x.md`
- Success criteria: As per completion protocol, at least: `make test` passes, `make lint` passes, summary written
```

Important notes:
- Tests for new types/new fields are not required, only logic needs tests
- Any task should leave the codebase in "operational" state

## Target Audience

Assume the primary reader of the Plan is a **junior developer**. Therefore, requirements should be explicit, unambiguous, and avoid jargon where possible. Provide enough detail for them to understand the feature's purpose and core logic.

## Output

*   **Format:** Markdown (`.md`)
*   **Location:** `/doc/[feature-name]`
*   **Filename:** `plan-[feature-name].md`

## Final instructions

1. Do NOT start implementing the Plan
2. Do the research as stated
3. Avoid asking questions unless there is a **very** strong reason

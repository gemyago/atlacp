# Instruction to orchestrate plan implementation

**Pre Condition**: This instruction is for AI agents with mode switching/running sub-agents capabilities only. If you don't have access to mode switching/running sub-agents, please report the limitation and do not proceed.

## Input

You will be given a plan usually in form of a document as an input. The plan will usually include a tasks list. Each task should be assumed to be `atomic` - this means it can be implemented fully independently and after the implementation the codebase is supposed to stay in a "green" state (all tests passing).

## The orchestration process

As a first step please identify `atomic` tasks and build a TODO for for yourself to orchestrate implementation.

For each task start a sub-agent with the exact instruction enclosed in `<sub-agent-instruction>` tags. **Note:** do not include enclosing tags, just contents. Do not include any additional text or context:
1. Start a coding agent with implementation instruction as follows:
  <sub-agent-instruction>
  Follow [.context/implement-plan-task.md](.context/implement-plan-task.md) to implement the following task: Task XX: <task description> from <plan reference>
  </sub-agent-instruction>
2. Once results received, finalize the task by starting debugging subagent as follows:
  <sub-agent-instruction>
  Please verify the implementation of Task XX: <task description> by following [.context/finalize-plan-task.md](.context/finalize-plan-task.md)
  </sub-agent-instruction>
3. If the task finalization fails, start a coding subagent to fix the codebase as follows:
  <sub-agent-instruction>
  Please fix the codebase to make it "green" by following [.context/fix-broken-codebase.md](.context/fix-broken-codebase.md)
  </sub-agent-instruction>
4. Repeat step 2 to finalize the task again.

Note: user may intervene (and cancel) any sub-agent flow if got stuck or otherwise got the wrong way. If you identify that user intervention took place, ask the user what to do next, don't proceed.

Once verification succeeded, report the progress to the user as follows:
- Step 1: succeeded
- Step 2: succeeded
- Step 3: succeeded/skipped (step 2 succeeded)
- Step 4: succeeded/skipped (step 2 succeeded)
- Summary: implementation completed!

Proceed with orchestration of a next task using the same process.

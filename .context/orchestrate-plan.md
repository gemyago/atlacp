# Instruction to orchestrate plan implementation

**Pre Condition**: This instruction is for AI agents with mode switching/running sub-agents capabilities only. If you don't have access to mode switching/running sub-agents, please report the limitation and do not proceed.

**IMPORTANT** AI must orchestrate the implementation of the plan, not implement the plan itself, this is a job of sub-agents. AI must only read relevant files (usually just plan) and then start sub-agents, that's it!

## Input

You will be given a plan usually in form of a document as an input. The plan will include a tasks list. Each task should be assumed to be `atomic` - this means it can be implemented fully independently and after the implementation the codebase is supposed to stay in a "green" state (all tests passing).

AI must create TODO list for each such task. Then AI must follow the exact same orchestration process (see next) for each task in the TODO list. AI must not attempt to optimize or the plan or deviate from the plan in any way. AI must focus on the orchestration and nothing else.

## The orchestration process

As a first step please identify `atomic` tasks and build a TODO for for yourself to orchestrate implementation.

For each task start a sub-agent with the exact instruction enclosed in `<sub-agent-instruction>` tags. **Note:** do not include enclosing tags, just contents. Do not include any additional text or context:
1. Start a coding agent with implementation instruction as follows:
  <sub-agent-instruction>
  Follow [.context/implement-plan-task.md](.context/implement-plan-task.md) to implement the following task: Task XX: <task description> from <plan reference>
  </sub-agent-instruction>
  Note: Include a full reference to the `implement-plan-task.md` file in the instruction, do not cite it's contents. Sub-agent will read it.
2. Once results received, finalize the task by starting debugging subagent as follows:
  <sub-agent-instruction>
  Please verify the implementation of Task XX: <task description> by following [.context/finalize-plan-task.md](.context/finalize-plan-task.md)
  </sub-agent-instruction>
  Note: Include a full reference to the `finalize-plan-task.md` file in the instruction, do not cite it's contents. Sub-agent will read it.
3. If the task finalization fails, start a coding subagent to fix the codebase as follows:
  <sub-agent-instruction>
  Please fix the codebase to make it "green" by following [.context/fix-broken-codebase.md](.context/fix-broken-codebase.md)
  </sub-agent-instruction>
  Note: Include a full reference to the `fix-broken-codebase.md` file in the instruction, do not cite it's contents. Sub-agent will read it.
4. Repeat step 2 to finalize the task again.

Note: user may intervene (and cancel) any sub-agent flow if got stuck or otherwise got the wrong way. If you identify that user intervention took place, ask the user what to do next, don't proceed.

Once implementation of the task is complete, report back success and proceed with orchestration of a next task using the same process.
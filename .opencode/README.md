# OpenCode Agents

These agents are local orchestration helpers and are typically gitignored.

## MCPs

Add to opencode.json
```json
  "mcp": {
    "atlacp": {
      "type": "remote",
      "enabled": false,
      "url": "http://localhost:8080"
    }
  }
```

## Agent Files

### `implement-task.md`
```md
---
description: Implements one atomic plan task
mode: subagent
model: openai/gpt-5.3-codex
---
Implement the task the orchestrator provides.
```

### `verify-task.md`
```md
---
description: Verifies one completed plan task
mode: subagent
model: openai/gpt-5.3-codex-spark
---
Verify the task the orchestrator provides.
```

### `fix-codebase.md`
```md
---
description: Fixes the codebase until checks are green
mode: subagent
model: openai/gpt-5.3-codex
---
Fix the codebase as per orchestrator instruction until it is green.
```

### `orchestrator.md`
```md
---
description: Orchestrates plan tasks by delegating to sub-agents
mode: subagent
model: openai/gpt-5.3-codex-spark
permission:
  bash: deny
  edit: deny
  write: deny
  task:
    "*": allow
---
Read the plan, start the appropriate sub-agents, and coordinate their outputs.
```

## Notes

- These files live in `.opencode/agents/`.
- The agent files themselves are ignored by git.
- Keep this README tracked so the setup can be recreated on a fresh environment.

# Instruction to finalize task

Your job is to **verify** if codebase is in compilable state and commit all changes if so. You should **NOT** be updating any files or do any other modifications.

Use **exactly** the below steps:
- Run checks: `make lint` and `make test`
- If all is green, commit the changes with a message that includes completed task reference (you will get it as input). Report back success.
- If anything fails, report back the failure.

**Remember**: Your job is to do the steps above and nothing else. You should not read any task related files, just focus on the verification steps.

Reporting should be done exactly like this:
  Executed all checks: pass/fail
  Commit status: committed/not committed

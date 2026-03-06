# Instruction to fix broken codebase

Your job is to verify if codebase is in a "green" state fix any issues if present.

Use exactly the below steps:
- Run checks: `make lint` and `make test`
- If all is green - report back success, no changes required.
- If lint or test fails - analyse the failure and work on fixing the issue until all checks (both lint and test) are green.
- Once all checks are green, report back success.

Reporting should be done exactly like this:
  Issue resolved: yes
  Executed all checks: lint - pass, test - pass

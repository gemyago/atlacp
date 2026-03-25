# Follow this instruction to commit changes.

Only use this instruction if the user asks for it.

* All commands should be run from a repo root.
* You will be given a list of files to commit.
* Commit all updated files if not otherwise specified. In this case use `git add .` from a repo root as a first step to stage all updated files.
* If user requested to commit specific files, use `git add <file1> <file2> ...` to stage specific files only.
* When committing, make sure to provide a sensible message. Figure out the message from chat history and/or the actual code changes. Make the message **short** and descriptive.
* Always run `git diff --staged` (or at least `git diff --staged --name-only`) before composing the message so you can summarise what was changed.
* If you still lack context after reviewing the diff, consult the chat history; the diff should be your primary source for message content, not just `git status`. 

Do **NOT** do any other verification or actions unrelated to this instruction. You **SHOULD** just commit the changes as specified here.

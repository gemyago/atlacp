# Instruction to create pull request

This is an instruction to follow when user is referencing it. Only use this instruction when explicitly requested by the user.

1. Look on a commit history between a base (if not mentioned otherwise, use **main**). You will need to run command like below
```bash
# Make sure remote is up to date
git fetch origin

# to get current branch, you will use it in step 5
git branch

# git log
git log origin/<base branch>...HEAD --oneline | cat
```

2. Review commit history (git log output above) and come up with a sensible PR title and description.
3. The description should summarize meaningful changes and follow format similar to below:
  ```md
  # Focus of the changes
  <Short summary of the changes (1-2 sentences), what they are about and why they are needed.>

  # High level summary of the changes
  * Short change description 1
  * Short change description 2
  * ...
  ```
  Note: Use your best judgement and feel free to adjust the format if changes are more complex and require more detailed description. The main point is to make it easy for reviewers (usually humans, and humans to not like to read poems) to understand what changes are about and why they are needed.
4. Prepare PR title and description in the following format:
  ```md
  **PR title**:
  <PR title>
  ---
  **PR description**:
  <PR description>
  ```
5. Push pending changes and create a PR with a command below:
```bash
git push origin <current branch> --set-upstream

gh pr create --title "<PR title>" --base <base branch> --head <current branch> --body "$(cat <<EOF
# Focus of the changes
<Short summary of the changes (1-2 sentences), what they are about and why they are needed.>

# High level summary of the changes
* Short change description 1
* Short change description 2
* ...
EOF
)"
```
Note: Using cat and EOF allows you to create a multi-line PR description without worrying about escaping special characters. Make sure to use actual PR title and description and other parameters as needed.

6. Show the PR to the user as a URL so user can click it, as well as full URL for copying.

# Bitbucket API Client Endpoints

POST /repositories/{username}/{repo_slug}/pullrequests
Client method: CreatePR(ctx, tokenProvider, CreatePRParams)

GET /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}
Client method: GetPR(ctx, tokenProvider, GetPRParams)

PUT /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}
Client method: UpdatePR(ctx, tokenProvider, UpdatePRParams)

POST /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/approve
Client method: ApprovePR(ctx, tokenProvider, ApprovePRParams)

POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/request-changes
Client method: RequestPRChanges(ctx, tokenProvider, RequestPRChangesParams)

POST /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/merge
Client method: MergePR(ctx, tokenProvider, MergePRParams)

GET /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/tasks
Client method: ListPullRequestTasks(ctx, tokenProvider, ListPullRequestTasksParams)

GET /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/tasks/{task_id}
Client method: GetTask(ctx, tokenProvider, GetTaskParams)

POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/tasks
Client method: CreatePullRequestTask(ctx, tokenProvider, CreatePullRequestTaskParams)

PUT /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/tasks/{task_id}
Client method: UpdateTask(ctx, tokenProvider, UpdateTaskParams)

GET /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments
Client method: ListPRComments(ctx, tokenProvider, ListPRCommentsParams) 
# Bitbucket API Client Endpoints

POST /repositories/{username}/{repo_slug}/pullrequests
Client method: CreatePR(ctx, tokenProvider, CreatePRParams)

GET /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}
Client method: GetPR(ctx, tokenProvider, GetPRParams)

PUT /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}
Client method: UpdatePR(ctx, tokenProvider, UpdatePRParams)

POST /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/approve
Client method: ApprovePR(ctx, tokenProvider, ApprovePRParams)

POST /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/merge
Client method: MergePR(ctx, tokenProvider, MergePRParams)

GET /repositories/{username}/{repo_slug}/pullrequests/{pull_request_id}/tasks
Client method: ListPullRequestTasks(ctx, tokenProvider, ListPullRequestTasksParams) 
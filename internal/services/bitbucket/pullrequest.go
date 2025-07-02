package bitbucket

import "time"

// PullRequestBranch represents a branch in a pull request.
type PullRequestBranch struct {
	Name string `json:"name"`
}

// PullRequestCommit represents a commit in a pull request.
type PullRequestCommit struct {
	Hash string `json:"hash"`
}

// PullRequestRepository represents repository information in a pull request.
type PullRequestRepository struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
}

// PullRequestSource represents the source of a pull request.
type PullRequestSource struct {
	Branch     PullRequestBranch      `json:"branch"`
	Commit     *PullRequestCommit     `json:"commit,omitempty"`
	Repository *PullRequestRepository `json:"repository,omitempty"`
}

// PullRequestDestination represents the destination of a pull request.
type PullRequestDestination struct {
	Branch     PullRequestBranch      `json:"branch"`
	Commit     *PullRequestCommit     `json:"commit,omitempty"`
	Repository *PullRequestRepository `json:"repository,omitempty"`
}

// PullRequestAuthor represents the author of a pull request.
type PullRequestAuthor struct {
	AccountID   string `json:"account_id"`
	DisplayName string `json:"display_name"`
	Nickname    string `json:"nickname"`
	Username    string `json:"username"`
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
}

// PullRequestSummary represents the summary/description of a pull request.
type PullRequestSummary struct {
	Raw    string `json:"raw"`
	Markup string `json:"markup"`
	HTML   string `json:"html"`
	Type   string `json:"type"`
}

// PullRequest represents a Bitbucket pull request.
type PullRequest struct {
	ID int `json:"id"`

	// Title is the title of the pull request. Omitempty is required for partial updates.
	Title string `json:"title,omitempty"`

	Description string `json:"description,omitempty"`

	// State is the current state of the pull request. Omitempty prevents sending empty state during updates.
	State string `json:"state,omitempty"`

	Author *PullRequestAuthor `json:"author,omitempty"`

	// Source contains branch information. The omitzero prevents "branch not found" errors during partial updates.
	Source PullRequestSource `json:"source,omitzero"`

	Destination       *PullRequestDestination `json:"destination,omitempty"`
	Reviewers         []PullRequestAuthor     `json:"reviewers,omitempty"`
	Participants      []Participant           `json:"participants,omitempty"`
	CloseSourceBranch bool                    `json:"close_source_branch,omitempty"`
	Summary           *PullRequestSummary     `json:"summary,omitempty"`
	CommentCount      int                     `json:"comment_count,omitempty"`
	TaskCount         int                     `json:"task_count,omitempty"`
	Type              string                  `json:"type,omitempty"`
	CreatedOn         *time.Time              `json:"created_on,omitempty"`
	UpdatedOn         *time.Time              `json:"updated_on,omitempty"`
	MergeCommit       *PullRequestCommit      `json:"merge_commit,omitempty"`
	Draft             *bool                   `json:"draft,omitempty"`
}

// Participant represents a pull request participant (for approval responses).
type Participant struct {
	User     PullRequestAuthor `json:"user"`
	Role     string            `json:"role"`
	Approved bool              `json:"approved"`
	State    string            `json:"state,omitempty"`
	Type     string            `json:"type"`
}

// PullRequestMergeParameters represents parameters for merging a pull request.
type PullRequestMergeParameters struct {
	Type              string `json:"type,omitempty"`
	Message           string `json:"message,omitempty"`
	CloseSourceBranch bool   `json:"close_source_branch,omitempty"`
	MergeStrategy     string `json:"merge_strategy,omitempty"`
}

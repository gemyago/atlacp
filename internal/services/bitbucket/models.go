package bitbucket

import (
	"time"
)

// Link represents a link to a resource related to an object.
type Link struct {
	Href string `json:"href"`
	Name string `json:"name,omitempty"`
}

// Links contains links to related resources.
type Links struct {
	Self   *Link `json:"self,omitempty"`
	HTML   *Link `json:"html,omitempty"`
	Avatar *Link `json:"avatar,omitempty"`
}

// Account represents a Bitbucket account.
type Account struct {
	Type        string `json:"type"`
	Links       *Links `json:"links,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	AccountID   string `json:"account_id,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
}

// TaskContent represents the content of a task.
type TaskContent struct {
	Raw    string `json:"raw"`
	Markup string `json:"markup,omitempty"`
	HTML   string `json:"html,omitempty"`
}

// Task represents a Bitbucket task.
type Task struct {
	ID         int64        `json:"id"`
	CreatedOn  time.Time    `json:"created_on"`
	UpdatedOn  time.Time    `json:"updated_on"`
	State      string       `json:"state"`
	Content    *TaskContent `json:"content"`
	Creator    *Account     `json:"creator"`
	Pending    bool         `json:"pending,omitempty"`
	ResolvedOn time.Time    `json:"resolved_on,omitempty"`
	ResolvedBy *Account     `json:"resolved_by,omitempty"`
}

// PullRequestTask represents a task on a pull request.
type PullRequestTask struct {
	Task
	Links *Links `json:"links,omitempty"`
}

// Comment represents a comment on Bitbucket.
type Comment struct {
	ID        int64        `json:"id,omitempty"`
	CreatedOn time.Time    `json:"created_on,omitempty"`
	UpdatedOn time.Time    `json:"updated_on,omitempty"`
	Content   *TaskContent `json:"content,omitempty"`
	User      *Account     `json:"user,omitempty"`
	Pending   bool         `json:"pending,omitempty"`
}

// PullRequestCommentTask represents a task related to a comment on a pull request.
type PullRequestCommentTask struct {
	PullRequestTask
	Comment *Comment `json:"comment,omitempty"`
}

// PaginatedTasks represents a paginated list of tasks.
type PaginatedTasks struct {
	Size     int                      `json:"size,omitempty"`
	Page     int                      `json:"page,omitempty"`
	PageLen  int                      `json:"pagelen,omitempty"`
	Next     string                   `json:"next,omitempty"`
	Previous string                   `json:"previous,omitempty"`
	Values   []PullRequestCommentTask `json:"values"`
}

// DiffStatPath handles Bitbucket's "old"/"new" fields which may be a string or object.
// CommitFile matches the Bitbucket OpenAPI "commit_file" definition.
type CommitFile struct {
	Type        string  `json:"type"`
	Path        string  `json:"path,omitempty"`
	Commit      *Commit `json:"commit,omitempty"`
	Attributes  string  `json:"attributes,omitempty"`
	EscapedPath string  `json:"escaped_path,omitempty"`
}

// Commit matches the Bitbucket OpenAPI "commit" definition.
type Commit struct {
	Hash         string         `json:"hash,omitempty"`
	Date         string         `json:"date,omitempty"`
	Author       interface{}    `json:"author,omitempty"`    // Could be expanded if needed
	Committer    interface{}    `json:"committer,omitempty"` // Could be expanded if needed
	Message      string         `json:"message,omitempty"`
	Summary      *CommitSummary `json:"summary,omitempty"`
	Parents      []*Commit      `json:"parents,omitempty"`
	Repository   interface{}    `json:"repository,omitempty"`   // Could be expanded if needed
	Participants interface{}    `json:"participants,omitempty"` // Could be expanded if needed
}

// CommitSummary matches the summary object in the commit schema.
type CommitSummary struct {
	Raw    string `json:"raw,omitempty"`
	Markup string `json:"markup,omitempty"`
	HTML   string `json:"html,omitempty"`
}

// UnmarshalJSON supports both string and object with "path" field.

// DiffStat represents a summary of changes made to a file between two commits.
type DiffStat struct {
	Type         string      `json:"type,omitempty"`
	Status       string      `json:"status,omitempty"`
	LinesAdded   int         `json:"lines_added,omitempty"`
	LinesRemoved int         `json:"lines_removed,omitempty"`
	Old          *CommitFile `json:"old,omitempty"`
	New          *CommitFile `json:"new,omitempty"`
	Path         string      `json:"path,omitempty"`
	EscapedPath  string      `json:"escaped_path,omitempty"`
	Hunks        []DiffHunk  `json:"hunks,omitempty"` // Detailed hunk information for the diff
	Links        *Links      `json:"links,omitempty"`
}

// FileContent represents the content of a file at a specific commit.
type FileContent struct {
	Path    string `json:"path"`
	Commit  string `json:"commit"`
	Content string `json:"content"`
}

// FileContentMeta provides metadata about a file's content.
type FileContentMeta struct {
	Size     int    `json:"size"`
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
}

// FileContentResult is the result returned by the service layer for file content.
type FileContentResult struct {
	Content string          `json:"content"`
	Meta    FileContentMeta `json:"meta"`
}

// Diff represents a raw diff as a string.
type Diff string

// DiffHunk represents a hunk of changes in a diff.
// This struct can be extended with more fields as needed.
type DiffHunk struct {
	Content  string `json:"content,omitempty"`
	NewLines int    `json:"new_lines,omitempty"`
	OldLines int    `json:"old_lines,omitempty"`
	Type     string `json:"type,omitempty"`
}

// PRCommentRequest represents the payload for adding a comment to a pull request.
// Supports both general and inline comments.
type PRCommentRequest struct {
	Content  string `json:"content"`
	FilePath string `json:"file_path,omitempty"` // Optional, for inline comments
	LineFrom int    `json:"line_from,omitempty"` // Optional, for inline comments
	LineTo   int    `json:"line_to,omitempty"`   // Optional, for inline comments
}

// PRChangeRequest represents the payload for requesting changes on a pull request.
// Bitbucket's API does not require a payload for this operation.

// PRComment represents a comment on a Bitbucket pull request.
// Supports both general PR comments and inline code comments.
type PRComment struct {
	ID      int64 `json:"id"`
	Content struct {
		Raw string `json:"raw"`
	} `json:"content"`
	Author    *Account  `json:"user"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
	Pending   bool      `json:"pending,omitempty"`

	// For inline comments
	Inline *InlineContext `json:"inline,omitempty"`

	// For threaded comments
	Parent *struct {
		ID int64 `json:"id"`
	} `json:"parent,omitempty"`
}

// InlineContext represents the file and line context for inline comments.
type InlineContext struct {
	Path string `json:"path"`
	From int    `json:"from,omitempty"`
	To   int    `json:"to,omitempty"`
}

// ListPRCommentsParams represents the parameters for listing PR comments.
type ListPRCommentsParams struct {
	Workspace string `json:"workspace"`
	RepoSlug  string `json:"repo_slug"`
	PRID      int64  `json:"pr_id"`
}

// ListPRCommentsResponse represents the response for listing PR comments.
// This matches the Bitbucket API response structure directly.
type ListPRCommentsResponse struct {
	Values []PRComment `json:"values"`
}

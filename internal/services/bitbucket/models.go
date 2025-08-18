package bitbucket

import "time"

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

// DiffStat represents a summary of changes made to a file between two commits.
type DiffStat struct {
	Type         string `json:"type,omitempty"`
	Status       string `json:"status,omitempty"`
	LinesAdded   int    `json:"lines_added,omitempty"`
	LinesRemoved int    `json:"lines_removed,omitempty"`
	Old          string `json:"old,omitempty"`
	New          string `json:"new,omitempty"`
	Path         string `json:"path,omitempty"`
	EscapedPath  string `json:"escaped_path,omitempty"`
	Hunks        []any  `json:"hunks,omitempty"` // TODO: define hunk structure if needed
	Links        *Links `json:"links,omitempty"`
}

// FileContent represents the content of a file at a specific commit.
type FileContent struct {
	Path    string `json:"path"`
	Commit  string `json:"commit"`
	Content string `json:"content"`
}

// Diff represents a raw diff as a string.
type Diff string

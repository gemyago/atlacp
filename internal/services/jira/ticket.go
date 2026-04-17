package jira

import "time"

// User represents a Jira user.
type User struct {
	AccountID    string `json:"accountId,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Active       bool   `json:"active,omitempty"`
	TimeZone     string `json:"timeZone,omitempty"`
	Self         string `json:"self,omitempty"`
}

// Status represents a Jira issue status.
type Status struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	StatusCategory struct {
		ID   int    `json:"id,omitempty"`
		Key  string `json:"key,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"statusCategory,omitzero"`
	Self string `json:"self,omitempty"`
}

// Priority represents a Jira issue priority.
type Priority struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Self    string `json:"self,omitempty"`
	IconURL string `json:"iconUrl,omitempty"`
}

// IssueType represents a Jira issue type.
type IssueType struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconUrl,omitempty"`
	Subtask     bool   `json:"subtask,omitempty"`
	Self        string `json:"self,omitempty"`
}

// Project represents a Jira project.
type Project struct {
	ID         string `json:"id,omitempty"`
	Key        string `json:"key,omitempty"`
	Name       string `json:"name,omitempty"`
	ProjectURL string `json:"projectUrl,omitempty"`
	Self       string `json:"self,omitempty"`
}

// Comment represents a Jira issue comment.
type Comment struct {
	ID           string    `json:"id,omitempty"`
	Author       User      `json:"author,omitzero"`
	Body         string    `json:"body,omitempty"`
	Created      time.Time `json:"created,omitzero"`
	Updated      time.Time `json:"updated,omitzero"`
	JSDPublic    bool      `json:"jsdPublic,omitempty"`
	Self         string    `json:"self,omitempty"`
	UpdateAuthor User      `json:"updateAuthor,omitzero"`
}

// Comments represents a collection of Jira issue comments.
type Comments struct {
	Comments   []Comment `json:"comments,omitempty"`
	MaxResults int       `json:"maxResults,omitempty"`
	Total      int       `json:"total,omitempty"`
	StartAt    int       `json:"startAt,omitempty"`
}

// Attachment represents a Jira issue attachment.
type Attachment struct {
	ID        string    `json:"id,omitempty"`
	Filename  string    `json:"filename,omitempty"`
	Author    User      `json:"author,omitzero"`
	Created   time.Time `json:"created,omitzero"`
	Size      int       `json:"size,omitempty"`
	MimeType  string    `json:"mimeType,omitempty"`
	Content   string    `json:"content,omitempty"`
	Thumbnail string    `json:"thumbnail,omitempty"`
	Self      string    `json:"self,omitempty"`
}

// Transition represents a Jira issue transition.
type Transition struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	To   Status `json:"to,omitzero"`
}

// Fields represents the fields of a Jira issue.
type Fields struct {
	Summary                       string       `json:"summary,omitempty"`
	Description                   string       `json:"description,omitempty"`
	Status                        Status       `json:"status,omitzero"`
	Priority                      Priority     `json:"priority,omitzero"`
	IssueType                     IssueType    `json:"issuetype,omitzero"`
	Project                       Project      `json:"project,omitzero"`
	Creator                       User         `json:"creator,omitzero"`
	Reporter                      User         `json:"reporter,omitzero"`
	Assignee                      User         `json:"assignee,omitzero"`
	Created                       time.Time    `json:"created,omitzero"`
	Updated                       time.Time    `json:"updated,omitzero"`
	ResolutionDate                time.Time    `json:"resolutiondate,omitzero"`
	Labels                        []string     `json:"labels,omitempty"`
	Comments                      Comments     `json:"comment,omitzero"`
	Attachments                   []Attachment `json:"attachment,omitempty"`
	FixVersions                   []any        `json:"fixVersions,omitempty"`
	Components                    []any        `json:"components,omitempty"`
	DueDate                       string       `json:"duedate,omitempty"`
	Watches                       any          `json:"watches,omitempty"`
	WorkRatio                     int          `json:"workratio,omitempty"`
	Subtasks                      []any        `json:"subtasks,omitempty"`
	Environment                   string       `json:"environment,omitempty"`
	TimeSpent                     int          `json:"timespent,omitempty"`
	AggregateTimeSpent            int          `json:"aggregatetimespent,omitempty"`
	TimeEstimate                  int          `json:"timeestimate,omitempty"`
	AggregateTimeOriginalEstimate int          `json:"aggregatetimeoriginalestimate,omitempty"`
	AggregateTimeEstimate         int          `json:"aggregatetimeestimate,omitempty"`
	TimeOriginalEstimate          int          `json:"timeoriginalestimate,omitempty"`
}

// Ticket represents a Jira issue/ticket.
type Ticket struct {
	ID             string            `json:"id,omitempty"`
	Key            string            `json:"key,omitempty"`
	Self           string            `json:"self,omitempty"`
	Fields         Fields            `json:"fields,omitzero"`
	RenderedFields any               `json:"renderedFields,omitempty"`
	Changelog      any               `json:"changelog,omitempty"`
	Transitions    []Transition      `json:"transitions,omitempty"`
	Names          map[string]string `json:"names,omitempty"`
	Schema         map[string]any    `json:"schema,omitempty"`
}

// TransitionRequest represents a request to transition a Jira issue.
type TransitionRequest struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
	Fields map[string]any `json:"fields,omitempty"`
	Update map[string]any `json:"update,omitempty"`
}

// LabelUpdateRequest represents a request to update labels on a Jira issue.
type LabelUpdateRequest struct {
	Update struct {
		Labels []LabelOperation `json:"labels,omitempty"`
	} `json:"update"`
}

// LabelOperation represents an add or remove operation for labels.
type LabelOperation struct {
	Add    string `json:"add,omitempty"`
	Remove string `json:"remove,omitempty"`
}

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
	} `json:"statusCategory,omitempty"`
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
	Author       User      `json:"author,omitempty"`
	Body         string    `json:"body,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
	JSDPublic    bool      `json:"jsdPublic,omitempty"`
	Self         string    `json:"self,omitempty"`
	UpdateAuthor User      `json:"updateAuthor,omitempty"`
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
	Author    User      `json:"author,omitempty"`
	Created   time.Time `json:"created,omitempty"`
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
	To   Status `json:"to,omitempty"`
}

// Fields represents the fields of a Jira issue.
type Fields struct {
	Summary                       string        `json:"summary,omitempty"`
	Description                   string        `json:"description,omitempty"`
	Status                        Status        `json:"status,omitempty"`
	Priority                      Priority      `json:"priority,omitempty"`
	IssueType                     IssueType     `json:"issuetype,omitempty"`
	Project                       Project       `json:"project,omitempty"`
	Creator                       User          `json:"creator,omitempty"`
	Reporter                      User          `json:"reporter,omitempty"`
	Assignee                      User          `json:"assignee,omitempty"`
	Created                       time.Time     `json:"created,omitempty"`
	Updated                       time.Time     `json:"updated,omitempty"`
	ResolutionDate                time.Time     `json:"resolutiondate,omitempty"`
	Labels                        []string      `json:"labels,omitempty"`
	Comments                      Comments      `json:"comment,omitempty"`
	Attachments                   []Attachment  `json:"attachment,omitempty"`
	FixVersions                   []interface{} `json:"fixVersions,omitempty"`
	Components                    []interface{} `json:"components,omitempty"`
	DueDate                       string        `json:"duedate,omitempty"`
	Watches                       interface{}   `json:"watches,omitempty"`
	WorkRatio                     int           `json:"workratio,omitempty"`
	Subtasks                      []interface{} `json:"subtasks,omitempty"`
	Environment                   string        `json:"environment,omitempty"`
	TimeSpent                     int           `json:"timespent,omitempty"`
	AggregateTimeSpent            int           `json:"aggregatetimespent,omitempty"`
	TimeEstimate                  int           `json:"timeestimate,omitempty"`
	AggregateTimeOriginalEstimate int           `json:"aggregatetimeoriginalestimate,omitempty"`
	AggregateTimeEstimate         int           `json:"aggregatetimeestimate,omitempty"`
	TimeOriginalEstimate          int           `json:"timeoriginalestimate,omitempty"`
}

// Ticket represents a Jira issue/ticket.
type Ticket struct {
	ID             string                 `json:"id,omitempty"`
	Key            string                 `json:"key,omitempty"`
	Self           string                 `json:"self,omitempty"`
	Fields         Fields                 `json:"fields,omitempty"`
	RenderedFields interface{}            `json:"renderedFields,omitempty"`
	Changelog      interface{}            `json:"changelog,omitempty"`
	Transitions    []Transition           `json:"transitions,omitempty"`
	Names          map[string]string      `json:"names,omitempty"`
	Schema         map[string]interface{} `json:"schema,omitempty"`
}

// TransitionRequest represents a request to transition a Jira issue.
type TransitionRequest struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
	Fields map[string]interface{} `json:"fields,omitempty"`
	Update map[string]interface{} `json:"update,omitempty"`
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

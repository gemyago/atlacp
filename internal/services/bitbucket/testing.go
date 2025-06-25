//go:build !release

package bitbucket

import "github.com/go-faker/faker/v4"

// PullRequestOpt is a function that configures a PullRequest.
type PullRequestOpt func(*PullRequest)

// WithPullRequestTitle sets the title of the pull request.
func WithPullRequestTitle(title string) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.Title = title
	}
}

// WithPullRequestDescription sets the description of the pull request.
func WithPullRequestDescription(description string) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.Description = description
	}
}

// WithPullRequestID sets the ID of the pull request.
func WithPullRequestID(id int) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.ID = id
	}
}

// WithPullRequestState sets the state of the pull request.
func WithPullRequestState(state string) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.State = state
	}
}

// WithPullRequestSourceBranch sets the source branch of the pull request.
func WithPullRequestSourceBranch(name string) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.Source.Branch.Name = name
	}
}

// WithPullRequestDestinationBranch sets the destination branch of the pull request.
func WithPullRequestDestinationBranch(name string) PullRequestOpt {
	return func(pr *PullRequest) {
		if pr.Destination == nil {
			pr.Destination = &PullRequestDestination{}
		}
		pr.Destination.Branch.Name = name
	}
}

// WithPullRequestCloseSourceBranch sets whether to close the source branch.
func WithPullRequestCloseSourceBranch(closeFlag bool) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.CloseSourceBranch = closeFlag
	}
}

// WithPullRequestAuthor sets the author of the pull request.
func WithPullRequestAuthor(author *PullRequestAuthor) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.Author = author
	}
}

// WithPullRequestReviewers sets the reviewers for the pull request.
func WithPullRequestReviewers(reviewers []PullRequestAuthor) PullRequestOpt {
	return func(pr *PullRequest) {
		pr.Reviewers = reviewers
	}
}

// PullRequestAuthorOpt is a function that configures a PullRequestAuthor.
type PullRequestAuthorOpt func(*PullRequestAuthor)

// WithAuthorUsername sets the username of the author.
func WithAuthorUsername(username string) PullRequestAuthorOpt {
	return func(a *PullRequestAuthor) {
		a.Username = username
	}
}

// WithAuthorDisplayName sets the display name of the author.
func WithAuthorDisplayName(displayName string) PullRequestAuthorOpt {
	return func(a *PullRequestAuthor) {
		a.DisplayName = displayName
	}
}

// NewRandomPullRequest generates a random PullRequest for testing.
func NewRandomPullRequest(opts ...PullRequestOpt) *PullRequest {
	pr := &PullRequest{
		ID:          int(faker.RandomUnixTime()),
		Title:       "PR-" + faker.Sentence(),
		Description: faker.Paragraph(),
		State:       "OPEN",
		Source: PullRequestSource{
			Branch: PullRequestBranch{
				Name: "feature/" + faker.Word(),
			},
		},
		Destination: &PullRequestDestination{
			Branch: PullRequestBranch{
				Name: "main",
			},
		},
		CloseSourceBranch: true,
		Author:            NewRandomPullRequestAuthor(),
		Type:              "pullrequest",
	}

	// Apply all options
	for _, opt := range opts {
		opt(pr)
	}

	return pr
}

// NewRandomPullRequestAuthor generates a random PullRequestAuthor for testing.
func NewRandomPullRequestAuthor(opts ...PullRequestAuthorOpt) *PullRequestAuthor {
	username := faker.Username()
	author := &PullRequestAuthor{
		AccountID:   faker.UUIDHyphenated(),
		DisplayName: faker.Name(),
		Nickname:    username,
		Username:    username,
		UUID:        faker.UUIDHyphenated(),
		Type:        "user",
	}

	// Apply all options
	for _, opt := range opts {
		opt(author)
	}

	return author
}

// NewRandomParticipant generates a random Participant for testing.
func NewRandomParticipant(approved bool) *Participant {
	state := "changes_requested"
	if approved {
		state = "approved"
	}

	return &Participant{
		User:     *NewRandomPullRequestAuthor(),
		Role:     "REVIEWER",
		Approved: approved,
		State:    state,
		Type:     "participant",
	}
}

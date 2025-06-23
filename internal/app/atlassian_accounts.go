package app

// AtlassianAccount represents configuration for a single Atlassian account.
type AtlassianAccount struct {
	// Friendly name of the account
	Name string `json:"name"`

	// Is this the default account
	Default bool `json:"default"`

	// Bitbucket-specific configuration (optional)
	Bitbucket *BitbucketAccount `json:"bitbucket,omitempty"`

	// Jira-specific configuration (optional)
	Jira *JiraAccount `json:"jira,omitempty"`
}

// BitbucketAccount contains Bitbucket-specific account configuration.
type BitbucketAccount struct {
	// API token for authentication
	Token string `json:"token"`

	// Workspace is the Bitbucket workspace/username for this account
	Workspace string `json:"workspace"`
}

// JiraAccount contains Jira-specific account configuration.
type JiraAccount struct {
	// API token for authentication
	Token string `json:"token"`

	// Domain is the Jira cloud instance domain (e.g., "mycompany" for mycompany.atlassian.net)
	Domain string `json:"domain"`
}

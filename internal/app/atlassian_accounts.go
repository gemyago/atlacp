package app

// AtlassianAccount represents configuration for a single Atlassian account.
type AtlassianAccount struct {
	// Friendly name of the account
	Name string `json:"name"`

	// Is this the default account
	Default bool `json:"default"`

	// Bitbucket-specific configuration (optional)
	Bitbucket *AtlassianToken `json:"bitbucket,omitempty"`

	// Jira-specific configuration (optional)
	Jira *AtlassianToken `json:"jira,omitempty"`
}

// AtlassianToken contains authentication token information.
type AtlassianToken struct {
	// Token type (e.g., "Bearer", "Basic"). Defaults to "Bearer" if not specified.
	Type string `json:"type,omitempty"`

	// Token value for authentication
	Value string `json:"value"`
}

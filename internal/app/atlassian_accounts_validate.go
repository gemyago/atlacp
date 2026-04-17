package app

import (
	"errors"
	"fmt"
)

// ValidateAtlassianAccounts checks that the slice satisfies the same invariants as
// the file-backed Atlassian accounts repository: at least one account, unique names,
// exactly one default, and required token fields when Bitbucket and/or Jira are present.
func ValidateAtlassianAccounts(accounts []AtlassianAccount) error {
	if len(accounts) == 0 {
		return errors.New("no accounts configured")
	}

	accountNames := make(map[string]bool)
	foundDefault := false

	for _, account := range accounts {
		if err := validateBasicAtlassianAccountProperties(account, accountNames); err != nil {
			return err
		}
		accountNames[account.Name] = true

		if err := validateAtlassianServiceConfigs(account); err != nil {
			return err
		}

		if account.Default {
			if foundDefault {
				return errors.New("multiple default accounts defined")
			}
			foundDefault = true
		}
	}

	if !foundDefault {
		return errors.New("no default account specified")
	}

	return nil
}

func validateBasicAtlassianAccountProperties(account AtlassianAccount, existingNames map[string]bool) error {
	if existingNames[account.Name] {
		return fmt.Errorf("duplicate account name: %s", account.Name)
	}

	if account.Name == "" {
		return errors.New("account missing name")
	}

	if account.Bitbucket == nil && account.Jira == nil {
		return fmt.Errorf("account %s must have at least one service configured", account.Name)
	}

	return nil
}

func validateAtlassianServiceConfigs(account AtlassianAccount) error {
	if account.Bitbucket != nil {
		if err := validateAtlassianBitbucketConfig(account); err != nil {
			return err
		}
	}

	if account.Jira != nil {
		if err := validateAtlassianJiraConfig(account); err != nil {
			return err
		}
	}

	return nil
}

func validateAtlassianBitbucketConfig(account AtlassianAccount) error {
	if account.Bitbucket.Value == "" {
		return fmt.Errorf("account %s is missing Bitbucket token value", account.Name)
	}
	if account.Bitbucket.Type == "" {
		return fmt.Errorf("account %s is missing Bitbucket token type", account.Name)
	}
	return nil
}

func validateAtlassianJiraConfig(account AtlassianAccount) error {
	if account.Jira.Value == "" {
		return fmt.Errorf("account %s is missing Jira token value", account.Name)
	}
	if account.Jira.Type == "" {
		return fmt.Errorf("account %s is missing Jira token type", account.Name)
	}
	return nil
}

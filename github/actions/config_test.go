package actions_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/actions/actions-runner-controller/github/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHubConfig(t *testing.T) {
	t.Run("when given a valid URL", func(t *testing.T) {
		tests := []struct {
			configURL string
			expected  *actions.GitHubConfig
		}{
			{
				configURL: "https://github.com/org/repo",
				expected: &actions.GitHubConfig{
					Scope:        actions.GitHubScopeRepository,
					Enterprise:   "",
					Organization: "org",
					Repository:   "repo",
					IsHosted:     true,
				},
			},
			{
				configURL: "https://github.com/org",
				expected: &actions.GitHubConfig{
					Scope:        actions.GitHubScopeOrganization,
					Enterprise:   "",
					Organization: "org",
					Repository:   "",
					IsHosted:     true,
				},
			},
			{
				configURL: "https://github.com/enterprises/my-enterprise",
				expected: &actions.GitHubConfig{
					Scope:        actions.GitHubScopeEnterprise,
					Enterprise:   "my-enterprise",
					Organization: "",
					Repository:   "",
					IsHosted:     true,
				},
			},
			{
				configURL: "https://www.github.com/org",
				expected: &actions.GitHubConfig{
					Scope:        actions.GitHubScopeOrganization,
					Enterprise:   "",
					Organization: "org",
					Repository:   "",
					IsHosted:     true,
				},
			},
			{
				configURL: "https://github.localhost/org",
				expected: &actions.GitHubConfig{
					Scope:        actions.GitHubScopeOrganization,
					Enterprise:   "",
					Organization: "org",
					Repository:   "",
					IsHosted:     true,
				},
			},
			{
				configURL: "https://my-ghes.com/org",
				expected: &actions.GitHubConfig{
					Scope:        actions.GitHubScopeOrganization,
					Enterprise:   "",
					Organization: "org",
					Repository:   "",
					IsHosted:     false,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.configURL, func(t *testing.T) {
				parsedURL, err := url.Parse(test.configURL)
				require.NoError(t, err)
				test.expected.ConfigURL = parsedURL

				cfg, err := actions.ParseGitHubConfigFromURL(test.configURL)
				require.NoError(t, err)
				assert.Equal(t, test.expected, cfg)
			})
		}
	})

	t.Run("when given an invalid URL", func(t *testing.T) {
		invalidURLs := []string{
			"https://github.com/",
			"https://github.com",
			"https://github.com/some/random/path",
		}

		for _, u := range invalidURLs {
			_, err := actions.ParseGitHubConfigFromURL(u)
			require.Error(t, err)
			assert.True(t, errors.Is(err, actions.ErrInvalidGitHubConfigURL))
		}
	})
}

func TestGitHubConfig_GitHubAPIURL(t *testing.T) {
	t.Run("when hosted", func(t *testing.T) {
		config, err := actions.ParseGitHubConfigFromURL("https://github.com/org/repo")
		require.NoError(t, err)

		result := config.GitHubAPIURL("/some/path")
		assert.Equal(t, "https://api.github.com/some/path", result.String())
	})
	t.Run("when not hosted", func(t *testing.T) {})
}

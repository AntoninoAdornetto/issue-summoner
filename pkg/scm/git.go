package scm

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	GH = "GitHub"
	GL = "GitLab"
	BB = "BitBucket"
)

// @TODO can GlobalUserName and RepoName functions be deleted?
// We are now using the device flow and the mentioned functions could be useless since
// we are creating an access token for the user after they authorize the application.

type GitConfig struct {
	UserName       string
	RepositoryName string
	Token          string
	Scm            string // GitHub || GitLab || BitBucket ...
}

// GitConfigManager interface allows us to have different adapters for each
// source code management system that we would like to use. We can have different
// implementations for GitHub, GitLab, BitBucket and so on.
// Authorize creates an access token with scopes that will allow us to read/write issues
// ReadToken checks if there is an access token in ~/.config/issue-summoner/config.json
type GitConfigManager interface {
	Authorize() error
	IsAuthorized() bool
}

func GetGitConfig(scm string) GitConfigManager {
	switch scm {
	default:
		return &GitHubManager{}
	}
}

type ScmTokenConfig struct {
	AccessToken string
}

type IssueSummonerConfig = map[string]ScmTokenConfig

// WriteToken accepts an access token and the source code management platform
// (GitHub, GitLab etc...) and will write the token to a configuration file.
// This will be used to authorize future requests for reporting issues.
func WriteToken(token string, scm string) error {
	config := make(map[string]ScmTokenConfig)

	usr, err := user.Current()
	if err != nil {
		return err
	}

	home := usr.HomeDir
	path := filepath.Join(home, ".config", "issue-summoner")

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	configFile := filepath.Join(path, "config.json")
	file, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	switch scm {
	default:
		config["GitHub"] = ScmTokenConfig{
			AccessToken: token,
		}
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

// @TODO refactor WriteToken & CheckForAccess functions.
// There is some DRY code in the two functions that I would like to refactor.
// Specifically for getting the current directory, home dir and joining the paths
// for the configuration file.
func CheckForAccess(scm string) (bool, error) {
	config := make(map[string]ScmTokenConfig)
	authorized := false

	usr, err := user.Current()
	if err != nil {
		return authorized, err
	}

	home := usr.HomeDir
	configFile := filepath.Join(home, ".config", "issue-summoner", "config.json")

	file, err := os.OpenFile(configFile, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return authorized, err
		} else {
			return authorized, errors.New("Error opening file")
		}
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return authorized, errors.New("Error decoding config file")
	}

	return config[scm].AccessToken != "", nil
}

// GlobalUserName uses the **git config** command to retrieve the global
// configuration options. Specifically, the user.name option. The userName is
// read and set onto the reciever's (GitConfig) UserName property. This will be used
func GlobalUserName() (string, error) {
	var out strings.Builder
	cmd := exec.Command("git", "config", "--global", "user.name")
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	userName := strings.TrimSpace(out.String())
	if userName == "" {
		return "", errors.New("global userName option not set. See man git config for more details")
	}

	return userName, nil
}

func RepoName() (string, error) {
	var out strings.Builder
	cmd := exec.Command("git", "remote", "-v")
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	repoName := extractRepoName(out.String())
	if repoName == "" {
		return "", errors.New("Failed to get repo name")
	}

	return repoName, nil
}

// extractRepoName takes the output from the `git remote -v` command as input (origins) and outputs the repository name.
// The function can handle both ssh and https origins.
// Git does not offer a command that outputs the repository name directly
func extractRepoName(origins string) string {
	for _, line := range strings.Split(origins, "\n") {
		if strings.Contains(line, "(push)") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				repoURL := fields[1]
				parts := strings.Split(repoURL, "/")
				if len(parts) > 1 {
					repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
					return repo
				}
			}
		}
	}
	return ""
}

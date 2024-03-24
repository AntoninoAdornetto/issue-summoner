package scm

import (
	"errors"
	"os/exec"
	"strings"
)

type GitConfig struct {
	UserName       string
	RepositoryName string
	Token          string
}

// GlobalUserName uses the **git config** command to retrieve the global
// configuration options. Specifically, the user.name option. The userName is
// read and set onto the reciever's (GitConfig) UserName property. This will be used
func (gc *GitConfig) GlobalUserName() error {
	var out strings.Builder
	cmd := exec.Command("git", "config", "--global", "user.name")
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	userName := out.String()
	if userName == "" {
		return errors.New("global userName option not set. See man git config for more details")
	}

	gc.UserName = userName
	return nil
}

func (gc *GitConfig) RepoName() error {
	var out strings.Builder
	cmd := exec.Command("git", "remote", "-v")
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	gc.RepositoryName = extractRepoName(out.String())
	return nil
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

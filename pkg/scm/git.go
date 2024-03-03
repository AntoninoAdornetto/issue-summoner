package scm

import (
	"os/exec"
	"strings"
)

type GitConfig struct {
	UserName       string
	RepositoryName string
	Token          string
}

func (gc *GitConfig) User() error {
	var out strings.Builder
	cmd := exec.Command("git", "config", "--global", "user.name")
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	gc.UserName = out.String()
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

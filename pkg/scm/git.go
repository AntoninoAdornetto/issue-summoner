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

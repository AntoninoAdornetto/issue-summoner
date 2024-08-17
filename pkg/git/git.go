package git

import (
	"errors"
	"fmt"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/common"
)

type sourceCodeManagement = string

const (
	Github    sourceCodeManagement = "github"
	Gitlab    sourceCodeManagement = "gitlab"
	Bitbucket sourceCodeManagement = "bitbucket"
)

type GitManager interface {
	Authorize() error
	Report(issue ReportIssueReq) (int64, error)
	Authenticated() bool
}

type ReportIssueReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewGitManager(scm sourceCodeManagement, repo *Repository) (GitManager, error) {
	conf, err := common.ReadConfig()
	if err != nil {
		return nil, err
	}

	switch scm {
	case Bitbucket:
		return nil, errors.New("bitbucket is not supported yet. Check back soon")
	case Gitlab:
		return nil, errors.New("gitlab is not supported yet. Check back soon")
	case Github:
		return &githubManager{conf: conf, repo: repo}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported scm. expected one of the following: %s %s %s but got %s",
			Github,
			Gitlab,
			Bitbucket,
			scm,
		)
	}
}

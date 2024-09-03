package git

import (
	"errors"
	"fmt"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/common"
)

type sourceCodeHost = string

const (
	Github    sourceCodeHost = "github"
	Gitlab    sourceCodeHost = "gitlab"
	Bitbucket sourceCodeHost = "bitbucket"
)

type GitManager interface {
	Authorize() error
	Report(issue ReportRequest, res chan ReportResponse)
	GetStatus(issueNum, index int, res chan StatusResponse)
	Authenticated() bool
}

type ReportRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Index int    // index location in [IssueManager.Issues] slice in the issue package
}

type ReportResponse struct {
	ID    int // issue number
	Err   error
	Index int // index location in [IssueManager.Issues] slice in the issue package
}

type StatusResponse struct {
	Resolved bool
	Err      error
	Index    int // index location in [IssueManager.Issues] slice in the issue package
}

func NewGitManager(sch sourceCodeHost, repo *Repository) (GitManager, error) {
	conf, err := common.ReadConfig()
	if err != nil {
		return nil, err
	}

	switch sch {
	case Bitbucket:
		return nil, errors.New("bitbucket is not supported yet. Check back soon")
	case Gitlab:
		return nil, errors.New("gitlab is not supported yet. Check back soon")
	case Github:
		return &githubManager{conf: conf, repo: repo}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported source code host. expected one of the following: %s %s %s but got %s",
			Github,
			Gitlab,
			Bitbucket,
			sch,
		)
	}
}

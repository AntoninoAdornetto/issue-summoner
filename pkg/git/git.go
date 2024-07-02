package git

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/utils"
)

type sourceCodeManager = string

const (
	GITHUB    sourceCodeManager = "github"
	GITLAB    sourceCodeManager = "gitlab"
	BITBUCKET sourceCodeManager = "bitbucket"
)

type GitManager interface {
	Authorize() error
	IsAuthorized() bool
	Report(req ReportRequest, res chan ReportResponse)
}

type ReportRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Index int
}

type ReportResponse struct {
	ID    int64 `json:"id"`
	Index int
	Err   error
}

// NewGitManager allows each source code management platform to support the methods in the GitManager
// interface. Meaning, github, bitbucket and gitlab can each have their own implementations of Authorizing
// and Reporting issues. The function will also assist in reading/writing to the configuration file. The
// configuration file contains access tokens that are needed to report issues. If a config file is not present,
// one will be created at the time of invoking the function.
func NewGitManager(scm sourceCodeManager, repo *Repository) (GitManager, error) {
	config, err := prepareConfig()
	if err != nil {
		return nil, err
	}

	switch scm {
	case BITBUCKET:
		return nil, errors.New("bitbucket is not supported as of yet. Check back soon")
	case GITLAB:
		return nil, errors.New("gitlab is not supported as of yet. Check back soon")
	case GITHUB:
		return &githubManager{config: config, repo: repo}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported scm. expected one of the following: %s %s %s but got %s",
			GITHUB,
			GITLAB,
			BITBUCKET,
			scm,
		)
	}
}

func prepareConfig() (map[string]utils.IssueSummonerConfig, error) {
	conf := make(map[string]utils.IssueSummonerConfig)
	data, err := utils.ReadIssueSummonerConfig()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		data, err = json.Marshal(utils.Config)
		if err != nil {
			return nil, err
		}
	}

	if err = json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

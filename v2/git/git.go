package git

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/utils"
)

type sourceCodeManager = string

const (
	GITHUB    sourceCodeManager = "github"
	GITLAB    sourceCodeManager = "gitlab"
	BITBUCKET sourceCodeManager = "bitbucket"
)

var (
	conf = map[string]IssueSummonerConfig{
		GITHUB:    {},
		GITLAB:    {},
		BITBUCKET: {},
	}
)

type GitManager interface {
	Authorize() error
	IsAuthorized() bool
	Report(issue CodeIssue) (ReportedIssue, error)
}

type CodeIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Index int
}

type ReportedIssue struct {
	ID    int64 `json:"id"`
	Index int
}

type IssueSummonerConfig struct {
	Auth authConfig `json:"auth"`
}

type authConfig struct {
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// NewGitManager will handle specific operations for each source code management platform
// that the program supports. Each scm should satisfy the GitManager interface for authorizing
// issue summoner to submit issues to both public/private repositories and support functionality
// for posting new issues to your code bases. If the issue summoner config file is not created
// this function will create both the configuration directory and json file. The config file
// contains details such as access tokens for making http requests to different scm's.
func NewGitManager(scm sourceCodeManager, repo *Repository) (GitManager, error) {
	conf := make(map[string]IssueSummonerConfig)
	data, err := utils.ReadIssueSummonerConfig()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		data, err = NewConfig()
		if err != nil {
			return nil, err
		}
	}

	if err = json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	// @TODO add support for GitLab and Bitbucket
	switch scm {
	case GITHUB:
		return &githubManager{config: conf, repo: repo}, nil
	default:
		return nil, nil
	}
}

func NewConfig() ([]byte, error) {
	return json.Marshal(conf)
}

func buildURL(base string, queryParams map[string]string, paths ...string) (string, error) {
	u, err := url.JoinPath(base, paths...)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	for key, val := range queryParams {
		params.Set(key, val)
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", u, params.Encode()), nil
	}

	return u, nil
}

func makeRequest(method, url string, body io.Reader, h http.Header) ([]byte, error) {
	var res []byte
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return res, err
	}

	req.Header = h
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()
	res, err = io.ReadAll(resp.Body)
	return res, err
}

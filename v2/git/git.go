package git

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type sourceCodeManager = string

const (
	GITHUB    sourceCodeManager = "github"
	GITLAB    sourceCodeManager = "gitlab"
	BITBUCKET sourceCodeManager = "bitbucket"
)

type GitManager interface {
	Authorize() error
}

type gitAccessToken struct {
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

func NewGitManager(scm sourceCodeManager, repo *Repository) (GitManager, error) {
	switch scm {
	case GITHUB:
		return &githubManager{repo: repo}, nil
	default:
		return nil, nil
	}
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

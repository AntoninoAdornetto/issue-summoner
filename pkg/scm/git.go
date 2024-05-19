package scm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	GITHUB    = "github"
	GITLAB    = "gitlab"
	BITBUCKET = "bitbucket"
)

type GitIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// GitConfigManager provides flexibility to have different implementations
// of Authorize and Report for each source code management platform supported
type GitConfigManager interface {
	Authorize() error
	Report(issues []GitIssue) <-chan int64
}

// @TODO add other source code management structs to NewGitManager once their implementations are created
func NewGitManager(scm, userName, repoName string) (GitConfigManager, error) {
	switch scm {
	case GITHUB:
		return &GitHubManager{repoName: repoName, userName: userName}, nil
	default:
		return nil, fmt.Errorf(
			"expected to receive scm with value of %s, %s, or %s but got %s",
			GITHUB,
			GITLAB,
			BITBUCKET,
			scm,
		)
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
	path, err := getConfigDirPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	// @TODO add remaining source code management platforms once other adapters are implemented
	switch scm {
	default:
		config[GITHUB] = ScmTokenConfig{
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

func ReadAccessToken(scm string) (string, error) {
	config := make(map[string]ScmTokenConfig)
	path, err := getConfigFilePath()
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return "", err
		} else {
			return "", errors.New("Error opening file")
		}
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return "", errors.New("Error decoding config file")
	}

	accessToken := config[scm].AccessToken
	if accessToken == "" {
		return "", errors.New("Access token does not exist")
	}

	return accessToken, nil
}

// ExtractUserRepoName takes the output from <git remote --verbose> command
// as input and attempts to extract the user name and repository name from out
func ExtractUserRepoName(out []byte) (string, string, error) {
	if len(out) == 0 {
		return "", "", errors.New(
			"expected to receive the output from <git remote -v> but got empty byte slice",
		)
	}

	line := bytes.Split(out, []byte("\n"))[0]
	// fields will give us -> ["origin", "url (https | ssh)", "(push) | (pull)"]
	// we only care about the url since it contains both the username and repo name
	fields := bytes.Fields(line)
	if len(fields) < 2 {
		return "", "", fmt.Errorf(
			"expected to receive the origin and url but got %s",
			string(fields[0]),
		)
	}

	url := fields[1]
	if bytes.HasPrefix(url, []byte("https")) {
		userName, repoName := extractFromHTTPS(url)
		return userName, repoName, nil
	}

	if bytes.HasPrefix(url, []byte("git")) {
		userName, repoName := extractFromSSH(url)
		return userName, repoName, nil
	}

	return "", "", fmt.Errorf("expected a https or ssh url but got %s", string(url))
}

func extractFromHTTPS(url []byte) (string, string) {
	split := bytes.SplitAfter(url, []byte("https://"))[1]
	sep := bytes.Split(split, []byte("/"))
	userName, repoName := sep[1], sep[2]
	return string(userName), string(bytes.TrimSuffix(repoName, []byte(".git")))
}

func extractFromSSH(url []byte) (string, string) {
	split := bytes.SplitAfter(url, []byte(":"))[1]
	sep := bytes.Split(split, []byte("/"))
	userName, repoName := sep[0], sep[1]
	return string(userName), string(bytes.TrimSuffix(repoName, []byte(".git")))
}

func getConfigDirPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	homeDir := usr.HomeDir
	return filepath.Join(homeDir, ".config", "issue-summoner"), nil
}

func getConfigFilePath() (string, error) {
	configDir, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

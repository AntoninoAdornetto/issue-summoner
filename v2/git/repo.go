package git

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type Repository struct {
	WorkTree                string
	Dir                     string
	RepoName                string
	UserName                string
	repositoryFormatVersion int
	remoteUrl               string
}

// NewRepository will attempt to locate the working directory of your
// local git project and parse the ini-like config file contained in the .git/ dir.
// the config file is checked for a repository version of 0 and some basic string manipulation
// is executed to get the user and repo name of the git project. This will be utilized later on
// to submit issues to different source code management platforms
func NewRepository(path string) (*Repository, error) {
	path, err := findRepositoryPath(path)
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		WorkTree:                path,
		Dir:                     filepath.Join(path, ".git"),
		repositoryFormatVersion: -1,
	}

	if err := repo.parseGitConfig(); err != nil {
		return nil, err
	}

	return repo, nil
}

func findRepositoryPath(wd string) (string, error) {
	if wd == "/" {
		return "", errors.New("expected to find a local git repo but found none.")
	}

	if _, err := os.Stat(filepath.Join(wd, ".git")); err != nil {
		if os.IsNotExist(err) {
			return findRepositoryPath(filepath.Join(wd, "../"))
		}
		return "", err
	}

	return wd, nil
}

func (repo *Repository) parseGitConfig() error {
	configFile, err := os.Open(filepath.Join(repo.Dir, "config"))
	if err != nil {
		return err
	}

	defer configFile.Close()
	scanner := bufio.NewScanner(configFile)

	for scanner.Scan() {
		if repo.ok() {
			break
		}

		line := scanner.Text()
		if err := repo.getConfigDetails(line); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (repo *Repository) getConfigDetails(line string) error {
	// split on strs -> e.g. "url = https://github.com/..."
	sep := strings.Split(line, "=")
	if len(sep) < 2 {
		return nil
	}

	for i := range sep {
		sep[i] = strings.TrimFunc(sep[i], unicode.IsSpace)
	}

	key, val := sep[0], sep[1]
	if strings.HasPrefix(key, ";") {
		return nil
	}

	switch {
	case repo.repositoryFormatVersion == -1 && strings.Contains(key, "repositoryformatversion"):
		return repo.handleRepoFormatVersion(val)
	case repo.remoteUrl == "" && strings.Contains(key, "url"):
		return repo.handleRepoUrl(val)
	default:
		return nil
	}
}

func (repo *Repository) handleRepoFormatVersion(val string) error {
	version, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	if version != 0 {
		return fmt.Errorf("expected repository format version of 0 but got %d", version)
	}

	repo.repositoryFormatVersion = version
	return nil
}

func (repo *Repository) handleRepoUrl(val string) error {
	if val == "" {
		return errors.New("expected remote url to be present but got empty string.")
	}

	repo.remoteUrl = val
	if err := repo.extractRepoAndUserName(); err != nil {
		return err
	}
	return nil
}

func (repo *Repository) extractRepoAndUserName() error {
	splitBuf := make([]byte, 0, 5)

	switch {
	case strings.HasPrefix(repo.remoteUrl, "https"):
		splitBuf = []byte(".com/")
	case strings.HasPrefix(repo.remoteUrl, "git@"):
		splitBuf = []byte(".com:")
	default:
		return fmt.Errorf(
			"expected https or ssh protocol but got unexpected url of %s",
			repo.remoteUrl,
		)
	}

	split := strings.SplitAfter(repo.remoteUrl, string(splitBuf))
	if len(split) < 2 {
		return fmt.Errorf(
			"unable to split url %s by separator %s",
			repo.remoteUrl,
			string(splitBuf),
		)
	}

	rm := split[1]
	rm = strings.TrimSuffix(rm, ".git")
	split = strings.Split(rm, "/")
	if len(split) < 2 {
		return errors.New("unable to extract username and repository name")
	}

	userName, repoName := split[0], split[1]
	repo.RepoName = repoName
	repo.UserName = userName
	return nil
}

func (repo *Repository) ok() bool {
	return repo.repositoryFormatVersion != -1 && repo.remoteUrl != ""
}

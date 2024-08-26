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
	WorkTree          string
	Dir               string
	RepoName          string
	UserName          string
	repoFormatVersion int
	remoteUrl         string
}

func NewRepository(path string) (*Repository, error) {
	workPath, err := findRepository(path)
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		WorkTree:          workPath,
		Dir:               filepath.Join(workPath, ".git"),
		repoFormatVersion: -1,
	}

	if err := repo.parseGitConfig(); err != nil {
		return nil, err
	}

	return repo, nil
}

// recursively move up the file tree till we locate a .git directory
// or reach the root dir.
func findRepository(path string) (string, error) {
	if path == "/" {
		return "", errors.New("failed to locate the work tree of your git project")
	}

	if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		if os.IsNotExist(err) {
			return findRepository(filepath.Join(path, "../"))
		}
		return "", err
	}

	return path, nil
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
		if err := repo.readLineDetails(line); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// the git config file is an ini-like file that we will read
// to grab the repository format version and remote url from.
// this is not meant to be a fully fledged ini parser nor did
// I want to write one. This is just a simple, quick and dirty
// solution to grab the data points we need for our use case.
func (repo *Repository) readLineDetails(line string) error {
	keyVal := strings.Split(line, "=")
	if len(keyVal) < 2 {
		return nil
	}

	for i := range keyVal {
		keyVal[i] = strings.TrimFunc(keyVal[i], unicode.IsSpace)
	}

	key, val := keyVal[0], keyVal[1]
	if strings.HasPrefix(key, ";") {
		// skip comments
		return nil
	}

	switch key {
	case "repositoryFormatVersion":
		return repo.extractRepoFormatVersion(val)
	case "url":
		return repo.extractRemoteUrl(val)
	default:
		return nil
	}
}

// repository format versions specify the rules for operating on the disk
func (repo *Repository) extractRepoFormatVersion(val string) error {
	version, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	if version != 0 {
		return fmt.Errorf(
			"repository format version of %d is currently unsupported. Expected version 0",
			version,
		)
	}

	repo.repoFormatVersion = version
	return nil
}

func (repo *Repository) extractRemoteUrl(val string) error {
	if val == "" {
		return errors.New("remote url value is empty")
	}

	repo.remoteUrl = val
	return repo.extractRepoDetails()
}

// extracts the user name and repo name that we will use for reporting
// issues to different source code hosting platforms
func (repo *Repository) extractRepoDetails() error {
	var buf []byte

	switch {
	case strings.HasPrefix(repo.remoteUrl, "https"):
		buf = []byte(".com/")
	case strings.HasPrefix(repo.remoteUrl, "git@"):
		buf = []byte(".com:")
	default:
		return fmt.Errorf(
			"expected https or ssh protocol but got unexpected url of %s",
			repo.remoteUrl,
		)
	}

	repoDetails := strings.SplitAfter(repo.remoteUrl, string(buf))
	if len(repoDetails) < 2 {
		return fmt.Errorf("failed to split url %s by separator %s", repo.remoteUrl, string(buf))
	}

	rm := repoDetails[1]
	rm = strings.TrimSuffix(rm, ".git")
	repoDetails = strings.Split(rm, "/")
	if len(repoDetails) < 2 {
		return fmt.Errorf(
			"failed to extract username and repo name from remote url %s",
			repo.remoteUrl,
		)
	}

	userName, repoName := repoDetails[0], repoDetails[1]
	repo.UserName = userName
	repo.RepoName = repoName
	return nil
}

// the only two properties we care about in the config file.
// remote url will be used to extract the user name and repo name
func (repo *Repository) ok() bool {
	return repo.repoFormatVersion != -1 && repo.remoteUrl != ""
}

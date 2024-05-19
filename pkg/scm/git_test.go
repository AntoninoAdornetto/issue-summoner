package scm_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/stretchr/testify/require"
)

const (
	ssh_remote_output = `origin	git@github.com:AntoninoAdornetto/issue-summoner.git (fetch)
	origin	git@github.com:AntoninoAdornetto/issue-summoner.git (push)`
	https_remote_output = `origin https://github.com/AntoninoAdornetto/issue-summoner.git (fetch)
  origin https://github.com/AntoninoAdornetto/issue-summoner.git (push)
	`
)

// should create a new GitHubManager struct
func TestNewGitManagerGitHub(t *testing.T) {
	gm, err := scm.NewGitManager(scm.GITHUB, "AntoninoAdornetto", "issue-summoner")
	require.NoError(t, err)
	require.IsType(t, &scm.GitHubManager{}, gm)
}

// should return an error when provided an unsupported source code management platform
func TestNewGitManagerUnsupported(t *testing.T) {
	gm, err := scm.NewGitManager("unsupported", "AntoninoAdornetto", "issue-summoner")
	require.Error(t, err)
	require.Empty(t, gm)
}

// should return the username and repo name when provided output
// from the git remote command that contains an https url
func TestExtractUserRepoNameHTTPS(t *testing.T) {
	expectedUser, expectedRepo := "AntoninoAdornetto", "issue-summoner"
	actualUser, actualRepo, err := scm.ExtractUserRepoName([]byte(https_remote_output))
	require.NoError(t, err)
	require.Equal(t, expectedRepo, actualRepo)
	require.Equal(t, expectedUser, actualUser)
}

// should return the username and repo name when provided output
// from the git remote command that contains an ssh url
func TestExtractUserRepoNameSSH(t *testing.T) {
	expectedUser, expectedRepo := "AntoninoAdornetto", "issue-summoner"
	actualUser, actualRepo, err := scm.ExtractUserRepoName([]byte(ssh_remote_output))
	require.NoError(t, err)
	require.Equal(t, expectedRepo, actualRepo)
	require.Equal(t, expectedUser, actualUser)
}

// should return empty user and repo name when provided empty byte slice as input
func TestExtractUserRepoNameNoOutput(t *testing.T) {
	userName, repoName, err := scm.ExtractUserRepoName([]byte{})
	require.Empty(t, userName)
	require.Empty(t, repoName)
	require.Error(t, err)
}

// should return empty user and repo name when provided unexpected output, such
// as a url that is neither https or ssh format
func TestExtractUserRepoNameUnknownProtocol(t *testing.T) {
	userName, repoName, err := scm.ExtractUserRepoName(
		[]byte("origin unknown.protocolformaturl, (push)"),
	)
	require.Empty(t, userName)
	require.Empty(t, repoName)
	require.Error(t, err)
}

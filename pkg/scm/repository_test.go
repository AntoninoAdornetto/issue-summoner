package scm_test

import (
	"path/filepath"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/stretchr/testify/require"
)

func TestNewRepository(t *testing.T) {
	path := "/home/user/gitproject"
	expectedWorkTree := path
	expectedDir := path + "/.git"

	repo := scm.NewRepository(path)
	require.Equal(t, expectedWorkTree, repo.WorkTree)
	require.Equal(t, expectedDir, repo.Dir)
}

func TestFindRepoSuccess(t *testing.T) {
	wd, err := filepath.Abs("../../testdata/exclude/")
	require.NoError(t, err)
	repo, err := scm.FindRepository(wd)
	require.NoError(t, err)
	require.NotNil(t, repo)
	baseWorkTree := filepath.Base(repo.WorkTree)
	baseDir := filepath.Base(repo.Dir)
	require.Equal(t, "issue-summoner", baseWorkTree)
	require.Equal(t, ".git", baseDir)
}

func TestFindRepoRootError(t *testing.T) {
	wd := "/"
	repo, err := scm.FindRepository(wd)
	require.Error(t, err)
	require.Nil(t, repo)
}

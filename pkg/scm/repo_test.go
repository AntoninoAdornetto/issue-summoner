package scm_test

import (
	"path/filepath"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/stretchr/testify/require"
)

func TestNewRepositoryWorkDirPath(t *testing.T) {
	path, err := filepath.Abs("../../")
	require.NoError(t, err)

	actual, err := scm.NewRepository(path)
	require.NoError(t, err)
	require.NotNil(t, actual)

	require.NotEmpty(t, actual.Dir)
	require.NotEmpty(t, actual.WorkTree)
	require.NotEmpty(t, actual.UserName)
	require.NotEmpty(t, actual.RepoName)

	require.Equal(t, "AntoninoAdornetto", actual.UserName)
	require.Equal(t, "issue-summoner", actual.RepoName)
}

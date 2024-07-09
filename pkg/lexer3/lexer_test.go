package lexer3_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func getTestDataSrc(t *testing.T, path string) ([]byte, string) {
	fileName := filepath.Base(path)
	srcFile, err := os.Open(path)
	require.NoError(t, err)

	defer srcFile.Close()
	srcCode, err := io.ReadAll(srcFile)
	require.NoError(t, err)
	return srcCode, fileName
}

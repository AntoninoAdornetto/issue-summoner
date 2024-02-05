package parser_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/parser"
	"github.com/stretchr/testify/require"
)

type MockIgnoreFileOpener struct {
	File *os.File
	Err  error
}

func (m MockIgnoreFileOpener) Open(fileName string) (*os.File, error) {
	return m.File, m.Err
}

func TestProcessIgnorePatterns(t *testing.T) {
	tempFile, err := os.CreateTemp("", "*.gitignore")
	require.NoError(t, err)

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Some example patterns that we can find in a .gitignore file
	patterns := []string{
		"/tmp",
		".pnp.js",
		"*.log",
		"src/old-impl/test/",
		"# Comment", // This should be ignored
		"",          // This should be ignored
	}

	expectedLength := 4

	expectedRegexpPatterns := []regexp.Regexp{
		*regexp.MustCompile("^\\/tmp"),
		*regexp.MustCompile(".pnp.js"),
		*regexp.MustCompile(".*\\.log"),
		*regexp.MustCompile("src/old-impl/test/"),
	}

	_, err = tempFile.WriteString(strings.Join(patterns, "\n"))
	require.NoError(t, err)

	_, err = tempFile.Seek(0, 0)
	require.NoError(t, err)

	mockIgnoreFileOpener := MockIgnoreFileOpener{File: tempFile}

	regexpPatterns, err := parser.ProcessIgnorePatterns(
		tempFile.Name(),
		mockIgnoreFileOpener,
	)

	require.NoError(t, err)
	require.Len(t, regexpPatterns, expectedLength)
	require.Equal(t, expectedRegexpPatterns, regexpPatterns)
}

package scm_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/stretchr/testify/require"
)

func TestParseIgnorePatterns(t *testing.T) {
	ignorePatterns := []string{
		"/tmp",
		".pnp.js",
		"*.log",
		"*.log*",
		"src/pkg/test/impl/",
		"# Comment, to be ignored",
		"",
		"   /tmp",
		"   .pnp.js",
		"   *.log",
		"   src/pkg/test/impl/",
		"   # Comment, to be ignored",
	}

	buf := bytes.NewBufferString(strings.Join(ignorePatterns, "\n"))
	expected := []regexp.Regexp{
		*regexp.MustCompile(`\/tmp`),
		*regexp.MustCompile(`.pnp.js`),
		*regexp.MustCompile(`.*.log`),
		*regexp.MustCompile(`.*.log*.`),
		*regexp.MustCompile(`src/pkg/test/impl/`),
		*regexp.MustCompile(`\/tmp`),
		*regexp.MustCompile(`.pnp.js`),
		*regexp.MustCompile(`.*.log`),
		*regexp.MustCompile(`src/pkg/test/impl/`),
	}

	actual, err := scm.ParseIgnorePatterns(buf)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

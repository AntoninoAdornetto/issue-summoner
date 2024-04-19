package issue_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/stretchr/testify/require"
)

// @TODO the line numbers of this assertion are off by 1. Fix in ParseLine function of comment.go
func TestScanSingleLineCommentsGo(t *testing.T) {
	r := generateSingleLineCommentImplFileGo()
	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)
	err = im.Scan(r, "/tmp/temp-dir/mock.go")
	require.NoError(t, err)

	// see generateSingleLineCommentImplFileGo for how the file
	// looks. It contains single line comments within a mock
	// source code file and provides an example of how single
	// line comments may be used. In the mock file, we added
	// 2 single line comments with annotations.
	expected := []issue.Issue{
		{
			Title:                "add ! (not) operator support for ignoring specific files/directories",
			Description:          "",
			AnnotationLineNumber: 13,
			StartLineNumber:      13,
			EndLineNumber:        13,
			ID:                   "/tmp/temp-dir/mock.go-13",
			FilePath:             "/tmp/temp-dir/mock.go",
		},
		{
			Title:                "update the formatIgnoreExpression expression to include ! operator support",
			Description:          "",
			AnnotationLineNumber: 25,
			StartLineNumber:      25,
			EndLineNumber:        25,
			ID:                   "/tmp/temp-dir/mock.go-25",
			FilePath:             "/tmp/temp-dir/mock.go",
		},
	}

	actual := im.GetIssues()
	require.Equal(t, expected, actual)
}

// should Walk the temp project created in /tmp dir and return
// a count of the number of times that Walk calls the Scan method
func TestWalkCountScans(t *testing.T) {
	root, err := setup()
	require.NoError(t, err)

	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)

	// the setup func generates 3 files and 3 directories.
	// 3 dirs (root temp dir, .git/, and pkg/)
	// 3 files (.exe file, impl go file, INDEX file that lives in .git/)
	// the expected number of times that Scan should be called is 2 times.
	// one time for the exe file and one time for the go impl file.
	// the only reason scan is called on the exe file is because this test
	// does not add any ignore patterns to pass into Walk. The next test
	// will make the assertion with ignore patterns.
	expected := 2
	actual, err := im.Walk(root, []regexp.Regexp{})
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	err = teardown(root)
	require.NoError(t, err)
}

// should Walk the temp project created in /tmp dir and return
// a count of the number of times that Walk calls filepath.WalkDir
// but this time we will add path validation in the mix.
// regular expressions are built based off the gitignore file.
// see ignore.go for examples on we that process is handled.

// should Walk the temp project create in /tmp dir and return
// a count of the number of times that Walk calls the Scan method.
// This time, we will add ignore patterns as a argument to Walk
// and the result is that Scan is only called 1 time on an impl file.
func TestWalkCountStepsWithIgnorePatterns(t *testing.T) {
	root, err := setup()
	require.NoError(t, err)

	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)

	// the setup func generates 3 files and 3 directories.
	// 3 dirs (root temp dir, .git/, and pkg/)
	// 3 files (.exe file, impl go file, INDEX file that lives in .git/)
	// the expected number of times that Scan should be called is 1 time.
	// one time for the go impl file. We will add an ignore pattern to
	// assert that Scan is not called on the executable.
	expected := 1
	actual, err := im.Walk(root, []regexp.Regexp{*regexp.MustCompile(`.*\.exe`)})
	require.NoError(t, err)
	require.Equal(t, expected, actual)
	err = teardown(root)
	require.NoError(t, err)
}

// should return an error when the root path does not exist
func TestWalkNoneExistentRoot(t *testing.T) {
	im, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
	require.NoError(t, err)
	require.NotNil(t, im)

	_, err = im.Walk("unknown-path", []regexp.Regexp{})
	require.Error(t, err)
}

// setup will create 6 files/dirs in total.
// 1. temp dir (temp-git-project)
// 2. temp pkg dir
// 3. temp imp file (issue.go)
// 4. temp git dir
// 5. temp INDEX file that resides in the temp .git dir
// 6. temp exe file (used to validate that walk can ignore/skip entire directories)
func setup() (string, error) {
	pathName, err := os.MkdirTemp("", "temp-git-project")
	if err != nil {
		return "", err
	}

	if err = buildSrc(pathName); err != nil {
		return "", err
	}

	return pathName, err
}

// produces a go implementation file that utilizes
// single line comments with a mock annotation.
// it is used to help assert that the PendingIssue Scan
// implementation behaves as expected. There are two
// valid issue annotations in the bytes buffer returned
// from this function. There are also 2 additonal single
// line comments that should be ignored.
func generateSingleLineCommentImplFileGo() io.Reader {
	implBytes := []byte(`package scm

		import (
			"bufio"
			"bytes"
			"io"
			"regexp"
			"unicode"
		)

		type IgnorePattern = regexp.Regexp

		// @TEST_TODO add ! (not) operator support for ignoring specific files/directories
		func ParseIgnorePatterns(r io.Reader) ([]IgnorePattern, error) {
			regexps := make([]IgnorePattern, 0)
			buf := &bytes.Buffer{}
			scanner := bufio.NewScanner(r)

			for scanner.Scan() {
				line := scanner.Bytes()
				if len(line) == 0 {
					continue
				}

				// @TEST_TODO update the formatIgnoreExpression expression to include ! operator support
				if err := formatIgnoreExpression(buf, line); err != nil {
					return regexps, err
				}

				// a single line comment that should be ignored
				if buf.Len() > 0 {
					re := regexp.MustCompile(buf.String())
					regexps = append(regexps, *re)
					buf.Reset()
				}
			}

			// another single line comment that should be ignored
			return regexps, scanner.Err()
		}
	`)
	r := bytes.NewReader(implBytes)
	return r
}

// creates temp source code files for our temporary dir
// this will be used to test the Walk func.
func buildSrc(path string) error {
	tempPkg, err := os.MkdirTemp(path, "pkg")
	if err != nil {
		return err
	}

	// we add a temp .git dir to ensure that our walk func
	// does not parse the files inside a .git dir
	tempGitDir, err := os.MkdirTemp(path, ".git")
	if err != nil {
		return err
	}

	gitIndex, err := os.CreateTemp(tempGitDir, "INDEX")
	if err != nil {
		return err
	}

	gitIndexBytes := []byte("git index file")
	_, err = gitIndex.Write(gitIndexBytes)
	if err != nil {
		return err
	}

	// another ignored file. This is to test the ignore pattern
	// functionality of the walk func
	_, err = os.CreateTemp(path, ".exe")
	if err != nil {
		return err
	}

	issueGo, err := os.CreateTemp(tempPkg, "issue")
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	issueSrc, err := os.Open(filepath.Join(wd, "issue.go"))
	if err != nil {
		return err
	}

	sourceBytes, err := io.ReadAll(issueSrc)
	if err != nil {
		return err
	}

	_, err = issueGo.Write(sourceBytes)
	if err != nil {
		return err
	}

	newFileName := issueGo.Name() + ".go"
	err = os.Rename(issueGo.Name(), newFileName)
	if err != nil {
		return err
	}

	return nil
}

func teardown(path string) error {
	return os.RemoveAll(path)
}

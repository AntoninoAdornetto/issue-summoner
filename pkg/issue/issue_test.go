package issue_test

import (
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

const (
	test_annotation = "@TEST_TODO"
)

func TestNewIssueManagerNotReportingPendMode(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, issue.ISSUE_MODE_PEND, false)
	require.NoError(t, err)
	require.NotNil(t, manager)
}

func TestNewIssueManagerReportingPendMode(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, issue.ISSUE_MODE_PEND, true)
	require.NoError(t, err)
	require.NotNil(t, manager)
}

func TestNewIssue(t *testing.T) {
	testCases := []struct {
		manager  issue.IssueManager
		token    lexer.Token
		comment  lexer.Comment
		expected issue.Issue
	}{
		{
			manager: issue.IssueManager{
				CurrentPath: "/project/test-project/src/module.c",
				CurrentBase: "module.c",
			},
			token: lexer.Token{
				TokenType:      lexer.SINGLE_LINE_COMMENT,
				Lexeme:         []byte("// @TEST_TODO test me"),
				Line:           15,
				StartByteIndex: 100,
				EndByteIndex:   120,
			},
			comment: lexer.Comment{
				Title:       []byte("test me"),
				Description: []byte(""),
			},
			expected: issue.Issue{
				ID:           "/project/test-project/src/module.c-100:120",
				Title:        "test me",
				IssueIndex:   0,
				StartIndex:   100,
				EndIndex:     120,
				LineNumber:   15,
				FilePath:     "/project/test-project/src/module.c",
				FileName:     "module.c",
				SubmissionID: -1,
				Index:        0,
			},
		},
		{
			manager: issue.IssueManager{
				CurrentPath: "/project/test-project/src/next-module.c",
				CurrentBase: "next-module.c",
				RecordCount: 1,
				Issues:      []issue.Issue{{Title: "1st issue"}},
			},
			token: lexer.Token{
				TokenType: lexer.MULTI_LINE_COMMENT,
				Lexeme: []byte(
					"/* @TEST_TODO write more tests\n\t* to ensure correct functionality\n\t* for our users\n\t*/",
				),
				Line:           50,
				StartByteIndex: 300,
				EndByteIndex:   390,
			},
			comment: lexer.Comment{
				Title:       []byte("write more tests"),
				Description: []byte("to ensure correct functionality for our users"),
			},
			expected: issue.Issue{
				ID:           "/project/test-project/src/next-module.c-300:390",
				Title:        "write more tests",
				Description:  "to ensure correct functionality for our users",
				IssueIndex:   1,
				StartIndex:   300,
				EndIndex:     390,
				LineNumber:   50,
				FilePath:     "/project/test-project/src/next-module.c",
				FileName:     "next-module.c",
				SubmissionID: -1,
				Index:        1,
			},
		},
	}

	for _, tc := range testCases {
		issue, err := tc.manager.NewIssue(tc.comment, tc.token)
		require.NoError(t, err)
		require.Equal(t, tc.expected, issue)
	}
}

func TestNewIssueWithTemplate(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, issue.ISSUE_MODE_PEND, true)
	require.NoError(t, err)
	require.NotNil(t, manager)

	comment := lexer.Comment{
		Title:       []byte("New Issue for github"),
		Description: []byte("found a bug"),
	}

	token := lexer.Token{
		TokenType:      lexer.MULTI_LINE_COMMENT,
		Lexeme:         []byte("/* New Issue for github\n\t* found a bug\n\t*/"),
		Line:           15,
		StartByteIndex: 100,
		EndByteIndex:   120,
	}

	/*
		Normally the ID, FilePath, FileName properties will have proper values.
		They do not populate during this test because we have not scanned the source files
	*/
	expected := issue.Issue{
		ID:           "-100:120",
		Title:        string(comment.Title),
		Description:  string(comment.Description),
		OS:           "linux",
		FilePath:     "",
		FileName:     "",
		LineNumber:   token.Line,
		IssueIndex:   0,
		StartIndex:   token.StartByteIndex,
		EndIndex:     token.EndByteIndex,
		SubmissionID: -1,
		Body:         "### Description\nfound a bug\n\n### Location\n\n***File name:***  ***Line number:*** 15\n\n### Environment\n\nlinux\n\n### Generated with :heart:\n\ncreated by [issue-summoner](https://github.com/AntoninoAdornetto/issue-summoner)\n\t",
	}

	actual, err := manager.NewIssue(comment, token)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestScanPendMode(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, issue.ISSUE_MODE_PEND, false)
	require.NoError(t, err)
	require.NotNil(t, manager)

	require.NoError(t, err)

	// actual implementation will use abs paths and not relative
	manager.CurrentPath = "../../testdata/test.c"
	manager.CurrentBase = "test.c"

	err = manager.Scan("../../testdata/test.c")
	require.NoError(t, err)

	os := runtime.GOOS

	// see test.c in the testdata directory for how the file looks prior
	// to scanning. The file contains c source code with comments that
	// are parsed. All comments that are annotated with @TEST_TODO should
	// appear in the expected slice. Also, wanted to mention that the
	// paths shown in ID & FilePath are not how they normally look. I am
	// using relative paths for tests so that you guys can't see the name
	// of my directories :)
	expectedIssues := []issue.Issue{
		{
			ID:           "../../testdata/test.c-62:95",
			Title:        "inline comment #1",
			Description:  "",
			LineNumber:   5,
			FileName:     "test.c",
			FilePath:     "../../testdata/test.c",
			StartIndex:   62,
			EndIndex:     95,
			IssueIndex:   0,
			OS:           os,
			Index:        0,
			SubmissionID: -1,
		},
		{
			ID:           "../../testdata/test.c-115:148",
			Title:        "inline comment #2",
			Description:  "",
			LineNumber:   6,
			FileName:     "test.c",
			FilePath:     "../../testdata/test.c",
			StartIndex:   115,
			EndIndex:     148,
			IssueIndex:   1,
			OS:           os,
			Index:        1,
			SubmissionID: -1,
		},
		{
			ID:           "../../testdata/test.c-192:252",
			Title:        "decode the message and clean up after yourself!",
			Description:  "",
			FileName:     "test.c",
			FilePath:     "../../testdata/test.c",
			LineNumber:   10,
			StartIndex:   192,
			EndIndex:     252,
			IssueIndex:   2,
			OS:           os,
			Index:        2,
			SubmissionID: -1,
		},
		{
			ID:           "../../testdata/test.c-269:561",
			Title:        "drop a star if you know about this code wars challenge",
			Description:  "Digital Cypher assigns to each letter of the alphabet unique number. Instead of letters in encrypted word we write the corresponding number Then we add to each obtained digit consecutive digits from the key",
			FileName:     "test.c",
			FilePath:     "../../testdata/test.c",
			LineNumber:   19,
			StartIndex:   269,
			EndIndex:     561,
			IssueIndex:   3,
			OS:           os,
			Index:        3,
			SubmissionID: -1,
		},
	}

	expectedIssueManager := &issue.IssueManager{
		CurrentPath: "../../testdata/test.c",
		CurrentBase: "test.c",
		RecordCount: len(expectedIssues),
	}

	require.Equal(t, expectedIssues, manager.Issues)
	require.Equal(t, expectedIssueManager.ReportMap, manager.ReportMap)
	require.Equal(t, expectedIssueManager.CurrentBase, manager.CurrentBase)
	require.Equal(t, expectedIssueManager.CurrentPath, manager.CurrentPath)
	require.Equal(t, expectedIssueManager.RecordCount, manager.RecordCount)
}

func TestScanPendModeWithReporting(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, issue.ISSUE_MODE_PEND, true)
	require.NoError(t, err)
	require.NotNil(t, manager)

	require.NoError(t, err)

	// actual implementation will use abs paths and not relative
	manager.CurrentPath = "../../testdata/test.c"
	manager.CurrentBase = "test.c"

	err = manager.Scan("../../testdata/test.c")
	require.NoError(t, err)

	os := runtime.GOOS

	// see test.c in the testdata directory for how the file looks prior
	// to scanning. The file contains c source code with comments that
	// are parsed. All comments that are annotated with @TEST_TODO should
	// appear in the expected slice. Also, wanted to mention that the
	// paths shown in ID & FilePath are not how they normally look. I am
	// using relative paths for tests so that you guys can't see the name
	// of my directories :)
	expectedIssues := []issue.Issue{
		{
			ID:          "../../testdata/test.c-62:95",
			Title:       "inline comment #1",
			Description: "",
			LineNumber:  5,
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			StartIndex:  62,
			EndIndex:    95,
			IssueIndex:  0,
			OS:          os,
		},
		{
			ID:          "../../testdata/test.c-115:148",
			Title:       "inline comment #2",
			Description: "",
			LineNumber:  6,
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			StartIndex:  115,
			EndIndex:    148,
			IssueIndex:  1,
			OS:          os,
		},
		{
			ID:          "../../testdata/test.c-192:252",
			Title:       "decode the message and clean up after yourself!",
			Description: "",
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			LineNumber:  10,
			StartIndex:  192,
			EndIndex:    252,
			IssueIndex:  2,
			OS:          os,
		},
		{
			ID:          "../../testdata/test.c-269:561",
			Title:       "drop a star if you know about this code wars challenge",
			Description: "Digital Cypher assigns to each letter of the alphabet unique number. Instead of letters in encrypted word we write the corresponding number Then we add to each obtained digit consecutive digits from the key",
			FileName:    "test.c",
			FilePath:    "../../testdata/test.c",
			LineNumber:  19,
			StartIndex:  269,
			EndIndex:    561,
			IssueIndex:  3,
			OS:          os,
		},
	}

	expectedIssueManager := &issue.IssueManager{
		CurrentPath: "../../testdata/test.c",
		CurrentBase: "test.c",
		RecordCount: len(expectedIssues),
	}

	require.Equal(t, expectedIssueManager.CurrentBase, manager.CurrentBase)
	require.Equal(t, expectedIssueManager.CurrentPath, manager.CurrentPath)
	require.Len(t, manager.ReportMap["../../testdata/test.c"], len(expectedIssues))
	require.Equal(t, expectedIssueManager.RecordCount, manager.RecordCount)
}

func TestConsolidateMap(t *testing.T) {
	issueManager := issue.IssueManager{
		ReportMap: map[string][]*issue.Issue{
			"../../testdata/test.c": {
				&issue.Issue{
					Title:        "Submission 1",
					SubmissionID: -1,
					ID:           "../../testdata/test.c:50-60",
					Index:        0,
				},
				&issue.Issue{
					Title:        "Submission 2",
					SubmissionID: -1,
					ID:           "../../testdata/test.c:84-90",
					Index:        1,
				},
				&issue.Issue{
					Title:        "Submission 3",
					SubmissionID: 8912,
					ID:           "../../testdata/test.c:200-300",
					Index:        2,
				},
				&issue.Issue{Title: "Submission 4", SubmissionID: -1},
				&issue.Issue{
					Title:        "Submission 5",
					SubmissionID: 83234,
					ID:           "../../testdata/test.c:400-500",
					Index:        3,
				},
				&issue.Issue{
					Title:        "Submission 6",
					SubmissionID: -1,
					Index:        3,
				},
				&issue.Issue{
					Title:        "Submission 7",
					SubmissionID: 100,
					Index:        4,
				},
			},
		},
	}

	issueManager.ConsolidateMap()

	expected := map[string][]*issue.Issue{
		"../../testdata/test.c": {
			&issue.Issue{
				Title:        "Submission 3",
				SubmissionID: 8912,
				ID:           "../../testdata/test.c:200-300",
				Index:        2,
			},
			&issue.Issue{
				Title:        "Submission 5",
				SubmissionID: 83234,
				ID:           "../../testdata/test.c:400-500",
				Index:        3,
			},
			&issue.Issue{Title: "Submission 7", SubmissionID: 100, Index: 4},
		},
	}

	for key, issues := range issueManager.ReportMap {
		for i, actual := range issues {
			require.Equal(t, expected[key][i], actual)
		}
	}
}

func TestWriteIssueIDs(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, issue.ISSUE_MODE_PEND, true)
	require.NoError(t, err)
	require.NotNil(t, manager)

	require.NoError(t, err)

	// actual implementation will use abs paths and not relative
	manager.CurrentPath = "../../testdata/test.c"
	manager.CurrentBase = "test.c"

	err = manager.Scan("../../testdata/test.c")
	require.NoError(t, err)

	file, err := os.OpenFile("../../testdata/test.c", os.O_RDWR, 0666)
	require.NoError(t, err)

	originalData, err := io.ReadAll(file)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	// restore original test file contents
	defer func() {
		file, err := os.OpenFile("../../testdata/test.c", os.O_RDWR, 0666)
		require.NoError(t, err)
		_, err = file.Write(originalData)
		require.NoError(t, err)
	}()

	for i, issue := range manager.ReportMap["../../testdata/test.c"] {
		issue.SubmissionID = int64(i + 1)
		// simulate a report request and update submission ids that will be written to the src file
		// normally the submissionID is updated after reporting to a src code manager but we want
		// our test to be deterministic
	}

	err = manager.WriteIssueIDs("../../testdata/test.c")
	require.NoError(t, err)

	expectedSourceCode := `#include <stdio.h>
#include <stdlib.h>

struct Person {
  int /* @TEST_TODO(#1) inline comment #1 */ age;
  char *name /* @TEST_TODO(#2) inline comment #2 */;
};

int main(int argc, char *argv[]) {
  // @TEST_TODO(#3) decode the message and clean up after yourself!
  return 0;
}

/*
 * @TEST_TODO(#4) drop a star if you know about this code wars challenge
 * Digital Cypher assigns to each letter of the alphabet unique number.
 * Instead of letters in encrypted word we write the corresponding number
 * Then we add to each obtained digit consecutive digits from the key
 * */
char *decode(const unsigned char *code, size_t n, unsigned key) {
  char *msg = calloc(n + 1, 1);
  char buf[64];
  int key_len = sprintf(buf, "%d", key);

  for (size_t i = 0; i < n; i++) {
    msg[i] = code[i] - buf[i % key_len] + '0' + 'a' - 1;
  }

  return msg;
}

// This comment should not be parsed since it does not have an annotation
	`

	newVersion, err := os.Open("../../testdata/test.c")
	require.NoError(t, err)

	newSrcCode, err := io.ReadAll(newVersion)
	require.NoError(t, err)
	require.NoError(t, newVersion.Close())
	require.Equal(t, expectedSourceCode, string(newSrcCode))
}

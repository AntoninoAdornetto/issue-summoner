package issue_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

const (
	test_annotation = "@TEST_TODO"
)

func TestNewIssueManagerNotReporting(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, false)
	require.NoError(t, err)
	require.NotNil(t, manager)
}

func TestNewIssueManagerReporting(t *testing.T) {
	manager, err := issue.NewIssueManager(test_annotation, true)
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
				ID:         "/project/test-project/src/module.c-100:120",
				Title:      "test me",
				IssueIndex: 0,
				StartIndex: 100,
				EndIndex:   120,
				LineNumber: 15,
				FilePath:   "/project/test-project/src/module.c",
				FileName:   "module.c",
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
				ID:          "/project/test-project/src/next-module.c-300:390",
				Title:       "write more tests",
				Description: "to ensure correct functionality for our users",
				IssueIndex:  1,
				StartIndex:  300,
				EndIndex:    390,
				LineNumber:  50,
				FilePath:    "/project/test-project/src/next-module.c",
				FileName:    "next-module.c",
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
	manager, err := issue.NewIssueManager(test_annotation, true)
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
		ID:          "-100:120",
		Title:       string(comment.Title),
		Description: string(comment.Description),
		OS:          "linux",
		FilePath:    "",
		FileName:    "",
		LineNumber:  token.Line,
		IssueIndex:  0,
		StartIndex:  token.StartByteIndex,
		EndIndex:    token.EndByteIndex,
		Body:        "### Description\nfound a bug\n\n### Location\n\n***File name:***  ***Line number:*** 15\n\n### Environment\n\nlinux\n\n### Generated with :heart:\n\ncreated by [issue-summoner](https://github.com/AntoninoAdornetto/issue-summoner)\n\t",
	}

	actual, err := manager.NewIssue(comment, token)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

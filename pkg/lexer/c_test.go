package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeTokenSingleLineComments(t *testing.T) {
	testCases := []struct {
		name       string
		srcCode    []byte
		fileName   string
		expected   []lexer.Token
		flags      lexer.U8
		annotation []byte
	}{
		{
			name: "should not create any tokens when consuming a non-comment notation byte",
			// int is not "//" or "/*" - denotes the opening notation of a comment in c-like languages
			srcCode:    []byte("int"),
			fileName:   "main.c",
			expected:   []lexer.Token{},
			annotation: testAnnotation,
			flags:      lexer.FLAG_SCAN,
		},
		{
			name:       "should not create any tokens when consuming single line comment bytes that do not have an issue annotation",
			srcCode:    []byte("// regular single line comment with no issue annotation in c"),
			fileName:   "main.c",
			expected:   []lexer.Token{},
			annotation: testAnnotation,
			flags:      lexer.FLAG_SCAN,
		},
		{
			name:       "should create the comment start, comment annotation and comment end tokens",
			srcCode:    []byte("// @TEST_ANNOTATION\n"),
			fileName:   "main.c",
			flags:      lexer.FLAG_SCAN,
			annotation: testAnnotation,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Start:  0,
					End:    1,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  3,
					End:    18,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{'\n'},
					Start:  19,
					End:    19,
					Line:   2,
				},
			},
		},
		{
			name: "should create the comment start, comment annotation and comment end tokens when more than 2 forward slashes are used to denote a SL comment",
			// the key take away here is that the opening SL comment notation has 4 forward slashes
			srcCode:    []byte("//// @TEST_ANNOTATION\n"),
			fileName:   "main.c",
			flags:      lexer.FLAG_SCAN,
			annotation: testAnnotation,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("////"),
					Start:  0,
					End:    3,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  5,
					End:    20,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{'\n'},
					Start:  21,
					End:    21,
					Line:   2,
				},
			},
		},
		{
			name:       "should create the comment start, comment annotation, comment title and comment end tokens",
			srcCode:    []byte("// @TEST_ANNOTATION check for edge cases"),
			fileName:   "main.c",
			flags:      lexer.FLAG_SCAN,
			annotation: testAnnotation,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Start:  0,
					End:    1,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  3,
					End:    18,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("check"),
					Start:  20,
					End:    24,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("for"),
					Start:  26,
					End:    28,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("edge"),
					Start:  30,
					End:    33,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("cases"),
					Start:  35,
					End:    39,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Start:  40,
					End:    40,
					Line:   1,
				},
			},
		},
		{
			name: "should create the comment start, comment annotation, comment title and comment end tokens when given the purge flag",
			// the src code for this test simulates an issue annotation that has been reported already. I.E., issue was published to github,gitlab,ect
			srcCode:  []byte("// @TEST_ANNOTATION(#98321) check for edge cases"),
			fileName: "main.c",
			flags:    lexer.FLAG_PURGE,
			// this pattern is constructed within the issue package
			annotation: []byte("@TEST_ANNOTATION\\(#\\d+\\)"),
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Start:  0,
					End:    1,
					Line:   1,
					Lexeme: []byte("//"),
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Start:  3,
					End:    18,
					Line:   1,
					Lexeme: []byte("@TEST_ANNOTATION"),
				},
				{
					Type:   lexer.TOKEN_OPEN_PARAN,
					Start:  19,
					End:    19,
					Line:   1,
					Lexeme: []byte{lexer.OPEN_PARAN},
				},
				{
					Type:   lexer.TOKEN_HASH,
					Start:  20,
					End:    20,
					Line:   1,
					Lexeme: []byte{lexer.HASH},
				},
				{
					Type:   lexer.TOKEN_ISSUE_NUMBER,
					Start:  21,
					End:    25,
					Line:   1,
					Lexeme: []byte("98321"),
				},
				{
					Type:   lexer.TOKEN_CLOSE_PARAN,
					Start:  26,
					End:    26,
					Lexeme: []byte{lexer.CLOSE_PARAN},
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Start:  28,
					End:    32,
					Line:   1,
					Lexeme: []byte("check"),
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Start:  34,
					End:    36,
					Line:   1,
					Lexeme: []byte("for"),
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Start:  38,
					End:    41,
					Line:   1,
					Lexeme: []byte("edge"),
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Start:  43,
					End:    47,
					Line:   1,
					Lexeme: []byte("cases"),
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Start:  48,
					End:    48,
					Line:   1,
					Lexeme: []byte{0},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseLexer := lexer.NewLexer(tc.annotation, tc.srcCode, tc.fileName, tc.flags)
			cLexer := lexer.Clexer{Base: baseLexer}
			err := cLexer.AnalyzeToken()
			require.NoError(t, err)
			require.Equal(t, tc.expected, baseLexer.Tokens)
		})
	}
}

func TestAnalyzeTokenMultiLineComments(t *testing.T) {
	testCases := []struct {
		name     string
		srcCode  []byte
		fileName string
		expected []lexer.Token
	}{
		{
			name: "should not create any tokens when consuming a non-comment notation byte",
			// int is not "//" or "/*" - denotes the opening notation of a comment in c-like languages
			srcCode:  []byte("int"),
			fileName: "main.c",
			expected: []lexer.Token{},
		},
		{
			name:     "should not create any tokens when consuming multi line comment bytes that do not have an issue annotation",
			srcCode:  []byte("/* regular single line comment with no issue annotation in c */"),
			fileName: "main.c",
			expected: []lexer.Token{},
		},
		{
			name:     "should create the comment start, comment annotation and comment end tokens",
			srcCode:  []byte("/* @TEST_ANNOTATION */"),
			fileName: "main.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Start:  0,
					End:    1,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  3,
					End:    18,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Start:  20,
					End:    21,
					Line:   1,
				},
			},
		},
		{
			name:     "should create the comment start, comment annotation and comment end tokens when the opening ML comment notation has extra chars in it",
			srcCode:  []byte("/*** @TEST_ANNOTATION */"),
			fileName: "main.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/***"),
					Start:  0,
					End:    3,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  5,
					End:    20,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Start:  22,
					End:    23,
					Line:   1,
				},
			},
		},
		{
			name: "should create comment start, annotation, title, and comment end tokens for multi line comment",
			// for this assertion we are using multi line comment notation but have not added any line breaks
			// the next test will assert with line breaks
			srcCode:  []byte("/* @TEST_ANNOTATION multi line comment */"),
			fileName: "main.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Start:  0,
					End:    1,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  3,
					End:    18,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("multi"),
					Start:  20,
					End:    24,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("line"),
					Start:  26,
					End:    29,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Start:  31,
					End:    37,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Start:  39,
					End:    40,
					Line:   1,
				},
			},
		},
		{
			name: "should create comment start, annotation, title, description, and comment end tokens for multi line comment",
			// here we introduce line breaks. Line breaks will build comment description tokens.
			srcCode: []byte(
				"/*\n\t* @TEST_ANNOTATION comment title\n\t* comment description 1\n\t* comment description 2\n*/",
			),
			fileName: "main.c",
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_START,
					Lexeme: []byte("/*"),
					Start:  0,
					End:    1,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  6,
					End:    21,
					Line:   2,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Start:  23,
					End:    29,
					Line:   2,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("title"),
					Start:  31,
					End:    35,
					Line:   2,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("comment"),
					Start:  40,
					End:    46,
					Line:   3,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("description"),
					Start:  48,
					End:    58,
					Line:   3,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("1"),
					Start:  60,
					End:    60,
					Line:   3,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("comment"),
					Start:  65,
					End:    71,
					Line:   4,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("description"),
					Start:  73,
					End:    83,
					Line:   4,
				},
				{
					Type:   lexer.TOKEN_COMMENT_DESCRIPTION,
					Lexeme: []byte("2"),
					Start:  85,
					End:    85,
					Line:   4,
				},
				{
					Type:   lexer.TOKEN_MULTI_LINE_COMMENT_END,
					Lexeme: []byte("*/"),
					Start:  87,
					End:    88,
					Line:   5,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseLexer := lexer.NewLexer(testAnnotation, tc.srcCode, tc.fileName, lexer.FLAG_SCAN)
			cLexer := lexer.Clexer{Base: baseLexer}
			err := cLexer.AnalyzeToken()
			require.NoError(t, err)
			require.Equal(t, tc.expected, baseLexer.Tokens)
		})
	}
}

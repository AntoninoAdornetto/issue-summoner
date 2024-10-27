package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

func TestPyAnalyzeTokenSingleLineComments(t *testing.T) {
	testCases := []struct {
		name       string
		srcCode    []byte
		fileName   string
		expected   []lexer.Token
		flags      lexer.U8
		annotation []byte
	}{
		{
			name:       "should not create any tokens when consuming bytes that do not contain comments",
			srcCode:    []byte("print('Hello, World!')"),
			fileName:   "main.py",
			expected:   []lexer.Token{},
			annotation: testAnnotation,
			flags:      lexer.FLAG_SCAN,
		},
		{
			name:       "should not create any tokens when consuming bytes for a single line comment that is not annotated",
			srcCode:    []byte("# single line comment with no issue annotation"),
			fileName:   "main.py",
			expected:   []lexer.Token{},
			annotation: testAnnotation,
			flags:      lexer.FLAG_SCAN,
		},
		{
			name:       "should create the comment start, comment annotation, and comment end tokens for an annotated SL comment",
			srcCode:    []byte("# @TEST_ANNOTATION\n"),
			fileName:   "main.py",
			flags:      lexer.FLAG_SCAN,
			annotation: testAnnotation,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte{'#'},
					Start:  0,
					End:    0,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_ANNOTATION"),
					Start:  2,
					End:    17,
					Line:   1,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{'\n'},
					Start:  18,
					End:    18,
					Line:   2,
				},
			},
		},
		{
			name:       "should create the comment start, comment annotation, and comment end tokens for an annotated SL comment with multiple leading comment notation bytes",
			srcCode:    []byte("#### @TEST_ANNOTATION\n"),
			fileName:   "main.py",
			flags:      lexer.FLAG_SCAN,
			annotation: testAnnotation,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("####"),
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
			srcCode:    []byte("## @TEST_ANNOTATION check for edge cases"),
			fileName:   "main.py",
			flags:      lexer.FLAG_SCAN,
			annotation: testAnnotation,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("##"),
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseLexer := lexer.NewLexer(tc.annotation, tc.srcCode, tc.fileName, tc.flags)
			pyLexer := lexer.PyLexer{Base: baseLexer}
			err := pyLexer.AnalyzeToken()
			require.NoError(t, err)
			require.Equal(t, tc.expected, baseLexer.Tokens)
		})
	}
}

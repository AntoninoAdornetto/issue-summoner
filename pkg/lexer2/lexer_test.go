package lexer2_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer2"
	"github.com/stretchr/testify/require"
)

func TestCLexicalAnaylzer(t *testing.T) {
	testCases := []struct {
		name             string
		testDataFilePath string
		expected         []lexer2.Token
	}{
		{
			name:             "Should return all tokens that make up the 1 single line comment in the singleline.c src code file",
			testDataFilePath: "../../testdata/singleline.c",
			expected: []lexer2.Token{
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   4,
					Start:  58,
					End:    59,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   4,
					Start:  61,
					End:    70,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("simple"),
					Line:   4,
					Start:  72,
					End:    77,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("single"),
					Line:   4,
					Start:  79,
					End:    84,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("line"),
					Line:   4,
					Start:  86,
					End:    89,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   4,
					Start:  91,
					End:    97,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Line:   4,
					Start:  97,
					End:    97,
				},
				{
					Type:   lexer2.TOKEN_EOF,
					Lexeme: []byte{},
					Line:   6,
					Start:  123,
					End:    123,
				},
			},
		},
		{
			name:             "Should return all tokens that make up the multiple single line comment in the singleline-multi.c src code file and ignore the non annotated comment",
			testDataFilePath: "../../testdata/singleline-multi.c",
			expected: []lexer2.Token{
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   3,
					Start:  21,
					End:    22,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   3,
					Start:  24,
					End:    33,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("implement"),
					Line:   3,
					Start:  35,
					End:    43,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("the"),
					Line:   3,
					Start:  45,
					End:    47,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("main"),
					Line:   3,
					Start:  49,
					End:    52,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("function"),
					Line:   3,
					Start:  54,
					End:    61,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Line:   3,
					Start:  61,
					End:    61,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte("//"),
					Line:   5,
					Start:  100,
					End:    101,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_ANNOTATION,
					Lexeme: []byte("@TEST_TODO"),
					Line:   5,
					Start:  103,
					End:    112,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("simple"),
					Line:   5,
					Start:  114,
					End:    119,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("single"),
					Line:   5,
					Start:  121,
					End:    126,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("line"),
					Line:   5,
					Start:  128,
					End:    131,
				},
				{
					Type:   lexer2.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("comment"),
					Line:   5,
					Start:  133,
					End:    139,
				},
				{
					Type:   lexer2.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{0},
					Line:   5,
					Start:  139,
					End:    139,
				},
				{
					Type:   lexer2.TOKEN_EOF,
					Lexeme: []byte{},
					Line:   10,
					Start:  246,
					End:    246,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(inner *testing.T) {
			src, fileName := readSrcFileC(inner, tc.testDataFilePath)
			base := lexer2.NewBaseLexer([]byte("@TEST_TODO"), src, fileName)
			cLexer, err := lexer2.NewLexicalAnalyzer(base)
			require.NoError(inner, err)
			tokens, err := base.AnalyzeTokens(cLexer)
			require.NoError(inner, err)

			if i > 0 {
				for _, token := range tokens {
					fmt.Println(string(token.Lexeme))
					fmt.Println(string(src[token.Start : token.End+1]))
					fmt.Printf("\n\n")
				}
			}
			require.Equal(inner, tc.expected, tokens)
		})
	}
}

func readSrcFileC(t *testing.T, path string) ([]byte, string) {
	fileName := filepath.Base(path)
	srcFile, err := os.Open(path)
	require.NoError(t, err)

	defer srcFile.Close()
	srcCode, err := io.ReadAll(srcFile)
	require.NoError(t, err)
	return srcCode, fileName
}

func getTestDataSrc(t *testing.T, path string) ([]byte, string) {
	fileName := filepath.Base(path)
	srcFile, err := os.Open(path)
	require.NoError(t, err)

	defer srcFile.Close()
	srcCode, err := io.ReadAll(srcFile)
	require.NoError(t, err)
	return srcCode, fileName
}

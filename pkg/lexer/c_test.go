package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

const (
	// 2 string tokens and 1 EOF token
	// token[0] Lexeme = "Hello World!"
	// token[1] Lexeme = "%s\\n" <- first arg to printf
	// token[2] Lexeme = "" <- End of file token
	c_src_code_str = `
	#include <stdio.h>
	int main() {
		char str[] = "Hello World!";
		printf("%s\n", str);
		return 0;
	}
	`

	c_src_code_single_line_comments = `
	#include <stdio.h>
	int main() {
		int x = 0; // @TEST_TODO first single line comment
		int y = 0; // @TEST_TODO second single line comment
		
		printf("X: %d\tY: %d", x, y);
		// @TEST_TODO third single line comment
	}
	`

	// multi line comments have been tricky and there are some edge cases
	// the main one being where multi line comment can be denoted between
	// source code. The Coords struct is a good example of this and has actually
	// been a scenario that has broke the program with past implementations.
	// we want to thoroughly test multi line comment parsing to ensure it's accuracy
	c_src_code_multi_line_comment = `
	#include <stdio.h>

	typedef struct {
		int /* @TEST_TODO inline 1 */ x /* @TEST_TODO inline 2 */;
	} Coords;

	/*
	 * @TEST_TODO multi line comment
	 * second line
	 * third line
	 * end line
	*/
	int main() {
		int x = 0;
		int y = 0;
		return x + y;
	}
	`

	// 2 string tokens and 1 EOF token (using single quotes and back tick delimiters)
	// token[0] Lexeme = 'Hello World!'
	// token[1] Lexeme = `${str} printed` <- first arg to console.log
	// token[2] Lexeme = "" <- End of file token
	js_src_code_str = "\nfunction main(){\n\tconst str = 'Hello World!';\nconsole.log(`${str} printed`)}"
)

// should return a slice of tokens containing all the string values
// that are contained in the src bytes passed into the lexer (C Code)
func TestAnalyzeTokenStrC(t *testing.T) {
	lex, err := lexer.NewLexer([]byte(c_src_code_str), "main.c")
	require.NoError(t, err)
	tokens, err := lex.AnalyzeTokens()
	require.NoError(t, err)

	expectedTokens := []lexer.Token{
		{
			TokenType:      lexer.STRING,
			Lexeme:         []byte("Hello World!"),
			Line:           4,
			StartByteIndex: 50,
			EndByteIndex:   63,
		},
		{
			TokenType:      lexer.STRING,
			Lexeme:         []byte("%s\\n"), // the string passed into printf
			Line:           5,
			StartByteIndex: 75,
			EndByteIndex:   80,
		},
		{
			TokenType: lexer.EOF,
		},
	}
	require.Equal(t, expectedTokens, tokens)
}

// should return a slice of tokens containing all the string values
// that are contained in the src bytes passed into the lexer (JS Code)
// languages such as JS can declare strings using double quotes, single quotes
// and string template literals. see js_src_code_str for the full example src code
func TestAnalyzeTokenStrJS(t *testing.T) {
	lex, err := lexer.NewLexer([]byte(js_src_code_str), "index.js")
	require.NoError(t, err)
	tokens, err := lex.AnalyzeTokens()
	require.NoError(t, err)

	expectedTokens := []lexer.Token{
		{
			TokenType:      lexer.STRING,
			Lexeme:         []byte("Hello World!"),
			Line:           3,
			StartByteIndex: 31,
			EndByteIndex:   44,
		},
		{
			TokenType: lexer.STRING,
			// the template literal passed into console.log
			Lexeme:         []byte("${str} printed"),
			Line:           4,
			StartByteIndex: 59,
			EndByteIndex:   74,
		},
		{
			TokenType: lexer.EOF,
		},
	}
	require.Equal(t, expectedTokens, tokens)
}

// should return an error when encountering a string that has not been terminated
// notice we return an error where the error occured, like a good CLI citizen should do
func TestAnalyzeTokenStrCUnterminatedStr(t *testing.T) {
	src := []byte("int x = 0;\nchar str[] = \"Unterminated String")
	lex, err := lexer.NewLexer(src, "main.c")
	require.NoError(t, err)
	tokens, err := lex.AnalyzeTokens()
	require.Error(t, err)
	require.ErrorContains(
		t,
		err,
		"[main.c line 2]: Error: unterminated string: \"Unterminated String",
	)
	require.Empty(t, tokens)
}

// should locate all single line comment tokens and store the comment contents
// as the lexeme for each item in the token slice output
func TestAnalyzeTokenSingleLineCommentC(t *testing.T) {
	lex, err := lexer.NewLexer([]byte(c_src_code_single_line_comments), "main.c")
	require.NoError(t, err)
	actualTokens, err := lex.AnalyzeTokens()
	require.NoError(t, err)

	expectedTokens := []lexer.Token{
		{
			TokenType:      lexer.SINGLE_LINE_COMMENT,
			Lexeme:         []byte("// @TEST_TODO first single line comment"),
			Line:           4,
			StartByteIndex: 48,
			EndByteIndex:   86,
		},
		{
			TokenType:      lexer.SINGLE_LINE_COMMENT,
			Lexeme:         []byte("// @TEST_TODO second single line comment"),
			Line:           5,
			StartByteIndex: 101,
			EndByteIndex:   140,
		},
		{
			TokenType:      lexer.STRING,
			Lexeme:         []byte("X: %d\\tY: %d"),
			Line:           7,
			StartByteIndex: 154,
			EndByteIndex:   167,
		},
		{
			TokenType:      lexer.SINGLE_LINE_COMMENT,
			Lexeme:         []byte("// @TEST_TODO third single line comment"),
			Line:           8,
			StartByteIndex: 179,
			EndByteIndex:   217,
		},
		{
			TokenType: lexer.EOF,
			Lexeme:    []byte(nil),
		},
	}
	require.Equal(t, expectedTokens, actualTokens)

	// we also want to confirm that the byte indices (StartByteIndex, EndByteIndex)
	// point to the expected byte values. I.E. StartByteIndex of a single line comment
	// should point to the byte of '/' and the end byte index should point to the last
	// byte in the comment, not the new line byte

	// 47 = expectedTokens[0].StartByteIndex
	// 85 = expectedTokens[0].EndByteIndex
	require.Equal(t, byte('/'), c_src_code_single_line_comments[48])
	require.Equal(t, byte('t'), c_src_code_single_line_comments[86])

	// 100 = expectedTokens[1].StartByteIndex
	// 139 = expectedTokens[1].EndByteIndex
	require.Equal(t, byte('/'), c_src_code_single_line_comments[101])
	require.Equal(t, byte('t'), c_src_code_single_line_comments[140])

	// 153 = expectedTokens[2].StartByteIndex
	// 166 = expectedTokens[2].EndByteIndex
	// string token
	require.Equal(t, byte('"'), c_src_code_single_line_comments[154])
	require.Equal(t, byte('"'), c_src_code_single_line_comments[167])
}

// should locate all multi line comment tokens and store the comment contents
// as the lexeme for each item in the token slice output
func TestAnalyzeTokenMultiLineCommentC(t *testing.T) {
	lex, err := lexer.NewLexer([]byte(c_src_code_multi_line_comment), "main.c")
	require.NoError(t, err)
	actualTokens, err := lex.AnalyzeTokens()
	require.NoError(t, err)

	// the first two tokens are the most important as they have been problamatic in the past
	expectedTokens := []lexer.Token{
		{
			TokenType:      lexer.MULTI_LINE_COMMENT,
			Lexeme:         []byte("/* @TEST_TODO inline 1 */"),
			Line:           5,
			StartByteIndex: 46,
			EndByteIndex:   70,
		},
		{
			TokenType:      lexer.MULTI_LINE_COMMENT,
			Lexeme:         []byte("/* @TEST_TODO inline 2 */"),
			Line:           5,
			StartByteIndex: 74,
			EndByteIndex:   98,
		},
		{
			TokenType: lexer.MULTI_LINE_COMMENT,
			Lexeme: []byte(
				"/*\n\t * @TEST_TODO multi line comment\n\t * second line\n\t * third line\n\t * end line\n\t*/",
			),
			Line:           13,
			StartByteIndex: 114,
			EndByteIndex:   197,
		},
		{
			TokenType: lexer.EOF,
			Lexeme:    []byte(nil),
		},
	}

	require.Equal(t, expectedTokens, actualTokens)

	// we also want to confirm that the byte indices (StartByteIndex, EndByteIndex)
	// point to the expected byte values. I.E. StartByteIndex of a multi line comment
	// should point to the byte of '/' and same with the EndByteIndex. We can also check
	// the index right after the start and the index right before the end to confirm the
	// multi line comment notation.

	// 46 = expectedTokens[0].StartByteIndex
	// 47 = closing byte for the start of a multi line comment
	// 69 = expectedTokens[0].EndByteIndex
	// 70 = closing byte for the end of a multi line comment
	require.Equal(t, string(byte('/')), string(c_src_code_multi_line_comment[46]))
	require.Equal(t, string(byte('*')), string(c_src_code_multi_line_comment[47]))
	require.Equal(t, string(byte('*')), string(c_src_code_multi_line_comment[69]))
	require.Equal(t, string(byte('/')), string(c_src_code_multi_line_comment[70]))
}

// should handle errors when we cannot find the closing bytes of a multi line comment
// and return a friendly message to the user
func TestAnalyzeTokenMultiLineCommentErrorC(t *testing.T) {
	src := []byte("int x = 0; /* @TEST_TODO no closing comment bytes")
	lex, err := lexer.NewLexer(src, "main.c")
	require.NoError(t, err)
	tokens, err := lex.AnalyzeTokens()
	require.Error(t, err)
	require.Empty(t, tokens)
	require.ErrorContains(
		t,
		err,
		"[main.c line 1]: Error: could not locate closing multi line comment: /* @TEST_TODO no closing comment byte",
	)
}

func TestParseCommentTokensSingleLineC(t *testing.T) {
	lex, err := lexer.NewLexer([]byte(c_src_code_single_line_comments), "main.c")
	require.NoError(t, err)
	tokens, err := lex.AnalyzeTokens()
	require.NoError(t, err)

	expectedComments := []lexer.Comment{
		{
			Title:          []byte("first single line comment"),
			Description:    []byte(nil),
			TokenIndex:     0,
			Source:         tokens[0].Lexeme,
			SourceFileName: "main.c",
		},
		{
			Title:          []byte("second single line comment"),
			Description:    []byte(nil),
			TokenIndex:     1,
			Source:         tokens[1].Lexeme,
			SourceFileName: "main.c",
		},
		{
			Title:       []byte("third single line comment"),
			Description: []byte(nil),
			// we skip index 2 because it's a string token
			TokenIndex:     3,
			Source:         tokens[3].Lexeme,
			SourceFileName: "main.c",
		},
	}

	actualComments, err := lex.Manager.ParseCommentTokens(lex, annotation)
	require.NoError(t, err)
	require.Equal(t, expectedComments, actualComments)
}

func TestParseCommentTokensMultiLineC(t *testing.T) {
	lex, err := lexer.NewLexer([]byte(c_src_code_multi_line_comment), "main.c")
	require.NoError(t, err)
	tokens, err := lex.AnalyzeTokens()
	require.NoError(t, err)

	expectedComments := []lexer.Comment{
		{
			Title:          []byte("inline 1"),
			Description:    []byte(nil),
			TokenIndex:     0,
			Source:         tokens[0].Lexeme,
			SourceFileName: "main.c",
		},
		{
			Title:          []byte("inline 2"),
			Description:    []byte(nil),
			TokenIndex:     1,
			Source:         tokens[1].Lexeme,
			SourceFileName: "main.c",
		},
		{
			Title:          []byte("multi line comment"),
			Description:    []byte("second line third line end line"),
			TokenIndex:     2,
			Source:         tokens[2].Lexeme,
			SourceFileName: "main.c",
		},
	}

	actualComments, err := lex.Manager.ParseCommentTokens(lex, annotation)
	require.NoError(t, err)
	require.Equal(t, expectedComments, actualComments)
}

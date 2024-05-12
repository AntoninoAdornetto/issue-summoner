package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

const (
	python_src_code_single_line_comments = `
	def main():
		x = 0 # @TEST_TODO first single line comment
		y = 0 # @TEST_TODO second single line comment
		
		print(f"X: {x}\tY: {y}");
		# @TEST_TODO third single line comment
	`

	python_src_code_multi_line_comment = `
	""""
	 @TEST_TODO multi line comment
	 second line
	 third line
	 end line
	""""
	def main():
		x = 0
		y = 0
		return x + y
	}
	`
)

// """
// @TEST_TODO multi line comment 2
// second line
// third line
// end line
// """

// should locate all single line comment tokens and store the comment contents
// as the lexeme for each item in the token slice output
func TestAnalyzeTokenSingleLineCommentPython(t *testing.T) {
	lex, err := lexer.NewLexer([]byte(python_src_code_single_line_comments), "main.py")
	require.NoError(t, err)
	actualTokens, err := lex.AnalyzeTokens()
	require.NoError(t, err)

	expectedTokens := []lexer.Token{
		{
			TokenType:      lexer.SINGLE_LINE_COMMENT,
			Lexeme:         []byte("# @TEST_TODO first single line comment"),
			Line:           3,
			StartByteIndex: 22,
			EndByteIndex:   59,
		},
		{
			TokenType:      lexer.SINGLE_LINE_COMMENT,
			Lexeme:         []byte("# @TEST_TODO second single line comment"),
			Line:           4,
			StartByteIndex: 69,
			EndByteIndex:   107,
		},
		{
			TokenType:      lexer.SINGLE_LINE_COMMENT,
			Lexeme:         []byte("# @TEST_TODO third single line comment"),
			Line:           7,
			StartByteIndex: 142,
			EndByteIndex:   179,
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

	// 22 = expectedTokens[0].StartByteIndex
	// 59 = expectedTokens[0].EndByteIndex
	require.Equal(t, byte('#'), python_src_code_single_line_comments[22])
	require.Equal(t, byte('t'), python_src_code_single_line_comments[59])

	// 69 = expectedTokens[1].StartByteIndex
	// 107 = expectedTokens[1].EndByteIndex
	require.Equal(t, byte('#'), python_src_code_single_line_comments[69])
	require.Equal(t, byte('t'), python_src_code_single_line_comments[107])

	// @TODO uncomment python string assertion once python lexer supports string tokens
	// // 153 = expectedTokens[2].StartByteIndex
	// // 166 = expectedTokens[2].EndByteIndex
	// // string token
	// require.Equal(t, byte('"'), python_src_code_single_line_comments[154])
	// require.Equal(t, byte('"'), python_src_code_single_line_comments[167])
}


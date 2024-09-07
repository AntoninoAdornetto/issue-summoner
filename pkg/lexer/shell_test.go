package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeTokenScan(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		expected   []lexer.Token
		flags      lexer.U8
		annotation []byte
	}{
		{
			name:       "should return the correct tokens from the shell script",
			path:       "./testdata/other/script.sh",
			flags:      lexer.FLAG_SCAN,
			annotation: testAnnotation,
			expected: []lexer.Token{
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_START,
					Lexeme: []byte{lexer.HASH},
					Line:   11,
					Start:  158,
					End:    158,
				},
				{
					Type:   lexer.TOKEN_COMMENT_ANNOTATION,
					Lexeme: testAnnotation,
					Line:   11,
					Start:  160,
					End:    175,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("iterate"),
					Line:   11,
					Start:  177,
					End:    183,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("through"),
					Line:   11,
					Start:  185,
					End:    191,
				},
				{
					Type:   lexer.TOKEN_COMMENT_TITLE,
					Lexeme: []byte("10"),
					Line:   11,
					Start:  193,
					End:    194,
				},
				{
					Type:   lexer.TOKEN_SINGLE_LINE_COMMENT_END,
					Lexeme: []byte{'\n'},
					Line:   12,
					Start:  195,
					End:    195,
				},
				{
					Type:   lexer.TOKEN_EOF,
					Lexeme: []byte{0},
					Line:   23,
					Start:  303,
					End:    303,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srcCode := getSrcCode(t, tc.path)
			base := lexer.NewLexer(tc.annotation, srcCode, tc.path, tc.flags)
			target, err := lexer.NewTargetLexer(base)
			require.NoError(t, err)
			tokens, err := base.AnalyzeTokens(target)
			require.NoError(t, err)
			require.Equal(t, tc.expected, tokens)
			for _, token := range tokens {
				if token.Type != lexer.TOKEN_EOF {
					require.Equal(t, token.Lexeme, srcCode[token.Start:token.End+1])
				}
			}
		})
	}
}

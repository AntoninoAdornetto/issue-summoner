package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

func TestBuildComments(t *testing.T) {
	testCases := []struct {
		name     string
		expected []lexer.Comment
		path     string
	}{
		{
			name: "should return the correct comments when parsing c code",
			path: "../../testdata/fixtures/c/mix.c",
			expected: []lexer.Comment{
				{
					Value:                "inline comment #1",
					TokenStartIndex:      0,
					TokenAnnotationIndex: 1,
					TokenEndIndex:        5,
				},
				{
					Value:                "inline comment #2",
					TokenStartIndex:      6,
					TokenAnnotationIndex: 7,
					TokenEndIndex:        11,
				},
				{
					Value:                "decode the message and clean up after yourself!",
					TokenStartIndex:      12,
					TokenAnnotationIndex: 13,
					TokenEndIndex:        22,
				},
				{
					Value:                "drop a star if you know about this code wars challenge Digital Cypher assigns to each letter of the alphabet unique number. Instead of letters in encrypted word we write the corresponding number Then we add to each obtained digit consecutive digits from the key",
					TokenStartIndex:      23,
					TokenAnnotationIndex: 24,
					TokenEndIndex:        70,
				},
			},
		},
	}

	for _, tc := range testCases {
		src := getSrcCode(t, tc.path)
		base := lexer.NewLexer(testAnnotation, src, tc.path)
		target, err := lexer.NewTargetLexer(base)
		require.NoError(t, err)

		tokens, err := base.AnalyzeTokens(target)
		require.NoError(t, err)

		actual := lexer.BuildComments(tokens)
		require.Equal(t, tc.expected, actual.Comments)
	}
}

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
			path: "./testdata/c/mix.c",
			expected: []lexer.Comment{
				{
					Title:                "inline comment #1",
					Description:          "", // single line comments don't have descriptions
					TokenStartIndex:      0,
					TokenAnnotationIndex: 1,
					TokenEndIndex:        5,
					LineNumber:           5,
				},
				{
					Title:                "inline comment #2",
					Description:          "", // single line comments don't have descriptions
					TokenStartIndex:      6,
					TokenAnnotationIndex: 7,
					TokenEndIndex:        11,
					LineNumber:           6,
				},
				{
					Title:                "decode the message and clean up after yourself!",
					Description:          "",
					TokenStartIndex:      12,
					TokenAnnotationIndex: 13,
					TokenEndIndex:        22,
					LineNumber:           10,
				},
				{
					// multi line comments have a description
					Title:                "drop a star if you know about this code wars challenge",
					Description:          "Digital Cypher assigns to each letter of the alphabet unique number. Instead of letters in encrypted word we write the corresponding number Then we add to each obtained digit consecutive digits from the key",
					TokenStartIndex:      23,
					TokenAnnotationIndex: 24,
					TokenEndIndex:        70,
					LineNumber:           14,
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

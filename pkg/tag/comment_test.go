package tag_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/stretchr/testify/require"
)

func TestGetCommentSyntax_CLanguages(t *testing.T) {
	programmingLanguages := []string{
		".c",
		".cpp",
		".java",
		".js",
		".jsx",
		".ts",
		".tsx",
		".cs",
		".go",
		".php",
		".swift",
		".kt",
		".rs",
		".m",
		".scala",
	}

	for _, lang := range programmingLanguages {
		commentSyntax := tag.GetCommentSyntax(lang)
		require.Equal(t, tag.CommentSyntaxMap["c-derived"], commentSyntax)
	}
}

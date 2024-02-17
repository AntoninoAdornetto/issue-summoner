package tag_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/stretchr/testify/require"
)

func TestGetCommentSyntax_CLanguages(t *testing.T) {
	fileExtensions := []string{
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

	for _, ext := range fileExtensions {
		commentSyntax := tag.CommentSyntax(ext)
		require.Equal(t, tag.CommentSyntaxMap["c-derived"], commentSyntax)
	}
}

func TestGetCommentSyntax_Default(t *testing.T) {
	unrecognizedExtensions := []string{
		".gitignore",
		"LICENSE",
		"Makefile",
	}

	for _, ext := range unrecognizedExtensions {
		commentSyntax := tag.CommentSyntax(ext)
		require.Equal(t, tag.CommentSyntaxMap["default"], commentSyntax)
	}
}

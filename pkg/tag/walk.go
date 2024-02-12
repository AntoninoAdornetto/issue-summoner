package tag

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type WalkTagManager interface {
	FindTags(path string, fileOperator TagFileOperator) ([]Tag, error)
	CompileSingleLineComment(fileInfo fs.FileInfo) regexp.Regexp
}

type WalkFileOperator interface {
	Open(fileName string) (*os.File, error)
	WalkDir(root string, fn fs.WalkDirFunc) error
}

type WalkParams struct {
	Root           string
	TagManager     WalkTagManager
	FileOperator   WalkFileOperator
	IgnorePatterns []GitIgnorePattern
}

func Walk(arg WalkParams) ([]Tag, error) {
	tags := make([]Tag, 0)

	err := arg.FileOperator.WalkDir(arg.Root, func(path string, d fs.DirEntry, err error) error {
		isValidPath := validatePath(path, arg.IgnorePatterns)

		if d.IsDir() {
			isGitDir := strings.Contains(d.Name(), ".git")

			if isGitDir || !isValidPath {
				return filepath.SkipDir
			}
			return nil
		}

		if !isValidPath {
			return nil
		}

		foundTags, err := arg.TagManager.FindTags(path, arg.FileOperator)
		if err != nil {
			return err
		}

		tags = append(tags, foundTags...)

		return err
	})

	return tags, err
}

func validatePath(path string, ignorePatterns []GitIgnorePattern) bool {
	for _, v := range ignorePatterns {
		matched := v.Match([]byte(path))
		if matched {
			return false
		}
	}
	return true
}

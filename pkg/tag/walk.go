package tag

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type WalkTagManager interface {
	FindTags(path string) ([]Tag, error)
}

type WalkFileOperator interface {
	Open(fileName string) (*os.File, error)
	WalkDir(root string, fn fs.WalkDirFunc) error
}

func Walk(root string, tm WalkTagManager, wo WalkFileOperator) ([]Tag, error) {
	tags := make([]Tag, 0)
	ignorePatterns, err := ProcessIgnorePatterns(filepath.Join(root, GitIgnoreFile), wo)
	if err != nil {
		log.Fatal(err)
	}

	err = wo.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		isGitDir := strings.Contains(d.Name(), ".git")
		isValidPath := validatePath(fmt.Sprintf("%s/", path), &ignorePatterns)

		if d.IsDir() {
			if isGitDir || !isValidPath {
				return filepath.SkipDir
			}
			return nil
		}

		if !isValidPath {
			return nil
		}

		foundTags, err := tm.FindTags(path)
		if err != nil {
			return err
		}

		tags = append(tags, foundTags...)

		return err
	})
	return tags, nil
}

func validatePath(path string, ignorePatterns *GitIgnorePattern) bool {
	for _, v := range *ignorePatterns {
		matched := v.Match([]byte(path))
		if matched {
			return false
		}
	}
	return true
}

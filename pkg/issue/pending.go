package issue

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type PendingIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *PendingIssue) Walk(root string, gitIgnore []regexp.Regexp) (int, error) {
	n := 0
	foundGitDir := false
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !foundGitDir {
			if skip := skipGitDir(d.Name(), d.IsDir()); skip {
				return filepath.SkipDir
			}
		}

		if skip := skipIgnoreMatch(path, gitIgnore); skip {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		n++
		return pi.Scan(file)
	})

	return n, err
}

func (pi *PendingIssue) Scan(r io.Reader) error {
	return nil
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

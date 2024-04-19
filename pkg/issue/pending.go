package issue

import (
	"bufio"
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

		if d.IsDir() {
			if !foundGitDir && skipGitDir(d.Name()) {
				foundGitDir = true
				return filepath.SkipDir
			}
			return nil
		}

		if skip := skipIgnoreMatch(path, gitIgnore); skip {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		n++
		return pi.Scan(file, path)
	})

	return n, err
}

func (pi *PendingIssue) Scan(r io.Reader, path string) error {
	lineNum := uint64(0)
	issues := make([]Issue, 0)
	scanner := bufio.NewScanner(r)
	notation := NewCommentNotation(filepath.Ext(path), pi.Annotation, scanner)

	for scanner.Scan() {
		lineNum++
		issue, err := notation.ParseLine(&lineNum)
		if err != nil {
			return err
		}

		if issue.AnnotationLineNumber > 0 {
			issue.ID = generateID(path, issue.AnnotationLineNumber)
			issue.FilePath = path
			issues = append(issues, issue)
		}
	}

	pi.Issues = append(pi.Issues, issues...)
	return scanner.Err()
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

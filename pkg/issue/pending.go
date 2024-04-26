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
	lineNum := 0
	issues := make([]Issue, 0)
	scanner := bufio.NewScanner(r)
	notation := NewCommentNotation(filepath.Ext(path))

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		prefixLoc, commentType := notation.FindPrefixLocations(line)
		if commentType == "" {
			continue
		}

		arg := NewCommentManagerParams{
			Annotation:      pi.Annotation,
			CommentType:     commentType,
			FilePath:        path,
			FileName:        filepath.Base(path),
			StartLineNumber: lineNum,
			Locations:       prefixLoc,
			Scanner:         scanner,
		}

		if commentType == SINGLE_LINE_COMMENT {
			arg.PrefixRe = notation.SingleLinePrefixRe
			arg.SuffixRe = notation.SingleLineSuffixRe
		} else {
			// @TODO create multi.go & implement multi-line comment parsing
			continue
		}

		cm, err := NewCommentManager(arg)
		if err != nil {
			return err
		}

		comments, err := cm.ParseComment(lineNum)
		if err != nil {
			return err
		}

		issues = append(issues, comments...)
	}

	pi.Issues = append(pi.Issues, issues...)
	return scanner.Err()
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

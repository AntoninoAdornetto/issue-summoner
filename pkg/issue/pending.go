package issue

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
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

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		n++
		return pi.Scan(src, path)
	})

	return n, err
}

func (pi *PendingIssue) Scan(src []byte, path string) error {
	base := filepath.Base(path)
	ext := filepath.Ext(base)

	// @TODO remove file ext check when additional language support is added. This is for testing phase
	if !lexer.IsAdoptedFromC(ext) {
		return nil
	}

	lex, err := lexer.NewLexer(src, base)
	if err != nil {
		return err
	}

	tokens, err := lex.AnalyzeTokens()
	if err != nil {
		return err
	}

	fmt.Printf("Found %d tokens in %s\n", len(tokens), base)
	return nil
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

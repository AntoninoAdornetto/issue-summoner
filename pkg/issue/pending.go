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

	/*
	* @TODO remove file extension check when additional language support is added.
	* The implementation has changed drastically and is now utilizing scanning/lexing
	* approach. I am starting out with languages that have adopted similar comment syntax
	* to C since it's the most common. Once more support is added, we can remove this check
	 */

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

	comments, err := lex.Manager.ParseCommentTokens(lex, []byte(pi.Annotation))
	if err != nil {
		return err
	}

	for _, c := range comments {
		token := tokens[c.TokenIndex]
		pi.Issues = append(pi.Issues, Issue{
			ID:          fmt.Sprintf("%s-%d:%d", base, token.StartByteIndex, token.EndByteIndex),
			Title:       string(c.Title),
			Description: string(c.Description),
			FileName:    base,
			FilePath:    path,
			LineNumber:  token.Line,
		})
	}

	return nil
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

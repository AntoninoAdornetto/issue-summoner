package issue

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
)

type PendingIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *PendingIssue) Walk(root string, gitIgnore []regexp.Regexp) (int, error) {
	n := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// @TODO Flag for Walking/Scanning hidden dirs? Revisit this thought
			if strings.HasPrefix(d.Name(), ".") {
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
	lex, err := lexer.NewLexer(src, base)
	if err != nil {
		return nil
	}

	tokens, err := lex.AnalyzeTokens()
	if err != nil {
		return err
	}
	for _, token := range tokens {
		fmt.Println(string(token.Lexeme))
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

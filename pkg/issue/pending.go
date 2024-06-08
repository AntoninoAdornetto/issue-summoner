package issue

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/AntoninoAdornetto/go-gitignore"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
)

type PendingIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *PendingIssue) Walk(root string) (int, error) {
	n := 0
	ignorer, err := ignore.NewIgnorer(root)
	if err != nil {
		return n, err
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// @TODO Flag for Walking/Scanning hidden dirs? Revisit this thought
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			isIgnored, err := ignorer.Match(path)
			if err != nil {
				return err
			}

			if isIgnored {
				return filepath.SkipDir
			}

			return nil
		}

		isIgnored, err := ignorer.Match(path)
		if err != nil {
			return err
		}

		if isIgnored {
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
	* The implementation has changed drastically and is now utilizing a scanning/lexing
	* approach. I am starting out with languages that have adopted similar comment syntax
	* to C, since it's the most common. Once more languages are added, we can remove this check
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
			StartIndex:  token.StartByteIndex,
			EndIndex:    token.EndByteIndex,
		})
	}

	return nil
}

func (pi *PendingIssue) WriteIssueID(id int64, issueIndex int) error {
	if len(pi.Issues) == 0 {
		return errors.New("cannot write issue_id with an empty issue slice")
	}

	if issueIndex < 0 || issueIndex > len(pi.Issues) {
		return fmt.Errorf(
			"issue index %d out of range. issue slice len: %d",
			issueIndex,
			len(pi.Issues),
		)
	}

	currentIssue := pi.Issues[issueIndex]
	file, err := os.OpenFile(currentIssue.FilePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	src, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	start, end := currentIssue.StartIndex, currentIssue.EndIndex
	comment := src[start : end+1]
	annotationId := fmt.Sprintf("(%d)", id)
	newAnnotation := "@ISSUE" + annotationId
	comment = bytes.Replace(comment, []byte(pi.Annotation), []byte(newAnnotation), 1)
	buf := make([]byte, 0)

	for i := 0; i < len(src); i++ {
		if i == start {
			for j, rn := range comment {
				buf = append(buf, rn)
				if j == len(newAnnotation)-1 {
					i = end
				}
			}
		} else {
			buf = append(buf, src[i])
		}
	}

	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	if err = file.Truncate(0); err != nil {
		return err
	}

	if _, err = file.Write(buf); err != nil {
		return err
	}

	return nil
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

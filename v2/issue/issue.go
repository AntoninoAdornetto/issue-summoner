package issue

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/AntoninoAdornetto/go-gitignore"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/charmbracelet/lipgloss"
)

type IssueManager struct {
	annotation  string
	currentPath string
	currentBase string
	Issues      []Issue
}

type Issue struct {
	ID          string
	Title       string
	Description string
	FilePath    string
	FileName    string
	LineNumber  int
	ColStart    int
	ColEnd      int
}

func NewIssueManager(annotation string) *IssueManager {
	return &IssueManager{annotation: annotation, Issues: make([]Issue, 0, 10)}
}

func (iMan *IssueManager) NewIssue(com lexer.Comment, token lexer.Token) Issue {
	id := fmt.Sprintf("%s-%d:%d", iMan.currentBase, token.StartByteIndex, token.EndByteIndex)
	return Issue{
		ID:          id,
		Title:       string(com.Title),
		Description: string(com.Description),
		FileName:    iMan.currentBase,
		FilePath:    iMan.currentPath,
		LineNumber:  token.Line,
		ColStart:    token.StartByteIndex,
		ColEnd:      token.EndByteIndex,
	}
}

func (iMan *IssueManager) Walk(root string) error {
	ignorer, err := ignore.NewIgnorer(root)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		if ignorer == nil {
			return nil
		}

		matched, err := ignorer.Match(path)
		if err != nil {
			return err
		}

		if matched {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		iMan.currentBase = d.Name()
		iMan.currentPath = path
		return iMan.scan(path)
	})
}

func (iMan *IssueManager) scan(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}

	defer src.Close()
	buf, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	lex, err := lexer.NewLexer(buf, iMan.currentBase)
	if err != nil {
		// @TODO return actual error once more languages are supported. For now just skip
		return nil
	}

	tokens, err := lex.AnalyzeTokens()
	if err != nil {
		return err
	}

	comments, err := lex.Manager.ParseCommentTokens(lex, []byte(iMan.annotation))
	if err != nil {
		return err
	}

	for _, com := range comments {
		token := tokens[com.TokenIndex]
		iMan.Issues = append(iMan.Issues, iMan.NewIssue(com, token))
	}

	return nil
}

func (iMan *IssueManager) Print(propStyle, valStyle lipgloss.Style) {
	for _, issue := range iMan.Issues {
		fmt.Printf("\n\n")
		paths := strings.Split(issue.FilePath, "/")
		fmt.Println(
			propStyle.Render("Filename: "),
			valStyle.Render(paths[len(paths)-1]),
		)
		fmt.Println(propStyle.Render("Title: "), valStyle.Render(issue.Title))
		fmt.Println(
			propStyle.Render("Description: "),
			valStyle.Render(issue.Description),
		)
		fmt.Println(
			propStyle.Render("Line number: "),
			valStyle.Render(fmt.Sprintf("%d", issue.LineNumber)),
		)
	}
}

package issue

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	ignore "github.com/AntoninoAdornetto/go-gitignore"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/charmbracelet/lipgloss"
)

type IssueManager struct {
	annotation  string
	currentPath string
	currentBase string
	template    *template.Template
	os          string
	Issues      []Issue
}

type Issue struct {
	ID          string
	Title       string
	Body        string
	Description string
	Environment string // @TODO change Environment to os -> operating system
	FilePath    string
	FileName    string
	LineNumber  int
	ColStart    int // not really the column index, a more accurate description is the start byte index
	ColEnd      int // not really the column index, a more accurate description is the end byte index
}

func NewIssueManager(annotation string, isReporting bool) (*IssueManager, error) {
	iMan := &IssueManager{annotation: annotation, Issues: make([]Issue, 0, 10)}
	if !isReporting {
		return iMan, nil
	}

	template, err := generateIssueTemplate()
	if err != nil {
		return nil, err
	}

	iMan.os = runtime.GOOS
	iMan.template = template
	return iMan, nil
}

func (iMan *IssueManager) NewIssue(com lexer.Comment, token lexer.Token) (Issue, error) {
	id := fmt.Sprintf("%s-%d:%d", iMan.currentBase, token.StartByteIndex, token.EndByteIndex)

	issue := Issue{
		ID:          id,
		Title:       string(com.Title),
		Description: string(com.Description),
		Environment: iMan.os,
		FileName:    iMan.currentBase,
		FilePath:    iMan.currentPath,
		LineNumber:  token.Line,
		ColStart:    token.StartByteIndex,
		ColEnd:      token.EndByteIndex,
	}

	// will not be nil if running the scan command. See params to NewIssueManager func
	if iMan.template == nil {
		return issue, nil
	}

	buf := bytes.Buffer{}
	if err := iMan.template.Execute(&buf, issue); err != nil {
		return issue, err
	}

	issue.Body = buf.String()
	return issue, nil
}

func (iMan *IssueManager) Walk(root string) error {
	ignorer, err := ignore.NewIgnorer(root)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

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
		issue, err := iMan.NewIssue(com, token)
		if err != nil {
			return err
		}
		iMan.Issues = append(iMan.Issues, issue)
	}

	return nil
}

// Print is invoked used when the verbose flag is passed to the scan command
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

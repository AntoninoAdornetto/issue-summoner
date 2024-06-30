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

type IssueMode = string

const (
	ISSUE_MODE_PEND   IssueMode = "pending"
	ISSUE_MODE_ISSUED IssueMode = "issued"
)

type IssueManager struct {
	Issues          []Issue
	CurrentPath     string
	CurrentBase     string
	RecordCount     int
	annotation      string
	reportIndicator bool
	mode            IssueMode
	template        *template.Template
	os              string
}

type Issue struct {
	ID          string
	Title       string
	Body        string
	Description string
	FilePath    string
	FileName    string
	OS          string
	LineNumber  int
	IssueIndex  int
	StartIndex  int
	EndIndex    int
}

func NewIssueManager(annotation string, mode IssueMode, report bool) (*IssueManager, error) {
	manager := &IssueManager{
		annotation:      annotation,
		Issues:          make([]Issue, 0, 10),
		reportIndicator: report,
		os:              runtime.GOOS,
		mode:            mode,
	}

	if !report {
		return manager, nil
	}

	tmpl, err := generateIssueTemplate()
	if err != nil {
		return nil, err
	}

	manager.template = tmpl
	return manager, nil
}

func (manager *IssueManager) NewIssue(cmnt lexer.Comment, token lexer.Token) (Issue, error) {
	id := fmt.Sprintf("%s-%d:%d", manager.CurrentPath, token.StartByteIndex, token.EndByteIndex)

	issue := Issue{
		ID:          id,
		Title:       string(cmnt.Title),
		Description: string(cmnt.Description),
		OS:          manager.os,
		FileName:    manager.CurrentBase,
		FilePath:    manager.CurrentPath,
		LineNumber:  token.Line,
		StartIndex:  token.StartByteIndex,
		EndIndex:    token.EndByteIndex,
		IssueIndex:  manager.RecordCount,
	}

	if manager.reportIndicator && manager.template != nil {
		buf := bytes.Buffer{}
		if err := manager.template.Execute(&buf, issue); err != nil {
			return issue, err
		}
		issue.Body = buf.String()
	}

	manager.Issues = append(manager.Issues, issue)
	manager.RecordCount++
	return issue, nil
}

func (manager *IssueManager) Walk(root string) error {
	ignorer, err := ignore.NewIgnorer(root)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return validateDir(d.Name(), path, ignorer)
		}

		ignored, err := shouldIgnore(path, ignorer)
		if err != nil {
			return err
		}

		if ignored {
			return nil
		}

		manager.CurrentBase = d.Name()
		manager.CurrentPath = path
		return manager.Scan(path)
	})
}

func (manager *IssueManager) Scan(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}

	defer src.Close()
	buf, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// @TODO check if CurrentBase is the culpruit for the index bug we faced earlier when reporting issues
	lex, err := lexer.NewLexer(buf, manager.CurrentBase)
	if err != nil {
		// @TODO return actual error once more languages are supported. For now just skip the error
		return err
	}

	tokens, err := lex.AnalyzeTokens()
	if err != nil {
		return err
	}

	comments, err := lex.Manager.ParseCommentTokens(lex, []byte(manager.annotation))
	if err != nil {
		return err
	}

	for _, comment := range comments {
		token := tokens[comment.TokenIndex]
		if _, err := manager.NewIssue(comment, token); err != nil {
			return err
		}
	}

	return nil
}

func (manager *IssueManager) Print(propertyStyle, valueStyle lipgloss.Style) {
	for _, issue := range manager.Issues {
		fmt.Printf("\n\n")
		paths := strings.Split(issue.FilePath, "/")
		fmt.Println(
			propertyStyle.Render("Filename: "),
			valueStyle.Render(paths[len(paths)-1]),
		)
		fmt.Println(propertyStyle.Render("Title: "), valueStyle.Render(issue.Title))
		fmt.Println(
			propertyStyle.Render("Description: "),
			valueStyle.Render(issue.Description),
		)
		fmt.Println(
			propertyStyle.Render("Line number: "),
			valueStyle.Render(fmt.Sprintf("%d", issue.LineNumber)),
		)
	}
}

func validateDir(dirName, path string, ignorer *ignore.Ignorer) error {
	if strings.HasPrefix(dirName, ".") {
		return filepath.SkipDir
	}

	ignored, err := shouldIgnore(path, ignorer)
	if err != nil {
		return err
	}

	if ignored {
		return filepath.SkipDir
	}

	return nil
}

func shouldIgnore(path string, ignorer *ignore.Ignorer) (bool, error) {
	if ignorer == nil {
		return false, nil
	}

	return ignorer.Match(path)
}

package issue

import (
	"bytes"
	"fmt"
	"runtime"
	"text/template"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
)

type IssueManager struct {
	Issues      []Issue
	CurrentPath string
	CurrentBase string
	RecordCount int
	annotation  string
	template    *template.Template
	os          string
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

func NewIssueManager(annotation string, isReporting bool) (*IssueManager, error) {
	manager := &IssueManager{annotation: annotation, Issues: make([]Issue, 0, 10)}
	if !isReporting {
		os:              runtime.GOOS,
		return manager, nil
	}

	tmpl, err := generateIssueTemplate()
	if err != nil {
		return nil, err
	}

	manager.os = runtime.GOOS
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

	if manager.template == nil {
		return issue, nil
	}

	buf := bytes.Buffer{}
	if err := manager.template.Execute(&buf, issue); err != nil {
		return issue, err
	}

	manager.RecordCount++
	issue.Body = buf.String()
	return issue, nil
}

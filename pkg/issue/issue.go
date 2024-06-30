package issue

import (
	"runtime"
	"text/template"
)

type IssueManager struct {
	Issues      []Issue
	annotation  string
	currentPath string
	currentBase string
	template    *template.Template
	os          string
	recordCount int
}

type Issue struct {
	ID          string
	Title       string
	Body        string
	Description string
	FilePath    string
	FileName    string
	LineNumber  int
	IssueIndex  int
	StartIndex  int
	EndIndex    int
}

func NewIssueManager(annotation string, isReporting bool) (*IssueManager, error) {
	manager := &IssueManager{annotation: annotation, Issues: make([]Issue, 0, 10)}
	if !isReporting {
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

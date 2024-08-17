package issue

import (
	"bytes"
	"errors"
	"runtime"
	"text/template"
)

type IssueMode = string

const (
	IssueModePurge IssueMode = "purge"
	IssueModeScan  IssueMode = "scan"
)

type Issue struct {
	ID          string
	Title       string
	Description string
	FilePath    string
	FileName    string
	LineNumber  int
	Environment string
	StartIndex  int
	EndIndex    int
}

type IssueManager interface {
	GetIssues() []Issue
	Scan(src []byte, path string) error
	Walk(root string) (int, error)
	WriteIssueID(id int64, issueIndex int) error
}

// NewIssueManager will return either a PendingIssue struct or ProcessedIssue struct
// that satisfies the methods defined in the IssueManager interface. The methods in
// said interface are used to report new issues an SCM or locate issues that have been
// reported to an SCM. Each struct will implement methods for walking the project directory
// and parsing source code files. The main difference is that pending issues will have an
// annotation with no id, since they haven't been pushed to an scm yet, and processed issues
// will have their original annotation plus an id so they can be located and removed from the
// source code file at a later time.
func NewIssueManager(issueType string, annotation string) (IssueManager, error) {
	switch issueType {
	case PENDING_ISSUE:
		return &PendingIssue{Annotation: annotation}, nil
	case PROCESSED_ISSUE:
		return &ProcessedIssue{Annotation: annotation}, nil
	switch mode {
	case IssueModeScan:
		manager.annotation = annotation
	case IssueModePurge:
		annotation = append(annotation, []byte("\\(\\d+\\)")...)
		manager.annotation = annotation
	default:
		return nil, errors.New("expected mode of \"report\" or \"purge\"")
	}
}

func (issue *Issue) ExecuteIssueTemplate(tmpl *template.Template) ([]byte, error) {
	buf := bytes.Buffer{}
	issue.Environment = runtime.GOOS
	err := tmpl.Execute(&buf, issue)
	return buf.Bytes(), err
}

package issue

import (
	"bytes"
	"errors"
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
	default:
		return nil, errors.New("Unsupported issue type. Use 'pending' or 'processed'")
	}
}

func (issue *Issue) ExecuteIssueTemplate(tmpl *template.Template) ([]byte, error) {
	buf := bytes.Buffer{}
	issue.Environment = runtime.GOOS
	err := tmpl.Execute(&buf, issue)
	return buf.Bytes(), err
}

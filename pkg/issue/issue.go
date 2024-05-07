package issue

import (
	"bytes"
	"errors"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

const (
	PENDING_ISSUE   = "pending"
	PROCESSED_ISSUE = "processed"
)

type Issue struct {
	ID          string
	Title       string
	Description string
	FilePath    string
	FileName    string
	LineNumber  int
	Environment string
}

type IssueManager interface {
	GetIssues() []Issue
	Scan(src []byte, path string) error
	Walk(root string, ignore []regexp.Regexp) (int, error)
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

func skipGitDir(name string) bool {
	return strings.Contains(name, ".git")
}

func skipIgnoreMatch(path string, patterns []regexp.Regexp) bool {
	for _, re := range patterns {
		if matched := re.MatchString(path); matched {
			return true
		}
	}
	return false
}

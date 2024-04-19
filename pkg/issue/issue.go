/*
The issue package is responsible for handling pending and processed issues.
Issues are objects that describe a task, concern, or area of code that
requires attention.

Issues are discovered by parsing each line in a source file. When a file is
first opened, we determine the syntax used to denote single and multi line comments
by checking the extenion of the file. I.E main.c, main.go, main.cpp and so on.
This will allow us to locate single and multi line comment symbols in each source file.

As we scan each line, we check if the line contains the prefix notation for a single
or multi line comment. If it contains a single line comment prefix, we parse and get
all text after the prefix and annotation, if it exists. If an issue annotation was discovered,
it qualifies for an issue and we will append a new issue object onto an issues slice.

In the case of a multi line comment prefix, we will continue to scan subsequent lines until
reaching the multi line comment suffix notation. At that point we will check if an issue annotation
was discovered. If so, we append that new issue item onto the issues slice.

There are two types of issues, pending and processed. Pending issues are those that have not yet been
uploaded to a source code management platform. Processed issues are issues that have been uploaded to a
source code management platform and will have a unique id number associated with the issue annotation.

See the issue_test.go file for examples.
*/
package issue

import (
	"errors"
	"io"
	"os"
	"regexp"
)

const (
	PENDING_ISSUE   = "pending"
	PROCESSED_ISSUE = "processed"
)

type Issue struct {
	ID                   string
	Title                string
	Description          string
	FileInfo             os.FileInfo
	StartLineNumber      uint64
	EndLineNumber        uint64
	AnnotationLineNumber uint64
}

type IssueManager interface {
	GetIssues() []Issue
	Scan(r io.Reader) error
	Walk(root string, ignore []regexp.Regexp) error
}

// NewIssueManager will return either a PendingIssue struct or ProcessedIssue struct
// that satisfies the methods defined in the IssueManager interface. The methods in
// said interface are used to report new issues an SCM or locate issues that have been
// reported to an SCM. Each struct will implement methods for walking the project directory
// and parsing source code files. The main difference is that pending issues can be uploaded
// to an SCM and processed issues can be resolved and the matching comment in the source code
// can be removed through it's methods.
func NewIssueManager(issueType string, annotation string) (IssueManager, error) {
	switch issueType {
	default:
		return nil, errors.New("Unsupported issue type. Use 'pending' or 'processed'")
	}
}

/*
The issue package is responsible for handling pending and processed issues.
Issues are objects that describe a task, concern, or area of code that require
some attention.

Issues are discovered by parsing single and multi line comments in source code files.
In order for a comment to be considered as an issue, the comment must have an annotation.
The annotation can be as simple as // @TODO. Or as complex as // @TICKET_123_REVIEW
Once an annotation is found, we parse the surrounding data and build an issue object
that will be used later on.

There are two types of issues, pending and processed. Pending issues have not yet been
uploaded to a source code management platform. Processed issues have been uploaded to
a source code management platform and will have an ID number right beside the annotation.
A simple processed issue may look like this: // @TODO(1234) where 1234 is the ID that the
source code management platform returned after making an http request to create the issue.
*/
package issue

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const (
	PENDING_ISSUE   = "pending"
	PROCESSED_ISSUE = "processed"
)

type Issue struct {
	ID                   string
	Title                string
	Description          string // @TODO use a string builder to reduce copying
	FileInfo             os.FileInfo
	StartLineNumber      uint64
	EndLineNumber        uint64
	AnnotationLineNumber uint64
	IsSingleLine         bool
	IsMultiLine          bool
}

// IssueManager is responsible for defining the methods
// we will use for parsing single and multi line comments.
type IssueManager interface {
	GetIssues() []Issue
	Scan(file *os.File) error
}

type ParseCommentParams struct {
	LineText      string
	LineType      string
	LineNum       *uint64
	Scanner       *bufio.Scanner
	Comment       Comment
	CommentPrefix string
	FileInfo      os.FileInfo
}

// GetIssueManager takes an issue type as input and returns
// a new struct that satisfies the IssueManager interface.
// The PendingIssue struct is in charge of issues that have not
// been reported to a source code management platform yet.
// ProcessedIssue struct is in charge of issues that have
// already been reported. An error is returned if an unsupported
// issueType is passed into the function
func GetIssueManager(issueType string, annotation string) (IssueManager, error) {
	switch issueType {
	case PENDING_ISSUE:
		return &PendingIssue{Issues: make([]Issue, 0), Annotation: annotation}, nil
	case PROCESSED_ISSUE:
		return &ProcessedIssue{Issues: make([]Issue, 0)}, nil
	default:
		return nil, errors.New("Unsupported issue type. Please use pending or processed")
	}
}

func (is *Issue) Init(lineType string, lineNum uint64, fi *os.FileInfo) {
	is.StartLineNumber = lineNum
	is.EndLineNumber = lineNum
	is.IsSingleLine = lineType == LINE_TYPE_SINGLE
	is.IsMultiLine = lineType == LINE_TYPE_MULTI_START
	is.FileInfo = *fi
	is.ID = fmt.Sprintf("%s-%d", is.FileInfo.Name(), lineNum)
}

func (is *Issue) Build(content string, isAnnotated bool, scannedLines uint64) {
	if isAnnotated {
		is.AnnotationLineNumber = is.StartLineNumber + scannedLines
		is.Title = content
	} else if is.Description == "" {
		is.Description = content
		is.EndLineNumber = scannedLines + is.StartLineNumber
	} else {
		is.Description = fmt.Sprintf("%s %s", is.Description, content)
		is.EndLineNumber = scannedLines + is.StartLineNumber
	}
}

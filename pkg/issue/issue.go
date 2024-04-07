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
	"os"
	"strings"
)

const (
	PENDING_ISSUE         = "pending"
	PROCESSED_ISSUE       = "processed"
	LINE_TYPE_SRC_CODE    = "c"
	LINE_TYPE_SINGLE      = "single"
	LINE_TYPE_MULTI_START = "multi-start"
	LINE_TYPE_MULTI_END   = "multi-end"
)

type Issue struct {
	ID                   string
	Title                string
	Description          string
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
	ParseComment(ParseCommentParams) error
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

func ParseSingleLineComment(line string, annotation string, prefix string) (string, bool) {
	index := 1
	isAnnotated := false
	fields := strings.Fields(line)

	for i, s := range fields {
		if s == annotation {
			isAnnotated = true
			index = i + 1
		}
	}

	return strings.Join(fields[index:], " "), isAnnotated
}

// EvalSourceLine evaluates the line (read from a bufio scanner) input
// and determines if the line is a single line comment **(LINE_TYPE_SINGLE)**,
// a multi line comment **(LINE_TYPE_MULTI)** or sorce code **(LINE_TYPE_SRC_CODE)**
// In addition, this function is used within the Scan implementations of the
// PendingIssue & ProcessedIssue structs to determine if we should proceed
// with the parsing of a single/multi comment line or skip the process entirely.
func EvalSourceLine(line string, c Comment) (string, string) {
	for _, s := range c.SingleLinePrefix {
		if strings.HasPrefix(line, s) {
			return LINE_TYPE_SINGLE, s
		}
	}

	for i := range c.MultiLineStartPrefix {
		isMultiStart := strings.HasPrefix(line, c.MultiLineStartPrefix[i])
		isMultiEnd := strings.HasSuffix(line, c.MultiLineEndPrefix[i])

		if isMultiStart {
			return LINE_TYPE_MULTI_START, c.MultiLineStartPrefix[i]
		}

		if isMultiEnd {
			return LINE_TYPE_MULTI_END, c.MultiLineEndPrefix[i]
		}
	}

	return LINE_TYPE_SRC_CODE, ""
}

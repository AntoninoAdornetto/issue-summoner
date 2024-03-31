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
	"errors"
	"os"
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
	IsSingleLine         bool
	IsMultiLine          bool
}

type IssueManager interface {
	Scan() ([]Issue, error)
}

func GetIssueManager(issueType string) (IssueManager, error) {
	switch issueType {
	case PENDING_ISSUE:
		return &PendingIssue{}, nil
	case PROCESSED_ISSUE:
		return &ProcessedIssue{}, nil
	default:
		return nil, errors.New("Unsupported issue type. Please use pending or processed")
	}
}

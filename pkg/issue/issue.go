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

package issue

import (
	"os"
)

type ProcessedIssue struct {
	Issues []Issue
}

func (pi *ProcessedIssue) Scan(file *os.File) ([]Issue, error) {
	return pi.Issues, nil
}

func (pi *ProcessedIssue) ParseComment(arg ParseCommentParams) error {
	return nil
}

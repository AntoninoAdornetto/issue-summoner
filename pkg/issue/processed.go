package issue

import (
	"os"
)

type ProcessedIssue struct {
	Issues []Issue
}

func (pi *ProcessedIssue) GetIssues() []Issue {
	return pi.Issues
}

func (pi *ProcessedIssue) Scan(file *os.File) error {
	return nil
}

func (pi *ProcessedIssue) ParseComment(arg ParseCommentParams) error {
	return nil
}

package issue

import (
	"io"
	"regexp"
)

type ProcessedIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *ProcessedIssue) Walk(root string, ignore []regexp.Regexp) error {
	return nil
}

func (pi *ProcessedIssue) Scan(r io.Reader) error {
	return nil
}

func (pi *ProcessedIssue) GetIssues() []Issue {
	return pi.Issues
}

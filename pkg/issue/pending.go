package issue

import (
	"io"
	"regexp"
)

type PendingIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *PendingIssue) Walk(root string, ignore []regexp.Regexp) error {
	return nil
}

func (pi *PendingIssue) Scan(r io.Reader) error {
	return nil
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

package issue

import (
	"errors"
	"runtime"
)

type IssueMode = string

const (
	IssueModePurge IssueMode = "purge"
	IssueModeScan  IssueMode = "scan"
)

type IssueManager struct {
	Issues     []Issue
	annotation []byte
	mode       IssueMode
	os         string
}

type Issue struct{}

func NewIssueManager(annotation []byte, mode IssueMode) (*IssueManager, error) {
	manager := &IssueManager{
		Issues: make([]Issue, 0),
		mode:   mode,
		os:     runtime.GOOS,
	}

	switch mode {
	case IssueModeScan:
		manager.annotation = annotation
	case IssueModePurge:
		annotation = append(annotation, []byte("\\(\\d+\\)")...)
		manager.annotation = annotation
	default:
		return nil, errors.New("expected mode of \"report\" or \"purge\"")
	}

	return manager, nil
}

package issue

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	ignore "github.com/AntoninoAdornetto/go-gitignore"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
)

type IssueMode = string

const (
	IssueModePurge  IssueMode = "purge"
	IssueModeReport IssueMode = "report"
	IssueModeScan   IssueMode = "scan"
)

type IssueManager struct {
	Issues      []Issue
	annotation  []byte
	currentBase string
	currentPath string
	mode        IssueMode
	os          string
	template    *template.Template
}

type Issue struct {
	ID          string
	Body        string
	Description string
	EndIndex    int
	FileName    string
	FilePath    string
	LineNumber  int
	OS          string
	StartIndex  int
	Title       string
}

func NewIssueManager(annotation []byte, mode IssueMode) (*IssueManager, error) {
	manager := &IssueManager{
		Issues: make([]Issue, 0),
		mode:   mode,
		os:     runtime.GOOS,
	}

	switch mode {
	case IssueModeScan, IssueModeReport:
		manager.annotation = annotation
	case IssueModePurge:
		annotation = append(annotation, []byte("\\(\\d+\\)")...)
		manager.annotation = annotation
	default:
		return nil, errors.New("expected mode of \"report\" or \"purge\"")
	}

	return manager, nil
}

func (mngr *IssueManager) appendIssue(comment *lexer.Comment) error {
	id := fmt.Sprintf("%s-%d:%d", mngr.currentPath, comment.TokenStartIndex, comment.TokenEndIndex)

	issue := Issue{
		Description: comment.Description,
		EndIndex:    comment.TokenEndIndex,
		FileName:    mngr.currentBase,
		FilePath:    mngr.currentPath,
		ID:          id,
		LineNumber:  comment.LineNumber,
		OS:          mngr.os,
		StartIndex:  comment.TokenStartIndex,
		Title:       comment.Title,
	}

	if mngr.mode == IssueModeReport && mngr.template != nil {
		buf := bytes.Buffer{}
		if err := mngr.template.Execute(&buf, issue); err != nil {
			return err
		}
		issue.Body = buf.String()
	}

	mngr.Issues = append(mngr.Issues, issue)
	return nil
}

func (mngr *IssueManager) Walk(root string) error {
	ignorer, err := ignore.NewIgnorer(root)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return validateDir(d.Name(), path, ignorer)
		}

		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		ignored, err := shouldIgnore(path, ignorer)
		if err != nil {
			return err
		}

		if ignored {
			return nil
		}

		mngr.currentBase = d.Name()
		mngr.currentPath = path
		return mngr.Scan(path)
	})
}

func validateDir(dirName, path string, ignorer *ignore.Ignorer) error {
	if strings.HasPrefix(dirName, ".") {
		return filepath.SkipDir
	}

	ignored, err := shouldIgnore(path, ignorer)
	if err != nil {
		return err
	}

	if ignored {
		return filepath.SkipDir
	}

	return nil
}

func shouldIgnore(path string, ignorer *ignore.Ignorer) (bool, error) {
	if ignorer == nil {
		return false, nil
	}

	return ignorer.Match(path)
}

func (mngr *IssueManager) Scan(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	flag, err := mngr.toBitFlag()
	if err != nil {
		return err
	}

	base := lexer.NewLexer(mngr.annotation, src, mngr.currentPath, flag)
	target, err := lexer.NewTargetLexer(base)
	if err != nil {
		// @TODO create error/warning message when encountering an unsupported file extension/programming language
		return nil
	}

	tokens, err := base.AnalyzeTokens(target)
	if err != nil {
		return err
	}

	c := lexer.BuildComments(tokens)
	for _, comment := range c.Comments {
		if err := mngr.appendIssue(&comment); err != nil {
			return err
		}
	}

	return nil
}

func (mngr *IssueManager) toBitFlag() (lexer.U8, error) {
	switch mngr.mode {
	case IssueModeReport, IssueModeScan:
		return lexer.FLAG_SCAN, nil
	case IssueModePurge:
		return lexer.FLAG_PURGE, nil
	default:
		return 0, errors.New("unsupported issue mode. expected scan or purge")
	}
}

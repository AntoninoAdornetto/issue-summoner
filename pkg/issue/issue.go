/*
THIS FILE IS RESPONSIBLE FOR LOCATING, REPORTING AND MANAGING ISSUES THAT RESIDE WITHIN
SOURCE CODE COMMENTS.

ISSUE ANNOTATIONS ARE LOCATED BY WALKING [Walk] THE WORKING TREE AND SCANNING/TOKENIZING
[Scan] SOURCE CODE COMMENTS. SEE THE LEXER PACKAGE FOR DETAILS ON THE LEXICAL TOKENIZATION PROCESS.

ISSUES CAN BE REPORTED TO VARIOUS PLATFORMS, SUCH AS GITHUB, GITLAB, BITBUCKET, ECT...
ONCE REPORTED TO AN SCH, THE ID'S ASSOCIATED WITH THE ISSUE AND PLATFORM THEY WERE PUBLISHED ON
ARE APPENDED AND WRITTEN TO THE ISSUE ANNOTATION. THIS ALLOWS ISSUE-SUMMONER TO CHECK
THE STATUS OF ISSUES AND REMOVE THE COMMENT/ISSUE ENTIRELY, ONCE IT IS MARKED AS RESOLVED.

EXAMPLE PRIOR TO REPORTING THE ISSUE:
// @MY_ISSUE_ANNOTATION resolve bug....

EXAMPLE AFTER REPORTING THE ISSUE TO AN SCH:
// @MY_ISSUE_ANNOTATION(#45323) resolve bug....

# SUPPORTED MODES

- `SCAN`: LOCATES ALL SRC CODE COMMENTS THAT CONTAIN AN ISSUE [Annotation] AND STORES THE
RESULTS IN THE [Issues] SLICE.

- `REPORT`: PRODUCES THE SAME LIST OF ISSUES FROM `SCAN` MODE, AND CREATES A MAP THAT
GROUPS ISSUES TOGETHER BY FILE PATH. THE MAP IS USED AFTER ALL SELECTED [Issues] HAVE BEEN
REPORTED TO A SOURCE CODE HOSTING PLATFORM, I.E. GITHUB, GITLAB ECT..., AND WRITES THE
ISSUE ID BACK TO THE SOURCE FILE. THE ISSUE ID IS OBTAINED FROM PUBLISHING AN ISSUE,
SEE GIT PACKAGE FOR MORE DETAILS. THE GROUPING, FROM THE MAP, HELPS WITH EXECUTING 1 WRITE
CALL PER FILEPATH. MEANING, IF THERE ARE 10 ISSUES BEING REPORTED AND THEY RESIDE IN 2 SOURCE
CODE FILES, THERE WOULD ONLY BE 2 WRITE FILE CALLS BECAUSE THE ISSUES ARE GROUPED BY FILE PATH
AND BATCHED TOGETHER TO AVOID MULTIPLE WRITES. @SEE [WriteIssues] FUNC.

- `PURGE`: CHECK THE STATUS OF ISSUES THAT WERE REPORTED USING ISSUE SUMMONER AND ATTEMPTS
TO REMOVE THE COMMENTS. COMMENTS ARE REMOVED IF THE ISSUE WAS REPORTED USING <issue-summoner report>
COMMAND AND THE ISSUE ID WAS WRITTEN BACK TO THE SOURCE FILE. THE ID IS USED TO CHECK THE STATUS
AND IF RESOLVED, THE COMMENT IS REMOVED. @SEE [Purge] FUNC.
*/
package issue

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
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
	IssueMap    map[string][]IssueMapEntry
	Annotation  []byte
	currentBase string
	currentPath string
	mode        IssueMode
	os          string
	template    *template.Template
}

type Issue struct {
	ID          string // Used as a key for the multi select tui component for issue selection <issue summoner report> cmd
	Title       string // Title of the issue
	Description string // Description of the issue
	Body        string // Contains the issue body/description to use for the issue filing. Is a markdown template
	FileName    string // base
	FilePath    string // relative to the working tree dir
	LineNumber  int    // Line number of where the comment resides
	OS          string // Used for env section of the issue markdown template
	Index       int    // index of the issue in [IssueManager.Issues]
	Comment     *lexer.Comment
}

type IssueMapEntry struct {
	Index      int // index of the issue in [IssueManager.Issues]
	ReportedID int // issue identifier after calling [git.Report] func
}

// NewIssueManager accepts an annotation as input, which is used to locate issues/action
// items that are contained within comments, and a [mode] that is used to determine
// the functionality/behavior of both the issue & lexer packages. @See top comment in this
// file for a description on the supported modes and their responsibilities.
func NewIssueManager(annotation []byte, mode IssueMode) (*IssueManager, error) {
	manager := &IssueManager{
		Issues:     make([]Issue, 0),
		IssueMap:   make(map[string][]IssueMapEntry),
		Annotation: annotation,
		mode:       mode,
		os:         runtime.GOOS,
	}

	switch mode {
	case IssueModeScan:
		break
	case IssueModeReport:
		tmpl, err := generateIssueTemplate()
		if err != nil {
			return nil, err
		}
		manager.template = tmpl
	case IssueModePurge:
		manager.Annotation = append(manager.Annotation, []byte("\\(#\\d+\\)")...)
	default:
		return nil, errors.New("expected mode of \"report\", \"scan\", or \"purge\"")
	}

	return manager, nil
}

func (mngr *IssueManager) appendIssue(comment *lexer.Comment) error {
	id := fmt.Sprintf("%s-%d:%d", mngr.currentPath, comment.TokenStartIndex, comment.TokenEndIndex)

	issue := Issue{
		Description: comment.Description,
		FileName:    mngr.currentBase,
		FilePath:    mngr.currentPath,
		ID:          id,
		LineNumber:  comment.LineNumber,
		OS:          mngr.os,
		Title:       comment.Title,
		Comment:     comment,
	}

	if len(mngr.Issues) > 0 {
		issue.Index = len(mngr.Issues)
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

	base := lexer.NewLexer(mngr.Annotation, src, path, flag)
	target, err := lexer.NewTargetLexer(base)
	if err != nil {
		// @TODO create error/warning message when encountering an unsupported file extension/programming language
		return nil
	}

	tokens, err := base.AnalyzeTokens(target)
	if err != nil {
		return err
	}

	c, err := lexer.BuildComments(tokens)
	if err != nil {
		return err
	}

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

// Groups [Issues] together by file path so that when we are writting issue ids
// back to where the issue [Annotation] is located, we can do so with fewer sys calls.
func (mngr *IssueManager) Group(index, id int) error {
	if index < 0 || index > len(mngr.Issues)-1 {
		return fmt.Errorf(
			"Failed to group issues by filepath: index %d out of bounds with length of %d",
			index,
			len(mngr.Issues),
		)
	}

	current := mngr.Issues[index]
	mngr.IssueMap[current.FilePath] = append(mngr.IssueMap[current.FilePath], IssueMapEntry{
		Index:      index,
		ReportedID: id,
	})

	return nil
}

func (mngr *IssueManager) WriteIssues(pathKey string) error {
	if _, ok := mngr.IssueMap[pathKey]; !ok {
		return fmt.Errorf("File path key (%s) does not exist in issue map", pathKey)
	}

	size := len(mngr.IssueMap[pathKey])
	if size == 0 {
		return fmt.Errorf("Expected Issue map to have at least 1 entry")
	}

	srcFile, err := os.OpenFile(pathKey, os.O_RDWR, 0666)
	if err != nil {
		return nil
	}

	defer srcFile.Close()

	srcCode, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	mngr.sortPathGroup(pathKey)
	entries := mngr.IssueMap[pathKey]

	buf := bytes.Buffer{}
	for i, entry := range entries {
		comment := mngr.Issues[entry.Index].Comment
		start, end := comment.AnnotationPos[0], comment.AnnotationPos[1]

		if i == 0 {
			buf.Write(srcCode[:end+1])
		} else {
			buf.Write(srcCode[start : end+1])
		}

		buf.WriteString(fmt.Sprintf("(#%d)", entry.ReportedID))

		if i < size-1 {
			next := entries[i+1]
			nextComment := mngr.Issues[next.Index].Comment
			nextStart := nextComment.AnnotationPos[0]
			buf.Write(srcCode[end+1 : nextStart])
		} else {
			buf.Write(srcCode[end+1:])
		}
	}

	if _, err := srcFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := srcFile.Write(buf.Bytes()); err != nil {
		return err
	}

	if err := srcFile.Sync(); err != nil {
		return err
	}

	return nil
}

// sortPathGroup is needed to restore order to our [IssueMap]. This is important
// because [Issues] are reported to source code hosting platforms using go routines
// and we can't guarantee when they will finish.
// Having them in sequential order will also help during the [WriteIssues] invokation
func (mngr *IssueManager) sortPathGroup(pathKey string) {
	if group, ok := mngr.IssueMap[pathKey]; ok {
		sort.Slice(group, func(i, j int) bool {
			return group[i].Index < group[j].Index
		})
		mngr.IssueMap[pathKey] = group
	}
}

func (mngr *IssueManager) Purge(pathKey string) error {
	if _, ok := mngr.IssueMap[pathKey]; !ok {
		return fmt.Errorf("File path key (%s) does not exist in issue map", pathKey)
	}

	size := len(mngr.IssueMap[pathKey])
	if size == 0 {
		return fmt.Errorf("Expected Issue map to have at least 1 entry")
	}

	srcFile, err := os.OpenFile(pathKey, os.O_RDWR, 0666)
	if err != nil {
		return nil
	}

	defer srcFile.Close()

	srcCode, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	mngr.sortPathGroup(pathKey)
	entries := mngr.IssueMap[pathKey]

	buf := bytes.Buffer{}
	lastIndex := 0

	for _, entry := range entries {
		comment := mngr.Issues[entry.Index].Comment
		start, end := comment.NotationStartIndex, comment.NotationEndIndex
		buf.Write(srcCode[lastIndex:start])

		lastIndex = end + 1
		if srcCode[end] == '\n' {
			buf.WriteByte('\n')
		}
	}

	buf.Write(srcCode[lastIndex:])
	if err := srcFile.Truncate(0); err != nil {
		return err
	}

	if _, err := srcFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := srcFile.Write(buf.Bytes()); err != nil {
		return err
	}

	if err := srcFile.Sync(); err != nil {
		return err
	}

	return nil
}

var (
	errFailedWrite = "Issue <%s> was reported to %s but the program failed to write id %d back to the src file at path %s"
	successWrite   = "Issue <%s> successfully reported to %s and annotated with issue number %d"
)

// returns the results of reporting issues to a source code hosting platform
// and writing the ids obtained from the source code hosting platform back to the src code files
func (mngr *IssueManager) Results(pathKey, sch string, failed bool) ([]string, error) {
	if _, ok := mngr.IssueMap[pathKey]; !ok {
		return nil, fmt.Errorf("File path key (%s) does not exist in issue map", pathKey)
	}

	msgs := make([]string, len(mngr.IssueMap[pathKey]))
	for i, entry := range mngr.IssueMap[pathKey] {
		var msg string
		issue := mngr.Issues[entry.Index]
		if failed {
			msg = fmt.Sprintf(errFailedWrite, issue.Title, sch, entry.ReportedID, pathKey)
		} else {
			msg = fmt.Sprintf(successWrite, issue.Title, sch, entry.ReportedID)
		}
		msgs[i] = msg
	}

	return msgs, nil
}

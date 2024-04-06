package issue

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type PendingIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

func (pi *PendingIssue) Scan(file *os.File) error {
	lineNum := uint64(0)
	scanner := bufio.NewScanner(file)
	comment := GetCommentSymbols(filepath.Ext(file.Name()))

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	for scanner.Scan() {
		lineNum++
		currentLine := scanner.Text()
		lineType, prefix := EvalSourceLine(
			strings.TrimLeftFunc(currentLine, unicode.IsSpace),
			comment,
		)

		if lineType == LINE_TYPE_SRC_CODE {
			continue
		}

		arg := ParseCommentParams{
			LineText:      currentLine,
			LineType:      lineType,
			LineNum:       &lineNum,
			Scanner:       scanner,
			Comment:       comment,
			CommentPrefix: prefix,
			FileInfo:      fileInfo,
		}

		if err := pi.ParseComment(arg); err != nil {
			return err
		}
	}

	return nil
}

// @TODO split single and multi line parsing into different functions
// there is an error atm where we fail to build the entire description for
// multi line comments. This is because we break out of the iteration when
// the line we are parsing is not the begining or ending of a multi line commment
func (pi *PendingIssue) ParseComment(arg ParseCommentParams) error {
	issue := Issue{StartLineNumber: *arg.LineNum}
	description := strings.Builder{}

	for {
		annotated := containsAnnotation(arg.LineText, pi.Annotation)
		if annotated && issue.AnnotationLineNumber == 0 {
			pi.processAnnotationMetaData(arg.LineText, arg.LineType, arg.LineNum, &issue)
		} else if issue.AnnotationLineNumber > 0 {
			descriptionText, err := buildDescription(arg.LineText, arg.LineType, arg.CommentPrefix)
			if err != nil {
				return err
			}
			description.WriteString(descriptionText)
		}

		issue.EndLineNumber = *arg.LineNum
		if !arg.Scanner.Scan() {
			break
		}

		nextLine := arg.Scanner.Text()
		arg.LineText = nextLine
		nextLineType, _ := EvalSourceLine(
			strings.TrimLeftFunc(nextLine, unicode.IsSpace),
			arg.Comment,
		)

		*arg.LineNum++
		if nextLineType == LINE_TYPE_SRC_CODE {
			break
		}
	}

	if issue.AnnotationLineNumber > 0 {
		issue.Description = strings.TrimLeftFunc(description.String(), unicode.IsSpace)
		issue.FileInfo = arg.FileInfo
		issue.ID = fmt.Sprintf("%s-%d", arg.FileInfo.Name(), issue.AnnotationLineNumber)
		pi.Issues = append(pi.Issues, issue)
	}

	return nil
}

func (pi *PendingIssue) processAnnotationMetaData(
	line string,
	lineType string,
	lineNum *uint64,
	issue *Issue,
) {
	remainingText := strings.SplitAfter(line, pi.Annotation)[1]
	issue.Title = strings.TrimLeftFunc(remainingText, unicode.IsSpace)
	issue.AnnotationLineNumber = *lineNum
	issue.IsMultiLine = lineType == LINE_TYPE_MULTI_START || lineType == LINE_TYPE_MULTI_END
	issue.IsSingleLine = lineType == LINE_TYPE_SINGLE
}

func buildDescription(line string, lineType string, prefix string) (string, error) {
	if lineType == LINE_TYPE_SRC_CODE {
		return "", errors.New("line type should be single or multi")
	}
	return processCommentDescription(line, prefix), nil
}

// processCommentDescription returns all text after the comment prefix
func processCommentDescription(line string, prefix string) string {
	remainingText := strings.SplitAfter(line, prefix)[1]
	return strings.TrimLeftFunc(remainingText, unicode.IsSpace)
}

func containsAnnotation(line string, annotation string) bool {
	return strings.Contains(line, annotation)
}

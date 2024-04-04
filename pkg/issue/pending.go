package issue

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type PendingIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *PendingIssue) Scan(file *os.File) ([]Issue, error) {
	lineNum := uint64(0)
	scanner := bufio.NewScanner(file)
	comment := GetCommentSymbols(filepath.Ext(file.Name()))

	// @TODO utilize the file info obj
	_, err := file.Stat()
	if err != nil {
		return pi.Issues, err
	}

	for scanner.Scan() {
		lineNum++
		currentLine := scanner.Text()
		lineType := EvalSourceLine(strings.TrimLeftFunc(currentLine, unicode.IsSpace), comment)

		if lineType == LINE_TYPE_SRC_CODE {
			continue
		}

		arg := ParseCommentParams{
			LineText: currentLine,
			LineNum:  lineNum,
			LineType: lineType,
			Scanner:  scanner,
			Comment:  comment,
		}

		if err := pi.ParseComment(arg); err != nil {
			return pi.Issues, err
		}
	}

	return pi.Issues, nil
}

// ParseComment will parse lines of text that are read from a buffer and check if a given line
// contains an issue annotation. If found, details about the issue are parsed from by reading
// subsequent lines until we reach a line that is not a comment. When a line is not a comment,
// the base case clause executed and we stop scanning the buffer.
func (pi *PendingIssue) ParseComment(arg ParseCommentParams) error {
	description := strings.Builder{}
	issue := Issue{StartLineNumber: arg.LineNum}

	for {
		annotated := containsAnnotation(arg.LineText, pi.Annotation)
		if annotated && issue.AnnotationLineNumber == 0 {
			issue.Title = buildTitle(arg.LineText, pi.Annotation)
			issue.AnnotationLineNumber = arg.LineNum
			issue.IsMultiLine = arg.LineType == LINE_TYPE_MULTI
			issue.IsSingleLine = arg.LineType == LINE_TYPE_SINGLE
		} else if issue.AnnotationLineNumber > 0 {
			nextDescription := buildDescription(arg.LineText, arg.Comment)
			description.WriteString(nextDescription)
		}

		if !arg.Scanner.Scan() {
			issue.EndLineNumber = arg.LineNum
			arg.LineNum = arg.LineNum - issue.StartLineNumber
			break
		}

		arg.LineText = arg.Scanner.Text()
		nextLineType := EvalSourceLine(
			strings.TrimLeftFunc(arg.LineText, unicode.IsSpace),
			arg.Comment,
		)

		if nextLineType == LINE_TYPE_SRC_CODE {
			issue.EndLineNumber = arg.LineNum
			arg.LineNum = arg.LineNum - issue.StartLineNumber
			break
		}

		arg.LineNum++
	}

	if issue.AnnotationLineNumber > 0 {
		issue.Description = strings.TrimLeftFunc(description.String(), unicode.IsSpace)
		pi.Issues = append(pi.Issues, issue)
	}

	return arg.Scanner.Err()
}

// buildTitle takes a line and issue annotation as input and
// returns a string that removes leading white space and the
// annotation.
// Example:
// buildTitle("// @ANNOTATION comment", "@ANNOTATION") -> "comment"
func buildTitle(line string, annotation string) string {
	remainingText := strings.SplitAfter(line, annotation)[1]
	return strings.TrimLeftFunc(remainingText, unicode.IsSpace)
}

// buildDescription removes leading white space and the single line
// comment symbols and returns a new string that can be used to append
// the issue description.
// @TODO remove hardcode indexing of singleLine comment symbol
// think of a new solution to build description text without having
// to hardcode and think of how we are going to handle multi line comments
func buildDescription(line string, c Comment) string {
	remainingText := strings.SplitAfter(line, c.SingleLineSymbols[0])[1]
	return strings.TrimLeftFunc(remainingText, unicode.IsSpace)
}

func containsAnnotation(line string, tag string) bool {
	return strings.Contains(line, tag)
}

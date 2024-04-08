package issue

import (
	"bufio"
	"os"
	"path/filepath"
)

type PendingIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *PendingIssue) GetIssues() []Issue {
	return pi.Issues
}

func (pi *PendingIssue) Scan(file *os.File) error {
	issue := &Issue{}
	lineNum := uint64(0)
	scanner := bufio.NewScanner(file)
	comment := CommentPrefixes(filepath.Ext(file.Name()))

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	for scanner.Scan() {
		lineNum++
		currentLine := scanner.Text()
		comment.SetLineTypeAndPrefix(currentLine)

		if comment.CurrentLineType == LINE_TYPE_SRC_CODE {
			continue
		}

		exitOn := LINE_TYPE_SRC_CODE
		issue.Init(comment.CurrentLineType, lineNum, &fileInfo)
		if issue.IsMultiLine {
			exitOn = LINE_TYPE_MULTI_END
		}

		linesScanned := pi.ProcessComment(issue, scanner, &comment, exitOn)
		lineNum += linesScanned
		if issue.AnnotationLineNumber > 0 {
			pi.Issues = append(pi.Issues, *issue)
		}

		issue = &Issue{}
	}

	return scanner.Err()
}

func (pi *PendingIssue) ProcessComment(
	is *Issue,
	s *bufio.Scanner,
	c *Comment,
	exitLine string,
) uint64 {
	scannedLines := uint64(0)
	currentLine := s.Text()
	content, isAnnotated := c.ExtractCommentContent(currentLine, pi.Annotation)
	is.Build(content, isAnnotated, scannedLines)

	for s.Scan() {
		scannedLines++
		nextLine := s.Text()
		content, isAnnotated = c.ExtractCommentContent(nextLine, pi.Annotation)
		if c.CurrentLineType == exitLine {
			break
		}

		is.Build(content, isAnnotated, scannedLines)
	}

	return scannedLines
}

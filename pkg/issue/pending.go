package issue

import (
	"bufio"
	"fmt"
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

		switch comment.CurrentLineType {
		case LINE_TYPE_SRC_CODE:
			continue
		case LINE_TYPE_SINGLE:
			issue.Init(comment.CurrentLineType, lineNum, &fileInfo)
			lineNum += pi.ProcessComment(issue, scanner, &comment, LINE_TYPE_SRC_CODE)
		case LINE_TYPE_MULTI_START:
			issue.Init(comment.CurrentLineType, lineNum, &fileInfo)
			lineNum += pi.ProcessComment(issue, scanner, &comment, LINE_TYPE_MULTI_END)
		}

		if issue.AnnotationLineNumber > 0 {
			pi.Issues = append(pi.Issues, *issue)
		}
		//
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

	if isAnnotated {
		is.AnnotationLineNumber = is.StartLineNumber + scannedLines
		is.Title = content
	} else if is.Description == "" {
		is.Description = content
		is.EndLineNumber = scannedLines + is.StartLineNumber
	} else {
		is.Description = fmt.Sprintf("%s %s", is.Description, content)
		is.EndLineNumber = scannedLines + is.StartLineNumber
	}

	for s.Scan() {
		scannedLines++
		nextLine := s.Text()
		content, isAnnotated = c.ExtractCommentContent(nextLine, pi.Annotation)
		if c.CurrentLineType == exitLine {
			break
		}

		if isAnnotated {
			is.AnnotationLineNumber = is.StartLineNumber + scannedLines
			is.Title = content
		} else if is.Description == "" {
			is.Description = content
			is.EndLineNumber = scannedLines + is.StartLineNumber
		} else {
			is.Description = fmt.Sprintf("%s %s", is.Description, content)
			is.EndLineNumber = scannedLines + is.StartLineNumber
		}
	}

	return scannedLines
}

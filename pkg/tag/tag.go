package tag

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Tag struct {
	Description       string
	StartLineNumber   uint64
	EndLineNumber     uint64
	AnnotationLineNum uint64
	LineNumber        uint64
	FileInfo          os.FileInfo
}

type TagManager struct {
	TagName string
	Mode    string
}

type PendedTagManager struct {
	TagManager
}

type PendedTagParser interface {
	ScanForTags(ScanForTagsParams) ([]Tag, error)
}

type IssuedTagManager struct {
	TagManager
}

const (
	IssueMode   string = "I"
	PendingMode string = "P"
)

func (tm *TagManager) ValidateMode(mode string) error {
	switch mode {
	case IssueMode, PendingMode:
		return nil
	default:
		return fmt.Errorf(
			"Error: mode %s is invalid\nI (Issue Mode) & P (Pending Mode) are the available options",
			mode,
		)
	}
}

type ScanForTagsParams struct {
	Path     string
	File     *os.File
	FileInfo os.FileInfo
}

func (pm *PendedTagManager) ScanForTags(detail ScanForTagsParams) ([]Tag, error) {
	tags := make([]Tag, 0)
	lineNum := uint64(0)
	scanner := bufio.NewScanner(detail.File)
	fileExtension := filepath.Ext(detail.File.Name())

	for scanner.Scan() {
		tagDescription := extractAnnotationDetails(scanner, fileExtension, pm.TagName)
		if tagDescription != "" {
			tags = append(tags, Tag{
				LineNumber:  lineNum + 1,
				FileInfo:    detail.FileInfo,
				Description: tagDescription,
			})
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return tags, nil
}

func extractAnnotationDetails(scanner *bufio.Scanner, ext string, tag string) string {
	builder := strings.Builder{}
	text := scanner.Text()
	trimmedText := strings.TrimSpace(text)
	hasAnnotation := strings.Contains(trimmedText, tag)
	commentSyntax := CommentSyntax(ext)

	if strings.HasPrefix(trimmedText, commentSyntax.SingleLineCommentSymbols) && hasAnnotation {
		description := strings.Join(strings.SplitAfter(trimmedText, tag), "")
		builder.WriteString(description)
	}

	return builder.String()
}

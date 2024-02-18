package tag

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Tag struct {
	Description       string
	StartLineNumber   uint64
	EndLineNumber     uint64
	AnnotationLineNum uint64
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
	commentSyntaxMp := CommentSyntax(fileExtension)

	for scanner.Scan() {
		tag := Tag{}
		builder := strings.Builder{}
		text := scanner.Text()
		lineNum++

		if len(text) == 0 || isAlphaNumeric(rune(text[0])) {
			continue
		}

		currentLine := strings.TrimLeftFunc(text, func(r rune) bool {
			return unicode.IsSpace(r)
		})

		hasAnnotation := hasTagAnnotation(currentLine, pm.TagName)

		if isSingleLineComment(currentLine, commentSyntaxMp) {
			tag.StartLineNumber = lineNum
			if hasAnnotation {
				tag.AnnotationLineNum = lineNum
				builder.WriteString(
					strings.Join(strings.SplitAfter(currentLine, pm.TagName)[1:], ""),
				)
			}

			for scanner.Scan() {
				currentLine = strings.TrimLeftFunc(scanner.Text(), func(r rune) bool {
					return unicode.IsSpace(r)
				})

				if !isSingleLineComment(currentLine, commentSyntaxMp) {
					break
				}

				if !hasAnnotation {
					hasAnnotation = hasTagAnnotation(currentLine, pm.TagName)
				}

				builder.WriteString("\n")
				nextLineDescription := strings.Join(
					strings.SplitAfter(currentLine, commentSyntaxMp.SingleLineCommentSymbols)[1:],
					"",
				)

				builder.WriteString(strings.TrimSpace(nextLineDescription))
				lineNum++
			}

			tag.EndLineNumber = lineNum
			tag.FileInfo = detail.FileInfo
			tag.Description = strings.TrimSpace(builder.String())
			if hasAnnotation {
				tags = append(tags, tag)
				lineNum++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return tags, nil
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func trimWhiteSpace(s string) string {
	return strings.TrimLeft(s, " ")
}

func isSingleLineComment(line string, commentSyntaxMap CommentLangSyntax) bool {
	return strings.HasPrefix(line, commentSyntaxMap.SingleLineCommentSymbols)
}

func hasTagAnnotation(line string, annotation string) bool {
	return strings.Contains(line, annotation)
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

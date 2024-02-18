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
	Title             string
	Description       string
	StartLineNumber   uint64
	EndLineNumber     uint64
	AnnotationLineNum uint64
	IsSingleLine      bool
	IsMultiLine       bool
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
		return fmt.Errorf("Invalid mode %s. See issue-summoner --help for more information", mode)
	}
}

type ScanForTagsParams struct {
	Path     string
	File     *os.File
	FileInfo os.FileInfo
}

func (pm *PendedTagManager) ScanForTags(
	path string,
	file *os.File,
	info os.FileInfo,
) ([]Tag, error) {
	tags := make([]Tag, 0)
	lineNum := uint64(0)
	scanner := bufio.NewScanner(file)
	commentSyntax := CommentSyntax(filepath.Ext(file.Name()))

	for scanner.Scan() {
		text := scanner.Text()
		lineNum++

		if shouldSkip(text) {
			continue
		}

		if isSingleLineComment(text, commentSyntax) {
			tag, linesScanned := pm.parseSingleLineCommentBlock(
				scanner,
				text,
				lineNum,
				commentSyntax,
			)
			if tag != nil {
				tag.FileInfo = info
				tags = append(tags, *tag)
			}
			lineNum += linesScanned
		} else if isMultiLineCommentStart(text, commentSyntax) {
			tag, linesScanned := pm.parseMultiLineCommentBlock(scanner, text, lineNum, commentSyntax)
			if tag != nil {
				tag.FileInfo = info
				tags = append(tags, *tag)
			}
			lineNum += linesScanned
		}
	}

	return tags, scanner.Err()
}

func (pm *PendedTagManager) parseSingleLineCommentBlock(
	scanner *bufio.Scanner,
	text string,
	lineNum uint64,
	cs CommentLangSyntax,
) (*Tag, uint64) {
	tag := &Tag{StartLineNumber: lineNum}
	var description strings.Builder
	linesScanned := uint64(0)

	for {
		if annotated := hasTagAnnotation(text, pm.TagName); annotated &&
			tag.AnnotationLineNum == 0 {
			tag.AnnotationLineNum = lineNum
			tag.Title = strings.TrimSpace(strings.SplitN(text, pm.TagName, 2)[1])
			tag.IsSingleLine = true
		} else if tag.AnnotationLineNum > 0 {
			nextDescription := strings.Join(strings.SplitAfter(text, cs.SingleLineCommentSymbols)[1:], "")
			description.WriteString(nextDescription)
		}

		if !scanner.Scan() {
			tag.EndLineNumber = lineNum
			break
		}

		text = scanner.Text()
		lineNum++
		linesScanned++

		if !isSingleLineComment(text, cs) {
			tag.EndLineNumber = lineNum - 1
			break
		}
	}

	if tag.AnnotationLineNum > 0 {
		tag.Description = strings.TrimSpace(description.String())
		return tag, linesScanned
	}

	return nil, linesScanned
}

func (pm *PendedTagManager) parseMultiLineCommentBlock(
	scanner *bufio.Scanner,
	text string,
	lineNum uint64,
	cs CommentLangSyntax,
) (*Tag, uint64) {
	tag := &Tag{StartLineNumber: lineNum}
	var description strings.Builder
	linesScanned := uint64(0)

	for {
		if annotated := hasTagAnnotation(text, pm.TagName); annotated &&
			tag.AnnotationLineNum == 0 {
			tag.AnnotationLineNum = lineNum
			tag.Title = strings.TrimSpace(strings.SplitN(text, pm.TagName, 2)[1])
			tag.IsMultiLine = true
		} else {
			trimmedText := strings.TrimSpace(text)
			newDescription := strings.TrimPrefix(trimmedText, cs.MultiLineCommentSymbols.CommentStartSymbol)
			description.WriteString(fmt.Sprintf(" %s", strings.TrimSpace(newDescription)))
		}

		if !scanner.Scan() {
			tag.EndLineNumber = lineNum
			break
		}

		text = scanner.Text()
		lineNum++
		linesScanned++

		if isMultiLineCommentEnd(text, cs) {
			tag.EndLineNumber = lineNum
			break
		}
	}

	if tag.AnnotationLineNum > 0 {
		tag.Description = strings.TrimSpace(description.String())
		return tag, linesScanned
	}

	return nil, linesScanned
}

func shouldSkip(line string) bool {
	return len(line) == 0 || isAlphaNumeric(rune(line[0]))
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func isSingleLineComment(line string, commentSyntax CommentLangSyntax) bool {
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	return strings.HasPrefix(trimmed, commentSyntax.SingleLineCommentSymbols)
}

func isMultiLineCommentStart(line string, commentSyntax CommentLangSyntax) bool {
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	return strings.HasPrefix(trimmed, commentSyntax.MultiLineCommentSymbols.CommentStartSymbol)
}

func isMultiLineCommentEnd(line string, commentSyntax CommentLangSyntax) bool {
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	return strings.HasSuffix(trimmed, commentSyntax.MultiLineCommentSymbols.CommentEndSymbol)
}

func hasTagAnnotation(line string, annotation string) bool {
	return strings.Contains(line, annotation)
}

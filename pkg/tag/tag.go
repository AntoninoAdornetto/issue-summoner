package tag

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
)

type Tag struct {
	LineNum  uint64
	FileInfo os.FileInfo
}

type TagManager struct {
	TagName string
	Mode    string
}

type PendedTagManager struct {
	TagManager
}

type PendedTagParser interface {
	FindTags(path string, fileOperator TagFileOperator) ([]Tag, error)
	CompileSingleLineComment(fileInfo fs.FileInfo) regexp.Regexp
}

type IssuedTagManager struct {
	TagManager
}

// type IssuedTagParser interface {
// 	FindTags(path string, fileOperator TagFileOperator) ([]Tag, error)
// 	CompileSingleLineComment(fileInfo fs.FileInfo) regexp.Regexp // @TODO - Implement
// }

type TagFileOperator interface {
	Open(fileName string) (*os.File, error)
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

/*
@TODO - Implement
PendedTagManagers `FindTags` method will search for tags that have not been reported to a source code manager. The function will account for both
single line & multi line comment syntax.
*/
func (pm *PendedTagManager) FindTags(path string, fileOperator TagFileOperator) ([]Tag, error) {
	tags := make([]Tag, 0)

	file, err := fileOperator.Open(path)
	if err != nil {
		return tags, fmt.Errorf("Error: failed to open file %s\n%s", path, err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return tags, fmt.Errorf(
			"Error: failed to get file information for %s\n%s",
			file.Name(),
			err,
		)
	}

	lineNum := uint64(0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		singleLine := pm.CompileSingleLineComment()

		if singleLine.Match(scanner.Bytes()) {
			tags = append(tags, Tag{
				LineNum:  lineNum + 1,
				FileInfo: fileInfo,
			})
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return tags, fmt.Errorf("Error: failed to scan file %s\n%s", file.Name(), err)
	}

	return tags, nil
}

func (pm *PendedTagManager) CompileSingleLineComment() regexp.Regexp {
	// @TODO - Utilize the constants in comment.go to properly build an expression based on the file extension
	return *regexp.MustCompile(fmt.Sprintf("^//(.*)%s(.*)$", pm.TagName))
}

/*
@TODO - Implement
IssuedTagManager `FindTags` method will search for tags that have been reported to a source code manager. The function will account for both
single line & multi line comment syntax.
*/
// func (im *IssuedTagManager) FindTags(path string) ([]Tag, error) {
// 	tags := make([]Tag, 0)
// 	return tags, nil
// }

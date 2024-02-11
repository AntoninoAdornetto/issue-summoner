package tag

import (
	"fmt"
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
	FindTags(path string) ([]Tag, error)
}

type IssuedTagManager struct {
	TagManager
}

type IssuedTagParser interface {
	FindTags(path string) ([]Tag, error)
	CompileSingleLineComment() regexp.Regexp // @TODO - Implement
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
func (pm *PendedTagManager) FindTags(path string) ([]Tag, error) {
	tags := make([]Tag, 0)
	return tags, nil
}

/*
@TODO - Implement
IssuedTagManager `FindTags` method will search for tags that have been reported to a source code manager. The function will account for both
single line & multi line comment syntax.
*/
func (im *IssuedTagManager) FindTags(path string) ([]Tag, error) {
	tags := make([]Tag, 0)
	return tags, nil
}

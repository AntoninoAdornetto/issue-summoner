package issue

import (
	"bufio"
	"errors"
	"regexp"
)

type CommentManager interface {
	ParseComment(startLineNumber int) ([]Issue, error)
}

type NewCommentManagerParams struct {
	Annotation      string
	CommentType     string
	FilePath        string
	FileName        string
	StartLineNumber int
	Locations       []int
	PrefixRe        *regexp.Regexp
	SuffixRe        *regexp.Regexp
	Scanner         *bufio.Scanner
}

func NewCommentManager(arg NewCommentManagerParams) (CommentManager, error) {
	switch arg.CommentType {
	case SINGLE_LINE_COMMENT:
		return &SingleLineComment{
			Annotation:               arg.Annotation,
			AnnotationIndicator:      false,
			FileName:                 arg.FileName,
			FilePath:                 arg.FilePath,
			PrefixRe:                 arg.PrefixRe,
			SuffixRe:                 arg.SuffixRe,
			Scanner:                  arg.Scanner,
			CommentNotationLocations: arg.Locations,
		}, nil
	default:
		return nil, errors.New("unsupported comment type. single/multi line comments are supported")
	}
}

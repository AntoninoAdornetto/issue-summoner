package issue

import (
	"bufio"
	"errors"
	"regexp"
)

type Comment struct {
	ID                   string
	Title                string
	Description          string
	FileName             string
	FilePath             string
	StartLineNumber      int
	EndLineNumber        int
	AnnotationLineNumber int
	ColumnLocations      [][]int
}

type CommentManager interface {
	ParseComment(startLineNumber int) ([]Comment, error)
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

	default:
	}
}

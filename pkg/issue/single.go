package issue

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"unicode"
)

type SingleLineComment struct {
	Annotation               string
	AnnotationIndicator      bool
	FileName                 string
	FilePath                 string
	PrefixRe                 *regexp.Regexp
	SuffixRe                 *regexp.Regexp
	Scanner                  *bufio.Scanner
	CommentNotationLocations []int
}

func (slc *SingleLineComment) ParseComment(start int) ([]Comment, error) {
	comments := make([]Comment, 0)
	line := slc.Scanner.Bytes()
	cutIndices, err := slc.FindCutLocations(line)
	if err != nil {
		return nil, err
	}

	line, err = slc.Slice(line, cutIndices)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	fields := bytes.Fields(line)
	for _, field := range fields {
		if _, err := slc.Write(&buf, field); err != nil {
			return nil, err
		}
	}

	if comment := slc.NewComment(&buf, start); comment != nil {
		comment.ColumnLocations = [][]int{cutIndices}
		comments = append(comments, *comment)
		return comments, nil
	}

	return nil, nil
}

func (slc *SingleLineComment) FindCutLocations(line []byte) ([]int, error) {
	if len(slc.CommentNotationLocations) != 2 {
		return nil, errors.New("single line comment must contain a start & end index location")
	}

	start := slc.CommentNotationLocations[1]
	cutIndices := []int{start}

	suffixIndices := slc.FindSuffixLocations(line)
	if suffixIndices != nil {
		end := suffixIndices[0]
		cutIndices = append(cutIndices, end)
		return suffixIndices, nil
	}

	return cutIndices, nil
}

func (slc *SingleLineComment) FindSuffixLocations(line []byte) []int {
	if slc.SuffixRe == nil {
		return nil
	}
	return slc.SuffixRe.FindSubmatchIndex(line)
}

func (slc *SingleLineComment) Slice(line []byte, indices []int) ([]byte, error) {
	switch len(indices) {
	case 1:
		start := indices[0]
		return bytes.TrimLeftFunc(line[start:], unicode.IsSpace), nil
	case 2:
		start, end := indices[0], indices[1]
		return bytes.TrimLeftFunc(line[start:end], unicode.IsSpace), nil
	default:
		return nil, errors.New("slice indices should have a len of 1 or 2 (start & end)")
	}
}

func (slc *SingleLineComment) Write(wr io.Writer, field []byte) (int, error) {
	switch {
	case !slc.AnnotationIndicator && bytes.Contains(field, []byte(slc.Annotation)):
		slc.AnnotationIndicator = true
		return 0, nil
	case slc.AnnotationIndicator:
		field = append(field, ' ')
		if _, err := wr.Write(field); err != nil {
			return 0, err
		}
		return len(field), nil
	default:
		return 0, nil
	}
}

func (slc *SingleLineComment) NewComment(buf *bytes.Buffer, lineNumber int) *Comment {
	if !slc.AnnotationIndicator || buf.Len() == 0 {
		return nil
	}

	title := buf.String()[:buf.Len()-1]
	return &Comment{
		ID:                   fmt.Sprintf("%s-%d", slc.FileName, lineNumber),
		Title:                title,
		Description:          "",
		FileName:             slc.FileName,
		FilePath:             slc.FilePath,
		StartLineNumber:      lineNumber,
		EndLineNumber:        lineNumber,
		AnnotationLineNumber: lineNumber,
	}
}

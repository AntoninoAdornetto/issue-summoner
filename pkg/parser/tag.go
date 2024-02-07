package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Tag struct {
	TagName     string
	LineNum     uint64
	FileInfo    os.FileInfo
	FileExt     string
	IsProcessed bool
}

type TagOpener interface {
	Open(fileName string) (*os.File, error)
}

type TagParser interface {
	FindTags(path string, re regexp.Regexp, fo TagOpener) ([]Tag, error)
}

func FindTags(path string, re regexp.Regexp, fo TagOpener) ([]Tag, error) {
	tags := make([]Tag, 0)

	file, err := fo.Open(path)
	if err != nil {
		return tags, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return tags, err
	}

	defer file.Close()

	lineNum := uint64(0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if re.Match(scanner.Bytes()) {
			tags = append(tags, Tag{
				LineNum:  lineNum + 1,
				FileInfo: fileInfo,
				FileExt:  filepath.Ext(fileInfo.Name()),
			})
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return tags, err
	}

	return tags, nil
}

func CompileTagRegexp(tagName string) regexp.Regexp {
	return *regexp.MustCompile(fmt.Sprintf("^(.*)%s(.*)$", tagName))
}

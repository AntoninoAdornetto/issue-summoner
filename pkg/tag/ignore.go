package tag

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	GitIgnoreFile string = ".gitignore"
)

type (
	GitIgnorePattern = regexp.Regexp
)

type IgnoreFileOpener interface {
	Open(fileName string) (*os.File, error)
}

/*
Takes a path (`gitIgnorePath`) and a file opener (`fo`) as input and returns a slice of compiled regular expressions -
that represent the patterns found in the .gitignore file.

The file opener is an interface that allows the caller to provide a custom implementation for opening files. see `utils/os_file_opener.go`

The function reads the .gitignore file line by line and compiles each line into a regular expression. Empty lines and
comments (`#`) are ignored.
*/
func ProcessIgnorePatterns(gitIgnorePath string, fo IgnoreFileOpener) ([]GitIgnorePattern, error) {
	patterns := make([]GitIgnorePattern, 0)

	file, err := fo.Open(gitIgnorePath)
	if err != nil {
		return patterns, fmt.Errorf(
			"Error: failed to open git ignore file path: (%s)\n%w",
			gitIgnorePath,
			err,
		)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := scanner.Text()

		if len(pattern) == 0 || strings.HasPrefix(pattern, "#") {
			continue
		}

		patterns = append(patterns, *regexp.MustCompile(formatIgnore(pattern)))
	}

	if err := scanner.Err(); err != nil {
		return patterns, fmt.Errorf("Error: failed to scan gitignore file. %w", err)
	}

	file.Close()
	return patterns, nil
}

/*
Transforms a single .gitignore pattern into a string that can be used to compile a regular expression.
*/
func formatIgnore(pattern string) string {
	res := strings.Builder{}

	if strings.HasPrefix(pattern, "/") {
		res.WriteString("\\")
	}

	for _, v := range pattern {
		switch v {
		case '*':
			res.WriteRune('.')
			res.WriteRune('*')
		case '.':
			if strings.HasSuffix(res.String(), "*") {
				res.WriteRune('\\')
			}
			res.WriteRune('.')
		default:
			res.WriteRune(v)
		}
	}

	return res.String()
}

/*
The functions in this file are responsible for building a slice of regular expressions
base on patterns in a .gitignore file. The patterns are used to help the program adhere
to the same rules that a .gitignore pattern applies to git repos. The result, we do not
parse source code that does not need to be parsed. For example, we would never want to
parse a god forsaken node modules folder!
*/
package scm

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"unicode"
)

type IgnorePattern = regexp.Regexp

/*
@TODO add ! (not) operator support for ignoring specific files/directories.

The ParseIgnorePatterns can handle most of the common patterns found in a gitignore
file. However, there are scenarios where this function will fail to build proper regexps.
Here is an example of some patterns that are not yet supported:

1. Ignore files in a specific directory, but not its subdirectories:
directory_to_ignore/*
!directory_to_ignore/*

2. Ignore files in a specific directory, except for one specific file:
directory_to_ignore/*
!directory_to_ignore/exception_file.txt

3. Ignore all files in a directory, including hidden files:
directory_to_ignore/**
!directory_to_ignore/
!directory_to_ignore/*
*/
func ParseIgnorePatterns(r io.Reader) ([]IgnorePattern, error) {
	regexps := make([]IgnorePattern, 0)
	buf := &bytes.Buffer{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		if err := writeIgnoreRegexpBytes(buf, line); err != nil {
			return regexps, err
		}

		if buf.Len() > 0 {
			re := regexp.MustCompile(buf.String())
			regexps = append(regexps, *re)
			buf.Reset()
		}
	}

	return regexps, scanner.Err()
}

func writeIgnoreRegexpBytes(buf *bytes.Buffer, b []byte) error {
	b = bytes.TrimLeftFunc(b, unicode.IsSpace)
	if err := prependExpression(buf, b); err != nil {
		return err
	}

	for _, char := range b {
		if char == '\n' || char == '#' {
			return nil
		}

		if err := writeAndCheck(buf, []byte{char}); err != nil {
			return err
		}
	}

	return appendExpression(buf)
}

func prependExpression(buf *bytes.Buffer, b []byte) error {
	if len(b) == 0 {
		return nil
	}

	first := b[0]
	switch first {
	case '/', '\\':
		// prevent unescaped & dangling backslash error when compiling the expression
		return writeEscapeChar(buf)
	case '*':
		// write the dot operator to match any character before the quantifier
		// the quantifer matches 0 or more of the preceeding token
		return matchAnyChar(buf)
	default:
		return nil
	}
}

func appendExpression(buf *bytes.Buffer) error {
	if buf.Len() == 0 {
		return nil
	}

	last := buf.Bytes()[buf.Len()-1]
	switch last {
	case '*':
		return matchAnyChar(buf)
	case '\\':
		return writeEscapeChar(buf)
	default:
		return nil
	}
}

func writeEscapeChar(buf *bytes.Buffer) error {
	return writeAndCheck(buf, []byte(`\`))
}

// matchAnyChar writes a dot (.) byte to the buffer
// dot (.) matches any character except line breaks
func matchAnyChar(buf *bytes.Buffer) error {
	return writeAndCheck(buf, []byte("."))
}

func writeAndCheck(buf *bytes.Buffer, b []byte) error {
	_, err := buf.Write(b)
	return err
}

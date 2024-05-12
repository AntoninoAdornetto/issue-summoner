/*
The functions in this file are responsible for building a slice of regular expressions
based on patterns in a .gitignore file. The patterns are used to help the program adhere
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
@TODO broken gitignore pattern parsing
ParseIgnorePatterns function is used to determine if a file should be scanned when walking
the project dir. I've discovered issues when testing the program on large open source projects,
such Elasticsearch (java), and Pandas (python). What I found is that some gitignore patterns will
either cause files to be scanned, when they shouldn't be, or ignore files that should be scanned.
A decision will need to be made about if a library will be utilized or rewrite the implementation
to account for all the different edgecases. We may not need to reinvent the wheel here.
go-ignore is an option but it's quite old, around 10 years old I believe. Check for other libs.

For the meantime, I will use regexp.Compile over MustCompile. If an error is encountered,
we skip through the remaining lines of the gitignore file.
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
			re, err := regexp.Compile(buf.String())
			if err != nil {
				buf.Reset()
				continue
			}
			regexps = append(regexps, *re)
			buf.Reset()
		}
	}

	return regexps, scanner.Err()
}

func writeIgnoreRegexpBytes(buf *bytes.Buffer, b []byte) error {
	b = bytes.TrimLeftFunc(b, unicode.IsSpace)
	if b[0] == '!' {
		return nil
	}

	if err := prependExpression(buf, b); err != nil {
		return err
	}

	for _, char := range b {
		if char == '\n' || char == '#' {
			return nil
		}

		if repeated(buf, char) {
			continue
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

func repeated(buf *bytes.Buffer, b byte) bool {
	return b == '*' && buf.Len() > 0 && buf.Bytes()[buf.Len()-1] == '*'
}

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
@TODO .gitignore path validation rules are broken
Throughout some testing on large open source projects, I have found the effectiveness of-
ignoring paths largely unsuccessful. There are edge cases we need to bake into these functions:

1. (!) Prefix to negate patterns. Matching files excluded by previous patters will become included again.

2. (/) Separator placement. Beginning, middle (or both) means the path pattern is relative to the directory level
of the gitignore file. If the separator is at the end, the pattern should match directories or files

3. (*) Asterisk matches anything except a slash.

4. (?) Matches any one character except "/"

5. [a-zA-Z] Range notation can be used to match one of the characters in a range

6. (**) Leading/Trailing Double Asterisks

7. Sub directories may contain their own gitignore file. Account for this.

This initial implementation was quick and dirty and it worked for small projects. However, it is now causing issues
and may lead to some files/directories not being scanned at all or scanning files/dirs that shouldn't be scanned.
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

/*
This file provides constants and structures for determining comment syntax in source code files.
The constants define file extensions while the CommentNotation struct specifies
syntax details for various programming languages. The NewCommentNotation function retrieves the
appropriate CommentNotation based on a file extension.

The CommentNotations map contains predefined syntax for common languages such as C, Python, and Markdown,
with a default syntax for unrecognized file types. When walking a project directory, the program reads
each file extension and uses the CommentNotations map to determine the comment syntax for parsing.
*/
package issue

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
)

const (
	file_ext_asm        = ".asm"
	file_ext_bash       = ".sh"
	file_ext_cpp        = ".cpp"
	file_ext_c          = ".c"
	file_ext_ch         = ".h"
	file_ext_cs         = ".cs"
	file_ext_go         = ".go"
	file_ext_haskell    = ".hs"
	file_ext_html       = ".html"
	file_ext_jai        = ".jai"
	file_ext_java       = ".java"
	file_ext_javascript = ".js"
	file_ext_jsx        = ".jsx"
	file_ext_kotlin     = ".kt"
	file_ext_lisp       = ".lisp"
	file_ext_lua        = ".lua"
	file_ext_obj_c      = ".m"
	file_ext_ocaml      = ".ml"
	file_ext_php        = ".php"
	file_ext_python     = ".py"
	file_ext_ruby       = ".rb"
	file_ext_rust       = ".rs"
	file_ext_markdown   = ".md"
	file_ext_r          = ".r"
	file_ext_scala      = ".scala"
	file_ext_swift      = ".swift"
	file_ext_typescript = ".ts"
	file_ext_tsx        = ".tsx"
	file_ext_vim        = ".vim"
	file_ext_zig        = ".zig"
	errStackUnderflow   = "error: notation stack underflow"
)

type CommentNotation struct {
	Annotation          string
	AnnotationIndicator bool
	SingleLinePrefix    string
	SingleLinePrefixRe  *regexp.Regexp
	SingleLineSuffix    string
	SingleLineSuffixRe  *regexp.Regexp
	MultiLinePrefix     string
	MultiLinePrefixRe   *regexp.Regexp
	MultiLineSuffix     string
	MultiLineSuffixRe   *regexp.Regexp
	NewLinePrefix       string
	Scanner             *bufio.Scanner
	Stack               NotationStack
}

type NotationStack struct {
	Items []string
	Top   int
}

var CommentNotations = map[string]CommentNotation{
	file_ext_c: {
		SingleLinePrefix: `\/\/`,
		MultiLinePrefix:  `\/\*`,
		MultiLineSuffix:  `\*\/`,
		NewLinePrefix:    `*`,
	},
	file_ext_python: {
		SingleLinePrefix: `#`,
		MultiLinePrefix:  `(\"\"\")|(\'\'\')`,
		MultiLineSuffix:  `(\"\"\")|(\'\'\')`,
	},
	file_ext_markdown: {
		SingleLinePrefix: `<!--`,
		SingleLineSuffix: `-->`,
	},
	"default": {
		SingleLinePrefix: `#`,
	},
}

func NewCommentNotation(ext string, annotation string, scanner *bufio.Scanner) CommentNotation {
	var cn CommentNotation

	switch ext {
	case file_ext_c,
		file_ext_cpp,
		file_ext_java,
		file_ext_javascript,
		file_ext_jsx,
		file_ext_typescript,
		file_ext_tsx,
		file_ext_cs,
		file_ext_go,
		file_ext_php,
		file_ext_swift,
		file_ext_kotlin,
		file_ext_rust,
		file_ext_obj_c,
		file_ext_scala:
		cn = CommentNotations[file_ext_c]
	case file_ext_python:
		cn = CommentNotations[file_ext_python]
	case file_ext_markdown:
		cn = CommentNotations[file_ext_markdown]
	default:
		cn = CommentNotations["default"]
	}

	cn.Annotation = annotation
	cn.Scanner = scanner
	cn.Stack = *InitNotationStack()
	cn.SingleLinePrefixRe = compileAndSetRegexp(cn.SingleLinePrefix)
	cn.SingleLineSuffixRe = compileAndSetRegexp(cn.SingleLineSuffix)
	cn.MultiLinePrefixRe = compileAndSetRegexp(cn.MultiLinePrefix)
	cn.MultiLineSuffixRe = compileAndSetRegexp(cn.MultiLineSuffix)

	return cn
}

func (c *CommentNotation) ParseLine(n *uint64) (Issue, error) {
	var err error
	issue := Issue{}
	fields := strings.Fields(c.Scanner.Text())
	if len(fields) == 0 {
		return issue, nil
	}

	start := -1
	if c.Stack.IsEmpty() {
		if start = c.FindPrefixIndex(fields); start == -1 {
			return issue, nil
		}
	}

	top, err := c.Stack.Peek()
	if err != nil {
		return issue, err
	}

	if top == c.SingleLinePrefix {
		err = c.BuildSingle(fields, &issue, start+1, n)
	}

	if top == c.MultiLinePrefix {
		issue.StartLineNumber = *n
		for c.Scanner.Scan() && !c.Stack.IsEmpty() {
			err = c.BuildMulti(fields, &issue, start+1, n)
			fields = strings.Fields(c.Scanner.Text())
			start = -1
			c.AnnotationIndicator = false
			*n++
		}
	}

	c.AnnotationIndicator = false
	return issue, err
}

func (c *CommentNotation) BuildSingle(fields []string, is *Issue, start int, n *uint64) error {
	content, err := c.ExtractFromSingle(fields, start)
	if err != nil {
		return err
	}

	if c.AnnotationIndicator {
		is.Title = content
		is.Description = ""
		is.AnnotationLineNumber = *n
		is.StartLineNumber = *n
		is.EndLineNumber = *n
	}

	return nil
}

func (c *CommentNotation) ExtractFromSingle(fields []string, start int) (string, error) {
	end := len(fields)

	for i := start; i < len(fields); i++ {
		if fields[i] == c.Annotation {
			c.AnnotationIndicator = true
			start = i + 1
		}

		if c.SingleLineSuffixRe != nil && c.SingleLineSuffixRe.MatchString(fields[i]) {
			end = i
		}
	}

	_, err := c.Stack.Pop()
	return strings.Join(fields[start:end], " "), err
}

func (c *CommentNotation) BuildMulti(fields []string, is *Issue, start int, n *uint64) error {
	content, err := c.ExtractFromMulti(fields, start)
	if err != nil {
		return err
	}

	if content == "" {
		return nil
	}

	if c.AnnotationIndicator && is.Title == "" {
		is.Title = content
		is.AnnotationLineNumber = *n
	} else if is.Description == "" {
		is.Description = content
	} else {
		is.Description += " " + content
	}

	is.EndLineNumber = *n
	return nil
}

func (c *CommentNotation) ExtractFromMulti(fields []string, start int) (string, error) {
	var err error
	end := len(fields)

	if len(fields) == 0 {
		return "", nil
	}

	if fields[0] == c.NewLinePrefix {
		start++
	}

	for i := start; i < len(fields); i++ {
		if fields[i] == c.Annotation {
			c.AnnotationIndicator = true
			start = i + 1
		}

		if c.MultiLineSuffixRe.MatchString(fields[i]) {
			_, err = c.Stack.Pop()
			end = i
		}
	}

	return strings.Join(fields[start:end], " "), err
}

func (c *CommentNotation) FindPrefixIndex(fields []string) int {
	for i, field := range fields {
		if c.SingleLinePrefixRe != nil && c.SingleLinePrefixRe.MatchString(field) {
			c.Stack.Push(c.SingleLinePrefix)
			return i
		}

		if c.MultiLinePrefixRe != nil && c.MultiLinePrefixRe.MatchString(field) {
			c.Stack.Push(c.MultiLinePrefix)
			return i
		}
	}
	return -1
}

func InitNotationStack() *NotationStack {
	stack := &NotationStack{}
	stack.Items = make([]string, 0)
	stack.Top = -1
	return stack
}

func (s *NotationStack) Push(notation string) {
	s.Items = append(s.Items, notation)
	s.Top++
}

func (s *NotationStack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", errors.New(errStackUnderflow)
	}
	item := s.Items[s.Top]
	s.Items = s.Items[:len(s.Items)-1]
	s.Top--
	return item, nil
}

func (s *NotationStack) Peek() (string, error) {
	if s.IsEmpty() {
		return "", errors.New(errStackUnderflow)
	}
	return s.Items[s.Top], nil
}

func (s *NotationStack) IsEmpty() bool {
	return s.Top == -1
}

func compileAndSetRegexp(exp string) *regexp.Regexp {
	if exp == "" {
		return nil
	}
	return regexp.MustCompile(exp)
}

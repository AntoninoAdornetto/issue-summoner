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
		NewLinePrefix:    `\*`,
	},
	file_ext_python: {
		SingleLinePrefix: `#`,
		MultiLinePrefix:  `(\"\"\")|(\'\'\')`,
		MultiLineSuffix:  `(\"\"\")|(\'\'\')`,
	},
	file_ext_markdown: {
		MultiLinePrefix: `<!--`,
		MultiLineSuffix: `-->`,
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
	cn.MultiLinePrefixRe = compileAndSetRegexp(cn.MultiLinePrefix)
	cn.MultiLineSuffixRe = compileAndSetRegexp(cn.MultiLineSuffix)

	return cn
}


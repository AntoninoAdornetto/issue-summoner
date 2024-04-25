package issue

import (
	"regexp"
)

const (
	SINGLE_LINE_COMMENT = "s"
	MULTI_LINE_COMMENT  = "m"
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
)

type CommentNotation struct {
	SingleLinePrefixRe *regexp.Regexp
	SingleLineSuffixRe *regexp.Regexp
	MultiLinePrefixRe  *regexp.Regexp
	MultiLineSuffixRe  *regexp.Regexp
	NewLinePrefixRe    *regexp.Regexp
}

var CommentNotations = map[string]CommentNotation{
	file_ext_c: {
		SingleLinePrefixRe: regexp.MustCompile(`\/\/`),
		MultiLinePrefixRe:  regexp.MustCompile(`\/\*`),
		MultiLineSuffixRe:  regexp.MustCompile(`\*\/`),
		NewLinePrefixRe:    regexp.MustCompile(`\*`),
	},
	file_ext_python: {
		SingleLinePrefixRe: regexp.MustCompile(`#`),
		MultiLinePrefixRe:  regexp.MustCompile(`['\"]{3}`),
		MultiLineSuffixRe:  regexp.MustCompile(`['\"]{3}`),
	},
	file_ext_markdown: {
		SingleLinePrefixRe: regexp.MustCompile(`<!--`),
		SingleLineSuffixRe: regexp.MustCompile(`-->`),
	},
	"default": {
		SingleLinePrefixRe: regexp.MustCompile(`#`),
	},
}

func NewCommentNotation(ext string) CommentNotation {
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
		return CommentNotations[file_ext_c]
	case file_ext_python:
		return CommentNotations[file_ext_python]
	case file_ext_markdown:
		return CommentNotations[file_ext_markdown]
	default:
		return CommentNotations["default"]
	}
}

func (cn *CommentNotation) FindPrefixIndexAndLineType(line []byte) ([]int, string) {
	if cn.SingleLinePrefixRe != nil {
		if locations := cn.SingleLinePrefixRe.FindSubmatchIndex(line); locations != nil {
			return locations, SINGLE_LINE_COMMENT
		}
	}

	if cn.MultiLinePrefixRe != nil {
		if locations := cn.MultiLinePrefixRe.FindSubmatchIndex(line); locations != nil {
			return locations, MULTI_LINE_COMMENT
		}
	}

	return nil, ""
}

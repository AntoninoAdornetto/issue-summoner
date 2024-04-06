/*
comment.go is responsible for determining the style of comments (single-line, multi-line, or both) that may reside in the files
we are scanning. The file extension of a particular file is used to determine the symbols for both single/multi line comments.

Comments are important because actionable annotations (if they exist) will reside within comment blocks.
*/
package issue

import (
	"strings"
	"unicode"
)

const (
	fileExtAsm        = ".asm"
	fileExtBash       = ".sh"
	fileExtCpp        = ".cpp"
	fileExtC          = ".c"
	fileExtCH         = ".h"
	fileExtCS         = ".cs"
	fileExtGo         = ".go"
	fileExtHaskell    = ".hs"
	fileExtHtml       = ".html"
	fileExtJai        = ".jai"
	fileExtJava       = ".java"
	fileExtJavaScript = ".js"
	fileExtJsx        = ".jsx"
	fileExtKotlin     = ".kt"
	fileExtLisp       = ".lisp"
	fileExtLua        = ".lua"
	fileExtObjC       = ".m"
	fileExtOcaml      = ".ml"
	fileExtPhp        = ".php"
	fileExtPython     = ".py"
	fileExtRuby       = ".rb"
	fileExtRust       = ".rs"
	fileExtMarkdown   = ".md"
	fileExtR          = ".R"
	fileExtScala      = ".scala"
	fileExtSwift      = ".swift"
	fileExtTypeScript = ".ts"
	fileExtTsx        = ".tsx"
	fileExtVim        = ".vim"
	fileExtZig        = ".zig"
)

// Comment contains prefixes that are used to denote
// single-line and multi-line comments.
// Some languages, such as Python, may offer more than one
// way to indicate a multi-line comment.
type Comment struct {
	SingleLinePrefix     []string // Prefixes for single-line comments.
	MultiLineStartPrefix []string // Prefixes for starting a multi-line comment.
	MultiLineEndPrefix   []string // Prefixes for ending a multi-line comment.
}

type CommentStack struct {
	Items []string
}

var CommentSymbols = map[string]Comment{
	fileExtC: {
		SingleLinePrefix:     []string{"//"},
		MultiLineStartPrefix: []string{"/*"},
		MultiLineEndPrefix:   []string{"*/"},
	},
	fileExtPython: {
		SingleLinePrefix:     []string{"#"},
		MultiLineStartPrefix: []string{"\"\"\"", "'''"},
		MultiLineEndPrefix:   []string{"\"\"\"", "'''"},
	},
	fileExtMarkdown: {
		MultiLineStartPrefix: []string{"<!--"},
		MultiLineEndPrefix:   []string{"-->"},
	},
	"default": {
		SingleLinePrefix:     []string{"#"},
		MultiLineStartPrefix: []string{"#"},
		MultiLineEndPrefix:   []string{"#"},
	},
}

func GetCommentSymbols(ext string) Comment {
	switch ext {
	case fileExtC,
		fileExtCpp,
		fileExtJava,
		fileExtJavaScript,
		fileExtJsx,
		fileExtTypeScript,
		fileExtTsx,
		fileExtCS,
		fileExtGo,
		fileExtPhp,
		fileExtSwift,
		fileExtKotlin,
		fileExtRust,
		fileExtObjC,
		fileExtScala:
		return CommentSymbols[fileExtC]
	case fileExtPython:
		return CommentSymbols[fileExtPython]
	case fileExtMarkdown:
		return CommentSymbols[fileExtMarkdown]
	default:
		return CommentSymbols["default"]
	}
}

// @TODO remove ParseCommentContents func. The same logic has been moved to issue.go
func (c Comment) ParseCommentContents(
	line string,
	builder *strings.Builder,
	stack CommentStack,
) (strings.Builder, error) {
	if single, singleSyntax := c.isSingle(line); single {
		// remove single-line comment syntax found from isSingle
		nonCommentLine := strings.SplitAfter(line, singleSyntax)[1]
		builder.WriteString(strings.TrimFunc(nonCommentLine, unicode.IsSpace))
	}
	return *builder, nil
}

// isSingle uses the Comment struct as a receiver to
// determine if the line (from a source code file) is
// a single line comment.
// @TODO remove isSingle func. The same logic has been moved to issue.go
func (c Comment) isSingle(line string) (bool, string) {
	if len(c.SingleLinePrefix) == 0 {
		return false, ""
	}

	for _, s := range c.SingleLinePrefix {
		single := strings.HasPrefix(line, s)
		if single {
			return true, s
		}
	}

	return false, ""
}

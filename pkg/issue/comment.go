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

// Comment contains symbols that are used to denote
// single-line and multi-line comments.
// Some languages, such as python, may offer more than 1
// way to indicate a multi line comment.
// For this reason, a string slice is used.
type Comment struct {
	SingleLineSymbols     []string
	MultiLineStartSymbols []string
	MultiLineEndSymbols   []string
}

type CommentStack struct {
	Items []string
}

var CommentSymbols = map[string]Comment{
	fileExtC: {
		SingleLineSymbols:     []string{"//"},
		MultiLineStartSymbols: []string{"/*"},
		MultiLineEndSymbols:   []string{"*/"},
	},
	fileExtPython: {
		SingleLineSymbols:     []string{"#"},
		MultiLineStartSymbols: []string{"\"\"\"", "'''"},
		MultiLineEndSymbols:   []string{"\"\"\"", "'''"},
	},
	fileExtMarkdown: {
		MultiLineStartSymbols: []string{"<!--"},
		MultiLineEndSymbols:   []string{"-->"},
	},
	"default": {
		SingleLineSymbols:     []string{"#"},
		MultiLineStartSymbols: []string{"#"},
		MultiLineEndSymbols:   []string{"#"},
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
func (c Comment) isSingle(line string) (bool, string) {
	if len(c.SingleLineSymbols) == 0 {
		return false, ""
	}

	for _, s := range c.SingleLineSymbols {
		single := strings.HasPrefix(line, s)
		if single {
			return true, s
		}
	}

	return false, ""
}

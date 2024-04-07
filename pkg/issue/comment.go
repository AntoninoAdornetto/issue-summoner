/*
comment.go is responsible for determining the style of comments (single-line, multi-line, or both) that may reside in the files
we are scanning. The file extension of a particular file is used to determine the symbols for both single/multi line comments.

Comments are important because actionable annotations (if they exist) will reside within comment blocks.
*/
package issue

import "strings"

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

	LINE_TYPE_SRC_CODE    = "c"
	LINE_TYPE_SINGLE      = "single"
	LINE_TYPE_MULTI_START = "multi-start"
	LINE_TYPE_MULTI_END   = "multi-end"
)

// Comment contains prefixes that are used to denote
// single-line and multi-line comments.
// Some languages, such as Python, may offer more than one
// way to indicate a multi-line comment.
type Comment struct {
	SingleLinePrefix     []string // Prefixes for single-line comments.
	MultiLineStartPrefix []string // Prefixes for starting a multi-line comment.
	MultiLineEndPrefix   []string // Prefixes or Suffix for ending a multi-line comment.
	CurrentPrefix        string
	CurrentLineType      string
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

func CommentPrefixes(ext string) Comment {
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

func (c *Comment) ParseLineComment(line string, annotation string) (string, bool) {
	fields := strings.Fields(line)
	evalComment(c, strings.Join(fields, " "))

	if c.CurrentLineType == LINE_TYPE_MULTI_END {
		return trimCommentEnd(fields, annotation, c.CurrentPrefix)
	}

	return trimCommentStart(fields, annotation, c.CurrentPrefix)
}

func trimCommentStart(fields []string, annotation string, prefix string) (string, bool) {
	start := 0

	if len(fields) == 0 {
		return "", false
	}

	for i, s := range fields {
		if s == prefix {
			start = i + 1
		}

		if s == annotation {
			return strings.Join(fields[i+1:], " "), true
		}
	}

	return strings.Join(fields[start:], " "), false
}

func trimCommentEnd(fields []string, annotation string, prefix string) (string, bool) {
	if len(fields) == 0 {
		return "", false
	}

	if fields[len(fields)-1] == prefix {
		fields = fields[:len(fields)-1]
	}

	return trimCommentStart(fields, annotation, prefix)
}

// @TODO associate this method with the Comment struct.
// it will help with unneccesary parsing of source code line types
// that we use in the Scan function of a pending issue
func evalComment(c *Comment, line string) {
	if c.CurrentPrefix != "" {
		if strings.HasPrefix(line, c.CurrentPrefix) || strings.HasSuffix(line, c.CurrentPrefix) {
			return
		}
	}

	for _, s := range c.SingleLinePrefix {
		if strings.HasPrefix(line, s) {
			c.CurrentLineType = LINE_TYPE_SINGLE
			c.CurrentPrefix = s
		}
	}

	for i := range c.MultiLineStartPrefix {
		isMultiStart := strings.HasPrefix(line, c.MultiLineStartPrefix[i])
		isMultiEnd := strings.HasSuffix(line, c.MultiLineEndPrefix[i])

		if isMultiStart {
			c.CurrentLineType = LINE_TYPE_MULTI_START
			c.CurrentPrefix = c.MultiLineStartPrefix[i]
		}

		if isMultiEnd {
			c.CurrentLineType = LINE_TYPE_MULTI_END
			c.CurrentPrefix = c.MultiLineEndPrefix[i]
		}
	}
}

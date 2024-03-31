/*
comment.go is responsible for determining the style of comments (single-line, multi-line, or both) that may reside in the files
we are scanning. The file extension of a particular file is used to determine the symbols for both single/multi line comments.

Comments are important because actionable annotations (if they exist) will reside within comment blocks.
*/
package issue

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

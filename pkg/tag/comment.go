/*
This file contains the comment syntax for various programming languages and file extensions.
The constants in this file are used to determine the comment syntax for a given file when parsing source code for tag annotations.
This will help determine if the tag annotation is within a single line comment or a multi-line comment.

@TODO Increase programming language support. Python, Haskel, etc.
*/
package tag

type CommentLangSyntax struct {
	SingleLineCommentSymbols string
	MultiLineCommentSymbols  MultiLineCommentSyntax
}

type MultiLineCommentSyntax struct {
	CommentStartSymbol string
	CommentEndSymbol   string
}

var CommentSyntaxMap = map[string]CommentLangSyntax{
	"c-derived": {
		SingleLineCommentSymbols: CommentCSingle,
		MultiLineCommentSymbols: MultiLineCommentSyntax{
			CommentStartSymbol: CommentCMultiStart,
			CommentEndSymbol:   CommentCMultiEnd,
		},
	},
	"default": {
		SingleLineCommentSymbols: "#",
		MultiLineCommentSymbols: MultiLineCommentSyntax{
			CommentStartSymbol: "#",
			CommentEndSymbol:   "#",
		},
	},
}

func CommentSyntax(fileExtension string) CommentLangSyntax {
	switch fileExtension {
	case FileExtC,
		FileExtCpp,
		FileExtJava,
		FileExtJavaScript,
		FileExtJsx,
		FileExtTypeScript,
		FileExtTsx,
		FileExtCS,
		FileExtGo,
		FileExtPhp,
		FileExtSwift,
		FileExtKotlin,
		FileExtRust,
		FileExtObjC,
		FileExtScala:
		return CommentSyntaxMap["c-derived"]
	default:
		return CommentSyntaxMap["default"]
	}
}

const (
	FileExtAsm        = ".asm"
	FileExtBash       = ".sh"
	FileExtCpp        = ".cpp"
	FileExtC          = ".c"
	FileExtCH         = ".h"
	FileExtCS         = ".cs"
	FileExtGo         = ".go"
	FileExtHaskell    = ".hs"
	FileExtHtml       = ".html"
	FileExtJai        = ".jai"
	FileExtJava       = ".java"
	FileExtJavaScript = ".js"
	FileExtJsx        = ".jsx"
	FileExtKotlin     = ".kt"
	FileExtLisp       = ".lisp"
	FileExtLua        = ".lua"
	FileExtObjC       = ".m"
	FileExtOcaml      = ".ml"
	FileExtPhp        = ".php"
	FileExtPython     = ".py"
	FileExtRuby       = ".rb"
	FileExtRust       = ".rs"
	FileExtR          = ".R"
	FileExtScala      = ".scala"
	FileExtSwift      = ".swift"
	FileExtTypeScript = ".ts"
	FileExtTsx        = ".tsx"
	FileExtVim        = ".vim"
	FileExtZig        = ".zig"
)

// Single Line comments
const (
	CommentCSingle = "//"
)

// Multi-Line comments (start)
const (
	CommentCMultiStart = "/*"
)

// Multi-Line comments (end)
const (
	CommentCMultiEnd = "*/"
)

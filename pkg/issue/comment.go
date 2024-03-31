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

type SingleLineComment = string

type MultiLineComment struct {
	StartSymbols string
	EndSymbols   string
}

type CommentManager struct {
	SingleLineSymbols SingleLineComment
	MultiLineSymbols  MultiLineComment
// Comment contains symbols that are used to denote
// single-line and multi-line comments.
// Some languages, such as python, may offer more than 1
// way to indicate a multi line comment.
// For this reason, a string slice is used.
}

var CommentSyntaxMap = map[string]CommentManager{
	FileExtC: {
		SingleLineSymbols: SingleLineC,
		MultiLineSymbols: MultiLineComment{
			StartSymbols: MultiLineStartC,
			EndSymbols:   MultiLineEndC,
		},
	},
	FileExtMarkdown: {
		SingleLineSymbols: "",
		MultiLineSymbols: MultiLineComment{
			StartSymbols: MultiLineStartMd,
			EndSymbols:   MultiLineEndMd,
		},
	},
}

// var CommentSyntaxMap = map[string]CommentLangSyntax{
// 	"c-derived": {
// 		SingleLineCommentSymbols: CommentCSingle,
// 		MultiLineCommentSymbols: MultiLineCommentSyntax{
// 			CommentStartSymbol: CommentCMultiStart,
// 			CommentEndSymbol:   CommentCMultiEnd,
// 		},
// 	},
// 	"default": {
// 		SingleLineCommentSymbols: "#",
// 		MultiLineCommentSymbols: MultiLineCommentSyntax{
// 			CommentStartSymbol: "#",
// 			CommentEndSymbol:   "#",
// 		},
// 	},
// }
//
// func CommentSyntax(fileExtension string) CommentLangSyntax {
// 	switch fileExtension {
// 	case FileExtC,
// 		FileExtCpp,
// 		FileExtJava,
// 		FileExtJavaScript,
// 		FileExtJsx,
// 		FileExtTypeScript,
// 		FileExtTsx,
// 		FileExtCS,
// 		FileExtGo,
// 		FileExtPhp,
// 		FileExtSwift,
// 		FileExtKotlin,
// 		FileExtRust,
// 		FileExtObjC,
// 		FileExtScala:
// 		return CommentSyntaxMap["c-derived"]
// 	default:
// 		return CommentSyntaxMap["default"]
// 	}
// }
//
//
// // Single Line comments
// const (
// 	CommentCSingle = "//"
// )
//
// // Multi-Line comments (start)
// const (
// 	CommentCMultiStart = "/*"
// )
//
// // Multi-Line comments (end)
// const (
// 	CommentCMultiEnd = "*/"
// )

package issue

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
	FileExtMarkdown   = ".md"
	FileExtR          = ".R"
	FileExtScala      = ".scala"
	FileExtSwift      = ".swift"
	FileExtTypeScript = ".ts"
	FileExtTsx        = ".tsx"
	FileExtVim        = ".vim"
	FileExtZig        = ".zig"

	SingleLineC     = "//"
	MultiLineStartC = "/**"
	MultiLineEndC   = "/**"

	MultiLineStartMd = "<!--"
	MultiLineEndMd   = "-->"
)

type SingleLineComment = string

type MultiLineComment struct {
	StartSymbols string
	EndSymbols   string
}

type CommentManager struct {
	SingleLineSymbols SingleLineComment
	MultiLineSymbols  MultiLineComment
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

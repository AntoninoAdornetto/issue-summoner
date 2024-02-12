/*
This file contains the comment syntax for various programming languages and file extensions.
The constants in this file are used alongside functions in the tag package to determine comment sytax and assist in parsing
information about a `Tag` from a given file.
*/
package tag

// @TODO - Create a map to allow for easy lookup of comment syntax based on file extension.
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
	FileExtVim        = ".vim"
	FileExtZig        = ".zig"
)

// Single Line comments
const (
	CommentAsmSingle        = ";"
	CommentBashSingle       = "#"
	CommentCppSingle        = "//"
	CommentCSingle          = "//"
	CommentCSSingle         = "//"
	CommentGoSingle         = "//"
	CommentHaskellSingle    = "--"
	CommentJavaSingle       = "//"
	CommentJaiSingle        = "//"
	CommentJavaScriptSingle = "//"
	CommentKotlinSingle     = "//"
	CommentLispSingle       = ";"
	CommentLuaSingle        = "--"
	CommentObjCSingle       = "//"
	CommentOcamlSingle      = "" // OCaml does not have a designated single-line comment syntax.
	CommentPHPSingle        = "//"
	CommentPythonSingle     = "#"
	CommentRubySingle       = "#"
	CommentRustSingle       = "//"
	CommentRSingle          = "#"
	CommentScalaSingle      = "//"
	CommentSwiftSingle      = "//"
	CommentTypeScriptSingle = "//"
	CommentVimSingle        = "\""
	CommentZigSingle        = "//"
)

// Multi-Line comments (start)
const (
	CommentCppMultiStart        = "/*"
	CommentCMultiStart          = "/*"
	CommentCSMultiStart         = "/*"
	CommentGoMultiStart         = "/*"
	CommentHaskellMultiStart    = "{-"
	CommentHtmlMultiStart       = "<!--"
	CommentJavaMultiStart       = "/*"
	CommentJavaScriptMultiStart = "/*"
	CommentKotlinMultiStart     = "/*"
	CommentLispMultiStart       = "#|"
	CommentLuaMultiStart        = "--[["
	CommentObjCMultiStart       = "/*"
	CommentOcamlMultiStart      = "(*"
	CommentPhpMultiStart        = "/*"
	CommentPythonMultiStart     = "'''"
	CommentRubyMultiStart       = "=begin"
	CommentRustMultiStart       = "/*"
	CommentScalaMultiStart      = "/*"
	CommentSwiftMultiStart      = "/*"
	CommentTypeScriptMultiStart = "/*"
	CommentZigMultiStart        = "/*"
)

// Multi-Line comments (end)
const (
	CommentCppMultiEnd        = "*/"
	CommentCMultiEnd          = "*/"
	CommentCSMultiEnd         = "*/"
	CommentGoMultiEnd         = "*/"
	CommentHaskellMultiEnd    = "-}"
	CommentHtmlMultiEnd       = "-->"
	CommentJavaMultiEnd       = "*/"
	CommentJavaScriptMultiEnd = "*/"
	CommentKotlinMultiEnd     = "*/"
	CommentLispMultiEnd       = "|#"
	CommentLuaMultiEnd        = "--]]"
	CommentObjCMultiEnd       = "*/"
	CommentOcamlMultiEnd      = "*)"
	CommentPhpMultiEnd        = "*/"
	CommentPythonMultiEnd     = "'''"
	CommentRubyMultiEnd       = "=end"
	CommentRustMultiEnd       = "*/"
	CommentScalaMultiEnd      = "*/"
	CommentSwiftMultiEnd      = "*/"
	CommentTypeScriptMultiEnd = "*/"
	CommentZigMultiEnd        = "*/"
)

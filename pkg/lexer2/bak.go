package lexer2

//
// import (
// 	"fmt"
// 	"path/filepath"
// )
//
// const (
// 	byte_null = 0
// )
//
// // Lexer acts as the base lexer for other lexers that are specific
// // to different programming languages. Each programming language may
// // have a different way to denote comment syntax for single and multi
// // line comments. Each lexer that we satisfy, depending on the programming
// // language can utilize the various receiver methods from the base lexer
// // to assist with scanning through the src code to parse comments and annotations
// type Lexer struct {
// 	Src      []byte
// 	FileName string
// 	Tokens   []Token
// 	Start    int
// 	Current  int
// 	Line     int
// }
//
// // @DOCS_TODO update LexicalAnalyzer comment, it doesn't make much sense right now
// // LexicalAnalyzer - each programming language that we build comment/annotation
// // Analyzers for must satisfy the methods contained in this interface. The end
// // result of properly implementing these methods will allow the program to support
// // a wide variety of programming lanaguages and the comment syntax those languages
// // require
// type LexicalAnalyzer interface {
// 	AnaylzeToken() error
// 	String(delim byte) error
// 	Comment() error
// }
//
// func NewBaseLexer(src []byte, fileName string) *Lexer {
// 	return &Lexer{
// 		Src:      src,
// 		FileName: fileName,
// 		Tokens:   make([]Token, 0),
// 		Start:    0,
// 		Current:  0,
// 		Line:     1,
// 	}
// }
//
// func NewLexicalAnalyzer(base *Lexer) (LexicalAnalyzer, error) {
// 	ext := filepath.Ext(base.FileName)
// 	switch {
// 	case ext == ".c":
// 		return &CLexer{BaseLexer: base}, nil
// 	default:
// 		return nil, fmt.Errorf(
// 			"unsupported file type of %s. please open a feature request if you would like support.",
// 			ext,
// 		)
// 	}
// }
//
// func (base *Lexer) AnalyzeTokens(target LexicalAnalyzer) ([]Token, error) {
// 	for range base.Src {
// 		base.Start = base.Current
// 		if err := target.AnaylzeToken(); err != nil {
// 			return nil, err
// 		}
//
// 		if base.next() == byte_null {
// 			break
// 		}
// 	}
//
// 	base.Tokens = append(base.Tokens, NewToken(TOKEN_EOF, base))
// 	return base.Tokens, nil
// }
//
// func (base *Lexer) next() byte {
// 	if base.end() {
// 		return byte_null
// 	}
// 	base.Current++
// 	return base.Src[base.Current]
// }
//
// func (base *Lexer) peek() byte {
// 	if base.end() {
// 		return byte_null
// 	}
// 	return base.Src[base.Current]
// }
//
// func (base *Lexer) peekNext() byte {
// 	if base.Current+1 >= len(base.Src) {
// 		return byte_null
// 	}
// 	return base.Src[base.Current+1]
// }
//
// func (base *Lexer) end() bool {
// 	return base.Current >= len(base.Src)-1
// }

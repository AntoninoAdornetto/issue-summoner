/*
Copyright © 2024 AntoninoAdornetto

The c.go file is responsible for satisfying the `LexicalTokenizer` interface in the `lexer.go` file.
The methods are a strict set of rules for handling single & multi line comments for c-like languages.
The result, if an issue annotation is located, is a slice of tokens that will provide information about
the action item contained in the comment. If a comment does not contain an issue annotation, all subsequent
tokens of the remaining comment bytes will be ignored and removed from the `DraftTokens` slice.
*/
package lexer

import (
	"fmt"
	"regexp"
)

var allowed = []string{
	".c",
	".h",
	".cpp",
	".java",
	".js",
	".jsx",
	".ts",
	".tsx",
	".cs",
	".go",
	".php",
	".swift",
	".kt",
	".rs",
	".m",
	".scala",
}

type CLexer struct{}

func (cl *CLexer) AnalyzeToken(lex *Lexer) error {
	b := lex.peek()
	switch b {
	case FORWARD_SLASH:
		return cl.Comment(lex)
	case DOUBLE_QUOTE:
		return cl.String(lex, DOUBLE_QUOTE)
	case QUOTE:
		return cl.String(lex, QUOTE)
	case BACK_TICK:
		return cl.String(lex, BACK_TICK)
	case NEWLINE:
		lex.Line++
		return nil
	default:
		return nil
	}
}

func (cl *CLexer) Comment(lex *Lexer) error {
	switch lex.peekNext() {
	case FORWARD_SLASH:
		return cl.SingleLineComment(lex)
	case ASTERISK:
		return cl.MultiLineComment(lex)
	default:
		return nil
	}
}

func (cl *CLexer) SingleLineComment(lex *Lexer) error {
	for !lex.isEnd() && lex.peekNext() != NEWLINE {
		lex.next()
	}
	comment := lex.Source[lex.Start : lex.Current+1]
	lex.addToken(SINGLE_LINE_COMMENT, comment)
	return nil
}

func (cl *CLexer) MultiLineComment(lex *Lexer) error {
	for !lex.isEnd() {
		b := lex.next()
		if b == NEWLINE {
			lex.Line++
		}

		if b == ASTERISK && lex.peekNext() == FORWARD_SLASH {
			lex.next()
			break
		}
	}

	if lex.isEnd() {
		src := lex.Source[lex.Start:lex.Current]
		return lex.report(fmt.Sprintf("could not locate closing multi line comment: %s", src))
	}

	comment := lex.Source[lex.Start : lex.Current+1]
	lex.addToken(MULTI_LINE_COMMENT, comment)
	return nil
}

func (cl *CLexer) String(lex *Lexer, delim byte) error {
	for !lex.isEnd() && lex.peekNext() != delim {
		b := lex.next()
		if b == NEWLINE {
			lex.Line++
		}
	}
	lex.next() // closing delimiter
	return nil
}

func (cl *CLexer) ParseCommentTokens(lex *Lexer, annotation []byte) ([]Comment, error) {
	comments := make([]Comment, 0)
	for i, token := range lex.Tokens {
		switch token.TokenType {
		case SINGLE_LINE_COMMENT:
			comment := token.ParseSingleLineCommentToken(annotation, trimCommentC)
			comment.Push(&comments, lex.FileName, i)
		case MULTI_LINE_COMMENT:
			comment := token.ParseMultiLineCommentToken(annotation, trimCommentC)
			comment.Push(&comments, lex.FileName, i)
		default:
			continue
		}
	}
	return comments, nil
}

func findAnnotationLocations(annotation []byte, commentText []byte) []int {
	re := regexp.MustCompile(string(annotation))
	return re.FindIndex(commentText)
}

func trimCommentC(r rune) bool {
	switch r {
	case rune(WHITESPACE), rune(TAB), rune(NEWLINE), rune(ASTERISK), rune(FORWARD_SLASH):
		return true
	default:
		return false
	}
}

func IsAdoptedFromC(ext string) bool {
	for _, lang := range allowed {
		if ext == lang {
			return true
		}
	}
	return false
}

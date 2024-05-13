package lexer

import (
	"bytes"
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
			comment := ParseSingleLineCommentToken(&token, annotation)
			pushComment(&comment, &comments, lex.FileName, i)
		case MULTI_LINE_COMMENT:
			comment := ParseMultiLineCommentToken(&token, annotation)
			pushComment(&comment, &comments, lex.FileName, i)
		default:
			continue
		}
	}
	return comments, nil
}

func ParseSingleLineCommentToken(token *Token, annotation []byte) Comment {
	loc := findAnnotationLocations(annotation, token.Lexeme)
	if loc == nil {
		return Comment{}
	}
	end := loc[1]
	title := bytes.TrimFunc(token.Lexeme[end:], trimComment)
	return Comment{
		Title:  title,
		Source: token.Lexeme,
	}
}

func ParseMultiLineCommentToken(token *Token, annotation []byte) Comment {
	loc := findAnnotationLocations(annotation, token.Lexeme)
	if loc == nil {
		return Comment{}
	}
	end := loc[1]
	content := bytes.TrimFunc(token.Lexeme[end:], trimComment)
	newLines := bytes.Split(content, []byte("\n"))

	comment := Comment{
		Title:  bytes.TrimFunc(newLines[0], trimComment),
		Source: token.Lexeme,
	}

	for i := 1; i < len(newLines); i++ {
		comment.Description = append(
			comment.Description,
			bytes.TrimLeftFunc(newLines[i], trimComment)...)

		if i != len(newLines)-1 {
			comment.Description = append(comment.Description, ' ')
		}
	}

	return comment
}

func pushComment(comment *Comment, comments *[]Comment, fileName string, index int) {
	if comment.Validate() {
		comment.Prepare(fileName, index)
		*comments = append(*comments, *comment)
	}
}

func findAnnotationLocations(annotation []byte, commentText []byte) []int {
	re := regexp.MustCompile(string(annotation))
	return re.FindIndex(commentText)
}

func trimComment(r rune) bool {
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

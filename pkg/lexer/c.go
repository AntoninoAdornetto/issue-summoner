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
	case '/':
		return cl.Comment(lex)
	case '"':
		return cl.String(lex, '"')
	case '\'':
		return cl.String(lex, '\'')
	case '`':
		return cl.String(lex, '`')
	case '\n':
		lex.Line++
		return nil
	default:
		return nil
	}
}

func (cl *CLexer) Comment(lex *Lexer) error {
	switch lex.peekNext() {
	case '/':
		return cl.SingleLineComment(lex)
	case '*':
		return cl.MultiLineComment(lex)
	default:
		return nil
	}
}

func (cl *CLexer) SingleLineComment(lex *Lexer) error {
	for !lex.isEnd() && lex.peekNext() != '\n' {
		lex.next()
	}
	comment := lex.Source[lex.Start : lex.Current+1]
	lex.addToken(SINGLE_LINE_COMMENT, comment)
	return nil
}

func (cl *CLexer) MultiLineComment(lex *Lexer) error {
	for !lex.isEnd() {
		b := lex.next()
		if b == '\n' {
			lex.Line++
		}

		if b == '*' && lex.peekNext() == '/' {
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
		if b == '\n' {
			lex.Line++
		}
	}
	_ = lex.next() // closing delimiter
	return nil
}

func (cl *CLexer) ParseCommentTokens(lex *Lexer, annotation []byte) ([]Comment, error) {
	comments := make([]Comment, 0)
	for i, token := range lex.Tokens {
		switch token.TokenType {
		case SINGLE_LINE_COMMENT:
			comment := token.ParseSingleLineCommentToken(annotation, trimC)
			pushComment(&comment, &comments, lex.FileName, i)
		case MULTI_LINE_COMMENT:
			comment := token.ParseMultiLineCommentToken(annotation, trimC)
			pushComment(&comment, &comments, lex.FileName, i)
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

// specific to c like languages
func trimC(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '*', '/':
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

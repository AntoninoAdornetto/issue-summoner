package lexer

import "fmt"
var allowed = []string{
	".c",
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
		lex.next()
		if lex.peek() == '\n' {
			lex.Line++
		}
	}

	if lex.isEnd() {
		src := lex.Source[lex.Start : lex.Current+1]
		return lex.report(fmt.Sprintf("unterminated string: %s", src))
	}

	_ = lex.next() // closing delimiter
	str := lex.Source[lex.Start+1 : lex.Current]
	lex.addToken(STRING, str)
	return nil
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

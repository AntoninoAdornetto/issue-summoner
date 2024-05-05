package lexer

import "fmt"

type CLexer struct{}

func (cl *CLexer) AnalyzeToken(lex *Lexer) error {
	b := lex.peek()
	switch b {
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

func IsAdoptedFromC(ext string) bool {
	for _, lang := range allowed {
		if ext == lang {
			return true
		}
	}
	return false
}

package lexer

import "fmt"

type PythonLexer struct{}

func (pl *PythonLexer) AnalyzeToken(lex *Lexer) error {
	b := lex.peek()
	switch b {
	case '#', '\'', '"':
		return pl.Comment(lex)
	case '\n':
		lex.Line++
		return nil
	default:
		return nil
	}
}

// @TODO implement python string lexer
func (pl *PythonLexer) String(lex *Lexer, delim byte) error {
	return nil
}

func (pl *PythonLexer) Comment(lex *Lexer) error {
	p := lex.peekNext()
	if !pl.isQuoted(p) {
		return pl.SingleLineComment(lex)
	}

	lex.next()
	p = lex.peekNext()
	if pl.isQuoted(p) {
		return pl.MultiLineComment(lex)
	}
	return nil
}

func (pl *PythonLexer) SingleLineComment(lex *Lexer) error {
	for !lex.isEnd() && lex.peekNext() != '\n' {
		lex.next()
	}
	comment := lex.Source[lex.Start : lex.Current+1]
	lex.addToken(SINGLE_LINE_COMMENT, comment)
	return nil
}

// @TODO fix broken multi line parsing. handling strings should resolve the bug
func (pl *PythonLexer) MultiLineComment(lex *Lexer) error {
	for !lex.isEnd() {
		b := lex.next()
		if b == '\n' {
			lex.Line++
		}

		// is current and next byte a quote(s)
		if pl.isQuoted(b) && pl.isQuoted(lex.peekNext()) {
			// move to next byte and check for the 3rd quote(s)
			lex.next()
			if pl.isQuoted(lex.peekNext()) {
				break
			}
		}
	}

	if lex.isEnd() {
		src := lex.Source[lex.Start:lex.Current]
		return lex.report(fmt.Sprintf("could not locate closing multi line comment: %s", src))
	}

	lex.next() // proceed to final quote
	comment := lex.Source[lex.Start : lex.Current+1]
	lex.addToken(MULTI_LINE_COMMENT, comment)
	return nil
}

func (pl *PythonLexer) isQuoted(b byte) bool {
	return b == '\'' || b == '"'
}

func (pl *PythonLexer) ParseCommentTokens(lex *Lexer, annotation []byte) ([]Comment, error) {
	comments := make([]Comment, 0)
	for i, token := range lex.Tokens {
		switch token.TokenType {
		case SINGLE_LINE_COMMENT:
			comment := token.ParseSingleLineCommentToken(annotation, trimPython)
			pushComment(&comment, &comments, lex.FileName, i)
		case MULTI_LINE_COMMENT:
			comment := token.ParseMultiLineCommentToken(annotation, trimPython)
			pushComment(&comment, &comments, lex.FileName, i)
		default:
			continue
		}
	}
	return comments, nil
}

func trimPython(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\'', '"', '#':
		return true
	default:
		return false
	}
}

package lexer

import "fmt"

type PyLexer struct {
	Base        *Lexer
	DraftTokens []Token
	annotated   bool
	line        int
}

func (py *PyLexer) AnalyzeToken() error {
	currentByte := py.Base.peek()
	switch currentByte {
	case QUOTE, DOUBLE_QUOTE, BACK_TICK:
		return py.String(currentByte)
	case HASH:
		return py.Comment()
	case NEWLINE:
		py.Base.Line++
		return nil
	default:
		return nil
	}
}

func (py *PyLexer) String(delim byte) error {
	return nil
}

// @TODO process multi line python comment tokens
func (py *PyLexer) Comment() error {
	currentByte := py.Base.peek()
	if currentByte == HASH {
		return py.singleLineComment()
	}

	return nil
}

func (py *PyLexer) singleLineComment() error {
	if err := py.Base.initTokenization(TOKEN_SINGLE_LINE_COMMENT_START, &py.DraftTokens); err != nil {
		return err
	}

	py.Base.next()
	for !py.Base.pastEnd() {
		lexeme := py.Base.nextLexeme()
		if err := py.processLexeme(lexeme, TOKEN_SINGLE_LINE_COMMENT); err != nil {
			return err
		}

		if next := py.Base.peekNext(); next == NEWLINE || next == 0 {
			next = py.Base.next()
			if next == NEWLINE {
				py.Base.Line++
			}

			py.Base.resetStartIndex()
			closeToken := NewToken(TOKEN_SINGLE_LINE_COMMENT_END, []byte{next}, py.Base)
			py.DraftTokens = append(py.DraftTokens, closeToken)
			break
		}

		py.Base.next()
	}

	if py.annotated {
		py.Base.promoteTokens(py.DraftTokens)
	}

	py.reset()
	return nil
}

func (py *PyLexer) processLexeme(lexeme []byte, commentType TokenType) error {
	if len(lexeme) == 0 {
		return nil
	}

	tokens, err := py.Base.processAnnotation(lexeme, py.annotated)
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		py.DraftTokens = append(py.DraftTokens, tokens...)
		py.annotated = true
		py.line = py.Base.Line
		return nil
	}

	switch commentType {
	case TOKEN_SINGLE_LINE_COMMENT:
		token := NewToken(TOKEN_COMMENT_TITLE, lexeme, py.Base)
		py.DraftTokens = append(py.DraftTokens, token)
		return nil
		// @TODO add multi line comment type check for appending ML tokens
	default:
		return fmt.Errorf(errTargetTokenize, string(lexeme), decodeTokenType(commentType))
	}
}

func (py *PyLexer) reset() {
	py.annotated = false
	py.DraftTokens = py.DraftTokens[:0]
	py.line = 0
}

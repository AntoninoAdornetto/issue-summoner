package lexer

import "fmt"

type ShellLexer struct {
	Base        *Lexer  // holds shared byte consumption methods
	DraftTokens []Token // Unvalidated tokens
	annotated   bool    // Issue annotation indicator
}

func (sh *ShellLexer) AnalyzeToken() error {
	currentByte := sh.Base.peek()
	switch currentByte {
	case QUOTE, DOUBLE_QUOTE, BACK_TICK:
		return sh.String(currentByte)
	case HASH:
		return sh.Comment()
	case NEWLINE:
		sh.Base.Line++
		return nil
	default:
		return nil
	}
}

func (sh *ShellLexer) String(delim byte) error {
	for !sh.Base.pastEnd() {
		if sh.Base.next() == NEWLINE {
			sh.Base.Line++
		}

		if sh.Base.peekNext() == delim {
			sh.Base.next()
			break
		}
	}

	if sh.Base.pastEnd() {
		return fmt.Errorf(errStringClose, delim, sh.Base.Src[sh.Base.Start:])
	}

	return nil
}

func (sh *ShellLexer) Comment() error {
	if err := sh.Base.initTokenization(TOKEN_SINGLE_LINE_COMMENT_START, &sh.DraftTokens); err != nil {
		return err
	}

	sh.Base.next()
	for !sh.Base.pastEnd() {
		lexeme := sh.Base.nextLexeme()
		if err := sh.processLexeme(lexeme, TOKEN_SINGLE_LINE_COMMENT); err != nil {
			return err
		}

		if next := sh.Base.peekNext(); next == NEWLINE || next == 0 {
			next = sh.Base.next()
			if next == NEWLINE {
				sh.Base.Line++
			}

			sh.Base.resetStartIndex()
			closeToken := NewToken(TOKEN_SINGLE_LINE_COMMENT_END, []byte{next}, sh.Base)
			sh.DraftTokens = append(sh.DraftTokens, closeToken)
			break
		}

		sh.Base.next()
	}

	if sh.annotated {
		sh.Base.promoteTokens(sh.DraftTokens)
	}

	sh.reset()
	return nil
}

func (sh *ShellLexer) processLexeme(lexeme []byte, commentType TokenType) error {
	if len(lexeme) == 0 {
		return nil
	}

	if commentType != TOKEN_SINGLE_LINE_COMMENT {
		return fmt.Errorf(errTargetTokenizeSl, string(lexeme), decodeTokenType(commentType))
	}

	tokens, err := sh.Base.processAnnotation(lexeme, sh.annotated)
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		sh.DraftTokens = append(sh.DraftTokens, tokens...)
		sh.annotated = true
		return nil
	}

	token := NewToken(TOKEN_COMMENT_TITLE, lexeme, sh.Base)
	sh.DraftTokens = append(sh.DraftTokens, token)
	return nil
}

func (sh *ShellLexer) reset() {
	sh.annotated = false
	sh.DraftTokens = sh.DraftTokens[:0]
}

func isShell(ext string) bool {
	switch ext {
	case ".bash",
		".sh",
		".zsh",
		".ps1":
		return true
	default:
		return false
	}
}

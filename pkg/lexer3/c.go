package lexer3

import (
	"bytes"
	"fmt"
	"unicode"
)

// CLexer - is designed to handle comment tokens for not only
// the c programming language but for any programming language that
// denotes single/multi line comments using the same syntax as C. e.g, ("//") ("/*" "*/")
// The base lexer assists the CLexer in consuming bytes to construct tokens as the src
// code files are scanned. Each token that is validated and contains an annotation is appended
// to the Base Tokens slice and not the Tokens slice contained in CLexer. The Tokens slice in CLexer
// is meant to be reused to prevent multiple memory allocations.
type CLexer struct {
	Base   *Lexer
	Tokens []Token
}

func (target *CLexer) AnaylzeToken() error {
	b := target.Base.peek()
	switch b {
	case QUOTE, DOUBLE_QUOTE, BACK_TICK:
		return target.String(b)
	case FORWARD_SLASH:
		return target.Comment()
	case NEWLINE:
		target.Base.Line++
		return nil
	}
	return nil
}

// String - is used to prevent one specific edge case where a src code
// string may contain comment notation. For example, there could be a string
// that contains two forward slashes or a forward slash followed by an asterisk.
// "//" or "/*" or "*/". This method is used to prevent that edge case from happening
// iterating through the string till the delim is reached. String tokens are not persisted.
func (target *CLexer) String(delim byte) error {
	for target.Base.peekNext() != delim {
		b := target.Base.next()
		switch b {
		case NEWLINE:
			target.Base.Line++
		case byte_null:
			state := target.Base.Src[target.Base.Start:target.Base.Current]
			msg := fmt.Sprintf("Failed to find closing string delimiter: %s", state)
			return target.Base.report(msg)
		default:
			continue
		}
	}

	target.Base.next()
	return nil
}

// Comment - once we have consumed a forward slash byte, we will peek into the next
// byte and check if it is another forward slash (single line comment) or an asterisk
// (multi line comment). Peeking the next byte does not consume it.
func (target *CLexer) Comment() error {
	switch target.Base.peekNext() {
	case FORWARD_SLASH:
		return target.SingleLineComment()
	case ASTERISK:
		return target.MultiLineComment()
	default:
		return nil
	}
}

func (target *CLexer) prepareSingleLineComment() {
	target.Base.next() // consume the second forward slash
	token := NewToken(TOKEN_SINGLE_LINE_COMMENT_START, target.Base)
	target.Tokens = append(target.Tokens, token)
}

func (target *CLexer) SingleLineComment() error {
	base := target.Base
	offset := 0
	isAnnotated := false

	target.prepareSingleLineComment()
	for !base.end() {
		current := base.next()

		if !unicode.IsSpace(rune(current)) {
			offset++
			continue
		}

		if err := target.makeSingleLineCommentToken(offset, &isAnnotated); err != nil {
			return err
		}

		offset = 0

		if current == NEWLINE {
			break
		}
	}

	if isAnnotated {
		target.appendClosingCommentToken(false)
		target.Base.Tokens = append(target.Base.Tokens, target.Tokens...)
	}

	target.clearTokens()
	return nil
}

func (target *CLexer) MultiLineComment() error {
	// base := target.Base
	return nil
}

func (target *CLexer) makeSingleLineCommentToken(offset int, isAnnotated *bool) error {
	lexeme, err := target.Base.extractRange(offset)
	if err != nil {
		return err
	}

	switch {
	case len(lexeme) == 0:
		return nil
	case !*isAnnotated && bytes.Equal(lexeme, target.Base.Annotation):
		*isAnnotated = true
		token := target.Base.makeToken(TOKEN_COMMENT_ANNOTATION, lexeme)
		target.Tokens = append(target.Tokens, token)
		return nil
	default:
		token := target.Base.makeToken(TOKEN_COMMENT_TITLE, lexeme)
		target.Tokens = append(target.Tokens, token)
		return nil
	}
}

// @TODO maketoken causing bug
func (target *CLexer) appendClosingCommentToken(isMulti bool) {
	if isMulti {
		token := target.Base.makeToken(
			TOKEN_MULTI_LINE_COMMENT_END,
			[]byte{ASTERISK, FORWARD_SLASH},
		)
		target.Tokens = append(target.Tokens, token)
	} else {
		token := target.Base.makeToken(TOKEN_SINGLE_LINE_COMMENT_END, []byte{NEWLINE})
		target.Tokens = append(target.Tokens, token)
	}
}

func (target *CLexer) clearTokens() {
	target.Tokens = target.Tokens[:0]
}

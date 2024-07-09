package lexer2

import (
	"bytes"
	"fmt"
	"unicode"
)

type CLexer struct {
	Base   *Lexer  // Holds our comment tokens & methods for consuming bytes
	Tokens []Token // Used to reduce memory allocations. Validated tokens are stored in Base Lexer
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

func (target *CLexer) SingleLineComment() error {
	offset := 0
	isAnnotated := false
	base := target.Base
	tokens := target.Tokens
	base.next() // consume the 2nd forward slash and start single line comment lexing

	tokens = append(
		tokens,
		base.makeToken(TOKEN_SINGLE_LINE_COMMENT_START, []byte{FORWARD_SLASH, FORWARD_SLASH}),
	)

	for target.Base.peekNext() != NEWLINE {
		current := base.next()

		if !unicode.IsSpace(rune(current)) {
			offset++
			continue
		}

		lexeme, err := base.extractRange(offset)
		if err != nil {
			return base.report(err.Error())
		}

		offset = 0

		if len(lexeme) == 0 {
			continue
		}

		if !isAnnotated && bytes.Equal(lexeme, base.Annotation) {
			tokens = append(tokens, base.makeToken(TOKEN_COMMENT_ANNOTATION, lexeme))
			isAnnotated = true
		} else {
			tokens = append(tokens, base.makeToken(TOKEN_COMMENT_TITLE, lexeme))
		}
	}
	// tokens := make([]Token, 0, 10)
	// tokens = append(tokens, target.Base.makeToken(TOKEN_SINGLE_LINE_COMMENT_START, []byte("//")))
	//
	// offset, isAnnotated := 0, false
	// for !target.Base.end() {
	// 	b := target.Base.next()
	//
	// 	if !unicode.IsSpace(rune(b)) {
	// 		offset++
	// 		continue
	// 	}
	//
	// 	lexeme, err := target.Base.extractRange(offset)
	// 	if err != nil {
	// 		return target.Base.report(err.Error())
	// 	}
	//
	// 	if !isAnnotated && bytes.Equal(lexeme, target.Base.Annotation) {
	// 		isAnnotated = true
	// 		tokens = append(tokens, target.Base.makeToken(TOKEN_COMMENT_ANNOTATION, lexeme))
	// 	} else {
	// 		tokens = append(tokens, target.Base.makeToken(TOKEN_COMMENT_TITLE, lexeme))
	// 	}
	//
	// 	if b == NEWLINE {
	// 		tokens = append(tokens, target.Base.makeToken(TOKEN_SINGLE_LINE_COMMENT_END, []byte{0}))
	// 		target.Base.Line++
	// 		break
	// 	}
	//
	// 	offset = 0
	// }
	//
	// we do not need to store comments that are not annotated
	if isAnnotated {
		target.Base.Tokens = append(target.Base.Tokens, tokens...)
	}

	return nil
}

func (target *CLexer) MultiLineComment() error {
	tokens := make([]Token, 0, 10)
	tokens = append(tokens, target.Base.makeToken(TOKEN_MULTI_LINE_COMMENT_START, []byte("/*")))

	offset, lineBroken, isAnnotated := 0, false, false
	for !target.Base.end() {
		b := target.Base.next()

		if b == ASTERISK && target.Base.peekNext() == FORWARD_SLASH {
			target.Base.next()
			if target.Base.next() == NEWLINE {
				target.Base.Line++
			}
			lexeme := []byte{ASTERISK, FORWARD_SLASH}
			tokens = append(tokens, target.Base.makeToken(TOKEN_MULTI_LINE_COMMENT_END, lexeme))
			break
		}

		if !unicode.IsSpace(rune(b)) {
			offset++
			continue
		}

		lexeme, err := target.Base.extractRange(offset)
		if err != nil {
			return target.Base.report(err.Error())
		}

		offset = 0

		if len(lexeme) == 0 || (len(lexeme) == 1 && lexeme[0] == ASTERISK) {
			continue
		}

		switch {
		case !isAnnotated && bytes.Equal(lexeme, target.Base.Annotation):
			tokens = append(tokens, target.Base.makeToken(TOKEN_COMMENT_ANNOTATION, lexeme))
			isAnnotated = true
		case !lineBroken:
			tokens = append(tokens, target.Base.makeToken(TOKEN_COMMENT_TITLE, lexeme))
		default:
			tokens = append(tokens, target.Base.makeToken(TOKEN_COMMENT_DESCRIPTION, lexeme))
		}

		if b == NEWLINE {
			target.Base.Line++
			lineBroken = true
		}
	}

	if target.Base.end() {
		return fmt.Errorf(
			"Failed to parse multi line comment in: %s",
			target.Base.Src[target.Base.Start:target.Base.Current],
		)
	}

	if isAnnotated {
		target.Base.Tokens = append(target.Base.Tokens, tokens...)
	}

	return nil
}

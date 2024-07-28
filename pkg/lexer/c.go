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
	"bytes"
	"fmt"
)

var (
	cSLCommentNotation      = []byte("//") // Single Line Comment notation for c-like languages
	cMLCommentNotationStart = []byte("/*") // opening Multi Line Comment notation
	cMLCommentNotationEnd   = []byte("*/") // closing Multi Line Comment notation
	// new lines of an Multi Line Comment usually start with an asterisk
	cMLCommentSeparator = []byte("*")
)

type Clexer struct {
	Base        *Lexer  // stores utility methods for consuming bytes
	DraftTokens []Token // if a token contains the issue annotation, DraftTokens are appended to base.Tokens
	annotated   bool    // indicator to denote when an issue annotation is located
	line        int     // current line number when consuming a single/multi line opening byte
}

func (c *Clexer) AnalyzeToken() error {
	currentByte := c.Base.peek()
	switch currentByte {
	case QUOTE, DOUBLE_QUOTE, BACK_TICK:
		return c.String(currentByte)
	case FORWARD_SLASH:
		return c.Comment()
	case NEWLINE:
		c.Base.Line++
		return nil
	default:
		return nil
	}
}

func (c *Clexer) String(delim byte) error {
	return nil
}

func (c *Clexer) Comment() error {
	switch c.Base.peekNext() {
	case FORWARD_SLASH:
		return c.tokenizeSLComment()
	case ASTERISK:
		return c.tokenizeMLComment()
	default:
		return nil
	}
}

// creates tokens of the single line comment notation and all text
// that comes before the next line break
func (c *Clexer) tokenizeSLComment() error {
	for !c.Base.pastEnd() {
		if c.Base.peek() == NEWLINE {
			c.Base.Line++
		}

		lexeme := c.Base.nextLexeme()
		if err := c.classifyToken(lexeme, TOKEN_SINGLE_LINE_COMMENT); err != nil {
			return err
		}

		if next := c.Base.peekNext(); next == NEWLINE || next == 0 {
			c.closeSLComment()
		}

		if c.breakTokenization() {
			break
		} else {
			c.Base.next()
		}
	}

	if c.annotated {
		c.promoteTokens()
	}

	return nil
}

// creates tokens of the multi line comment notation and all text
// in between the opening and closing comment notation.
func (c *Clexer) tokenizeMLComment() error {
	c.line = c.Base.Line

	for !c.Base.pastEnd() {
		if c.Base.peek() == NEWLINE {
			c.Base.Line++
		}

		lexeme := c.Base.nextLexeme()
		if err := c.classifyToken(lexeme, TOKEN_MULTI_LINE_COMMENT); err != nil {
			return err
		}

		if c.breakTokenization() {
			break
		} else {
			c.Base.next()
		}
	}

	if c.annotated {
		c.promoteTokens()
	}

	return nil
}

func (c *Clexer) classifyToken(lexeme []byte, tokenType TokenType) error {
	if len(lexeme) == 0 {
		return nil
	}

	token := NewToken(TOKEN_UNKNOWN, lexeme, c.Base)
	return c.classifyTokenType(&token, tokenType)
}

func (c *Clexer) classifyTokenType(token *Token, target TokenType) error {
	isSLComment := containsBits(TOKEN_SINGLE_LINE_COMMENT, target)
	isMLComment := !isSLComment && containsBits(TOKEN_MULTI_LINE_COMMENT, target)

	switch {
	case isSLComment:
		c.classifySLComment(token)
	case isMLComment:
		// for now, we are not going to make separator tokens (*)
		if bytes.Equal(cMLCommentSeparator, token.Lexeme) {
			return nil
		} else {
			c.classifyMLComment(token)
		}
	default:
		return c.reportClassificationError(target)
	}

	c.DraftTokens = append(c.DraftTokens, *token)
	return nil
}

func (c *Clexer) classifySLComment(token *Token) {
	if c.isCommonTokenType(token) {
		return
	} else {
		token.Type = TOKEN_COMMENT_TITLE
	}
}

func (c *Clexer) classifyMLComment(token *Token) {
	if c.isCommonTokenType(token) {
		return
	}

	lineDelta := c.Base.Line - c.line
	if lineDelta == 0 || lineDelta == 1 {
		token.Type = TOKEN_COMMENT_TITLE
	} else {
		token.Type = TOKEN_COMMENT_DESCRIPTION
	}
}

func (c *Clexer) isCommonTokenType(token *Token) bool {
	switch {
	case !c.annotated && bytes.Equal(c.Base.Annotation, token.Lexeme):
		token.Type = TOKEN_COMMENT_ANNOTATION
		c.annotated = true
		return true
	case bytes.Equal(cSLCommentNotation, token.Lexeme):
		token.Type = TOKEN_SINGLE_LINE_COMMENT_START
		return true
	case bytes.Equal(cMLCommentNotationStart, token.Lexeme):
		token.Type = TOKEN_MULTI_LINE_COMMENT_START
		return true
	case bytes.Equal(cMLCommentNotationEnd, token.Lexeme):
		token.Type = TOKEN_MULTI_LINE_COMMENT_END
		return true
	default:
		return false
	}
}

// closeSLComment peeks at the next byte in base.Src. If said byte is a
// new line or 0 (end of src file), then we will append the DraftTokens
// slice with a closing single line comment token.
func (c *Clexer) closeSLComment() {
	next := c.Base.peekNext()

	if next == NEWLINE {
		c.Base.next()
		c.Base.Line++
	}

	c.Base.resetStartIndex()
	token := Token{
		Type:   TOKEN_SINGLE_LINE_COMMENT_END,
		Lexeme: []byte{0},
		Start:  c.Base.Start + 1,
		End:    c.Base.Current + 1,
		Line:   c.Base.Line,
	}

	c.DraftTokens = append(c.DraftTokens, token)
}

func (c *Clexer) breakTokenization() bool {
	if len(c.DraftTokens) == 0 {
		return false
	}

	last := c.DraftTokens[len(c.DraftTokens)-1]
	return containsBits(last.Type, TOKEN_SINGLE_LINE_COMMENT_END) ||
		containsBits(last.Type, TOKEN_MULTI_LINE_COMMENT_END)
}

// promoteTokens is only called when we have found an issue annotation.
// all draft tokens are promoted to the Base Lexers token slice and are
// considered valid comment tokens that the program can safely use
func (c *Clexer) promoteTokens() {
	c.Base.resetStartIndex()
	c.Base.Tokens = append(c.Base.Tokens, c.DraftTokens...)
	c.reset()
}

func (c *Clexer) reportClassificationError(target TokenType) error {
	msg := fmt.Sprintf(
		"classification error: should have TOKEN_SINGLE_LINE_COMMENT or TOKEN_MULTI_LINE_COMMENT but got %s",
		decodeTokenType(target),
	)
	return c.Base.reportError(msg)
}

func (c *Clexer) reset() {
	c.annotated = false
	c.DraftTokens = c.DraftTokens[:0]
	c.line = 0
}

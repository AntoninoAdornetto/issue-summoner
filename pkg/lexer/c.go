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

type Clexer struct {
	Base        *Lexer  // holds shared byte consumption methods
	DraftTokens []Token // Unvalidated tokens
	annotated   bool    // Issue annotation indicator
	line        int     // Current Line number
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

// @TEST_TODO Test the CLexer String func
func (c *Clexer) String(delim byte) error {
	for !c.Base.pastEnd() && c.Base.peekNext() != delim {
		b := c.Base.next()
		if b == NEWLINE {
			c.Base.Line++
		}
	}

	c.Base.next()
	return nil
}

func (c *Clexer) Comment() error {
	switch c.Base.peekNext() {
	case FORWARD_SLASH:
		return c.singleLineComment()
	case ASTERISK:
		return c.multiLineComment()
	default:
		return nil
	}
}

func (c *Clexer) initComment(tokenType TokenType) error {
	start, err := c.Base.startCommentLex(tokenType)
	if err != nil {
		return err
	}

	c.DraftTokens = append(c.DraftTokens, start)
	c.Base.next()
	return nil
}

func (c *Clexer) singleLineComment() error {
	if err := c.initComment(TOKEN_SINGLE_LINE_COMMENT_START); err != nil {
		return err
	}

	for !c.Base.pastEnd() {
		lexeme := c.Base.nextLexeme()
		if err := c.processLexeme(lexeme, TOKEN_SINGLE_LINE_COMMENT); err != nil {
			return err
		}

		if next := c.Base.peekNext(); next == NEWLINE || next == 0 {
			c.closeSLComment()
			break
		}

		c.Base.next()
	}

	if c.annotated {
		c.promoteTokens()
	}

	c.reset()
	return nil
}

func (c *Clexer) multiLineComment() error {
	if err := c.initComment(TOKEN_MULTI_LINE_COMMENT_START); err != nil {
		return err
	}

	for !c.Base.pastEnd() {
		currentByte := c.Base.peek()

		if currentByte == NEWLINE {
			c.Base.Line++
		}

		if currentByte == ASTERISK && c.Base.peekNext() == FORWARD_SLASH {
			c.closeMLComment()
			break
		}

		lexeme := c.Base.nextLexeme()
		if err := c.processLexeme(lexeme, TOKEN_MULTI_LINE_COMMENT); err != nil {
			return err
		}

		c.Base.next()
	}

	if c.annotated {
		c.promoteTokens()
	}

	c.reset()
	return nil
}

func (c *Clexer) processLexeme(lexeme []byte, commentType TokenType) error {
	if len(lexeme) == 0 {
		return nil
	}

	tokens, err := c.Base.processAnnotation(lexeme, c.annotated)
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		c.DraftTokens = append(c.DraftTokens, tokens...)
		c.annotated = true
		c.line = c.Base.Line
		return nil
	}

	switch commentType {
	case TOKEN_SINGLE_LINE_COMMENT:
		c.processSingleLineComment(lexeme)
	case TOKEN_MULTI_LINE_COMMENT:
		c.processMultiLineComment(lexeme)
	default:
		return fmt.Errorf(errTargetTokenize, string(lexeme), decodeTokenType(commentType))
	}

	return nil
}

func (c *Clexer) processSingleLineComment(lexeme []byte) {
	token := NewToken(TOKEN_COMMENT_TITLE, lexeme, c.Base)
	c.DraftTokens = append(c.DraftTokens, token)
}

func (c *Clexer) processMultiLineComment(lexeme []byte) {
	// ignore multi line comment separator (*)
	if bytes.Equal(lexeme, []byte{'*'}) {
		return
	}

	// lineDelta remains at 0 until an issue annotation is located.
	// this is helpful because we know that subsequent lines will
	// part of the comments description
	lineDelta := c.Base.Line - c.line

	var token Token
	if lineDelta == 0 {
		token = NewToken(TOKEN_COMMENT_TITLE, lexeme, c.Base)
	} else {
		token = NewToken(TOKEN_COMMENT_DESCRIPTION, lexeme, c.Base)
	}

	c.DraftTokens = append(c.DraftTokens, token)
}

func (c *Clexer) closeSLComment() {
	next := c.Base.next()
	c.Base.resetStartIndex()
	lexeme := make([]byte, 1)

	if next == NEWLINE {
		c.Base.Line++
		lexeme[0] = NEWLINE
	} else {
		lexeme[0] = 0
	}

	token := NewToken(TOKEN_SINGLE_LINE_COMMENT_END, lexeme, c.Base)
	c.DraftTokens = append(c.DraftTokens, token)
}

func (c *Clexer) closeMLComment() {
	c.Base.resetStartIndex()
	c.Base.next()
	lexeme := []byte{ASTERISK, FORWARD_SLASH}
	token := NewToken(TOKEN_MULTI_LINE_COMMENT_END, lexeme, c.Base)
	c.DraftTokens = append(c.DraftTokens, token)
}

// promoteTokens is invoked when c.annotated is true. Meaning, an issue
// annotation was discovered within the comment and it is safe to append all
// current DraftTokens into the Base Lexers primary token slice.
func (c *Clexer) promoteTokens() {
	c.Base.resetStartIndex()
	c.Base.Tokens = append(c.Base.Tokens, c.DraftTokens...)
}

func (c *Clexer) reset() {
	c.annotated = false
	c.DraftTokens = c.DraftTokens[:0]
	c.line = 0
}

func derivedFromC(ext string) bool {
	switch ext {
	case ".c",
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
		".scala":
		return true
	default:
		return false
	}
}

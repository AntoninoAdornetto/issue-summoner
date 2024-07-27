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
	Base        *Lexer  // stores utility methods for consuming bytes
	DraftTokens []Token // if a token contains the issue annotation, DraftTokens are appended to base.Tokens
	annotated   bool    // indicator to denote when an issue annotation is located
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
	c.reset()
	switch c.Base.peekNext() {
	case FORWARD_SLASH:
		return c.tokenizeSingleLineComment()
	default:
		return nil
	}
}

func (c *Clexer) tokenizeSingleLineComment() error {
	for !c.Base.pastEnd() && c.Base.peek() != NEWLINE {
		lexeme := c.Base.nextLexeme()
		if len(lexeme) == 0 {
			c.Base.next()
			continue
		}

		token := NewToken(TOKEN_UNKNOWN, lexeme, c.Base)
		if err := c.annotateTokenType(&token, TOKEN_SINGLE_LINE_COMMENT); err != nil {
			return err
		}
		c.Base.next()
	}

	if c.annotated {
		c.Base.resetStartIndex()
		c.closeSingleLineCommentToken()
		c.Base.Tokens = append(c.Base.Tokens, c.DraftTokens...)
	}

	return nil
}

func (c *Clexer) annotateTokenType(token *Token, tokenType TokenType) error {
	isSingle := containsBits(tokenType, TOKEN_SINGLE_LINE_COMMENT)
	isMulti := !isSingle && containsBits(tokenType, TOKEN_MULTI_LINE_COMMENT)

	switch {
	case isSingle:
		c.annotateSingleLineComment(token)
	case isMulti:
	default:
		return c.reportTokenTypeError(tokenType)
	}

	c.DraftTokens = append(c.DraftTokens, *token)
	return nil
}

var (
	cSingleLineCommentNotation = []byte("//")
)

func (c *Clexer) annotateSingleLineComment(token *Token) {
	switch {
	case !c.annotated && bytes.Equal(c.Base.Annotation, token.Lexeme):
		token.Type = TOKEN_COMMENT_ANNOTATION
		c.annotated = true
	case bytes.Equal(cSingleLineCommentNotation, token.Lexeme):
		token.Type = TOKEN_SINGLE_LINE_COMMENT_START
	default:
		token.Type = TOKEN_COMMENT_TITLE
	}
}

func (c *Clexer) closeSingleLineCommentToken() {
	token := NewToken(TOKEN_SINGLE_LINE_COMMENT_END, make([]byte, 0, 1), c.Base)

	if c.Base.peek() == NEWLINE {
		c.Base.Line++
		token.Lexeme = append(token.Lexeme, byte('\n'))
	} else {
		token.Lexeme = append(token.Lexeme, 0)
	}

	c.DraftTokens = append(c.DraftTokens, token)
}

}

}

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
	cMLCommentNotationStart = []byte("/*") // Multi Line Comment prefix notation
	cMLCommentSeparator     = []byte("*")  // Multi Line Comment separator
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
		return c.tokenizeSLComment()
	case ASTERISK:
		return c.tokenizeMLComment()
	default:
		return nil
	}
}

func (c *Clexer) tokenizeSLComment() error {
	for !c.Base.pastEnd() {
		lexeme := c.Base.nextLexeme()
		// @TESTING finish the impl
		if !c.annotated && bytes.HasPrefix(lexeme, c.Base.Annotation) {
			fmt.Println("___DEBUG____ CONTAINS PREFIX __________", string(lexeme))
			tkns, _ := c.Base.nextIssueTokens(lexeme)
			for _, t := range tkns {
				fmt.Println("LEXEME: ", string(t.Lexeme))
				fmt.Println("START: ", t.Start)
				fmt.Println("END: ", t.End)
			}
		}

		if err := c.classifyToken(lexeme, TOKEN_SINGLE_LINE_COMMENT); err != nil {
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

func (c *Clexer) tokenizeMLComment() error {
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
		if err := c.classifyToken(lexeme, TOKEN_MULTI_LINE_COMMENT); err != nil {
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

func (c *Clexer) classifyToken(lexeme []byte, tokenType TokenType) error {
	if len(lexeme) == 0 {
		return nil
	}

	token := NewToken(TOKEN_UNKNOWN, lexeme, c.Base)
	return c.classifyCommentToken(&token, tokenType)
}

func (c *Clexer) classifyCommentToken(token *Token, target TokenType) error {
	switch {
	case containsBits(target, TOKEN_SINGLE_LINE_COMMENT):
		c.classifySLComment(token)
	case containsBits(target, TOKEN_MULTI_LINE_COMMENT):
		// not concerned with separators... for now at least
		if bytes.Equal(token.Lexeme, cMLCommentSeparator) {
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
	if !c.isCommonTokenType(token) {
		token.Type = TOKEN_COMMENT_TITLE
	}
}

func (c *Clexer) classifyMLComment(token *Token) {
	lineDelta := c.Base.Line - c.line
	// lineDelta remains at 0 until an issue annotation is located.
	// this is helpful because we know that subsequent lines will
	// part of the comments description and thus allow us to classify
	// it's type correctly

	switch {
	case c.isCommonTokenType(token):
		break
	case lineDelta == 0:
		token.Type = TOKEN_COMMENT_TITLE
	default:
		token.Type = TOKEN_COMMENT_DESCRIPTION
	}
}

func (c *Clexer) isCommonTokenType(token *Token) bool {
	switch {
	case !c.annotated && c.Base.matchAnnotation(token):
		c.line = c.Base.Line
		token.Type = TOKEN_COMMENT_ANNOTATION
		c.annotated = true
		return true
	case bytes.Equal(token.Lexeme, cSLCommentNotation):
		token.Type = TOKEN_SINGLE_LINE_COMMENT_START
		return true
	case bytes.Equal(token.Lexeme, cMLCommentNotationStart):
		token.Type = TOKEN_MULTI_LINE_COMMENT_START
		return true
	default:
		return false
	}
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

func (c *Clexer) reportClassificationError(target TokenType) error {
	msg := fmt.Sprintf(
		"wanted TOKEN_SINGLE_LINE_COMMENT or TOKEN_MULTI_LINE_COMMENT by got %s",
		decodeTokenType(target),
	)
	return c.Base.reportError(msg)
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

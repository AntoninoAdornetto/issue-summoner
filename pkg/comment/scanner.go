package comment

import "fmt"

type TokenType = int

const (
	TOKEN_FORWARD_SLASH TokenType = iota
	TOKEN_ASTERISK
	TOKEN_ERROR
	TOKEN_EOF
)

type CommentScanner struct {
	Start, Current []byte
	Line           int
}

type Token struct {
	tokenType    TokenType
	Start        []byte
	length, line int
}

func NewCommentScanner(src []byte) *CommentScanner {
	return &CommentScanner{Start: src, Current: src, Line: 1}
}

func (c *CommentScanner) ScanTokens() {
	line := -1
	for {
		token := c.ScanToken()
		if token.line != line {
			fmt.Printf("%4d", token.line)
			line = token.line
		} else {
			fmt.Printf(" | ")
		}
		fmt.Printf("%2d '%.*s'\n", token.tokenType, token.length, token.Start)
	}
}

func (c *CommentScanner) ScanToken() Token {
	c.Start = c.Current

	if c.isAtEnd() {
		return *NewToken(TOKEN_EOF, c)
	}

	char := c.advance()
	fmt.Printf("%s", string(char))
	switch char {
	case '/':
		return *NewToken(TOKEN_FORWARD_SLASH, c)
	case '*':
		return *NewToken(TOKEN_ASTERISK, c)
	}

	return newErrorToken("Unexpected character.", c)
}

func (c *CommentScanner) advance() byte {
	start := len(c.Current) - len(c.Start)
	fmt.Println("START NUM ", start)
	c.Current = c.Current[start : start+1]
	return c.Current[len(c.Current)-1]
}

func NewToken(tokenType TokenType, c *CommentScanner) *Token {
	return &Token{
		tokenType: tokenType,
		Start:     c.Start,
		length:    len(c.Current) - len(c.Start),
		line:      c.Line,
	}
}

func newErrorToken(msg string, c *CommentScanner) Token {
	return Token{
		tokenType: TOKEN_ERROR,
		Start:     []byte(msg),
		length:    len(msg),
		line:      c.Line,
	}
}

func (c *CommentScanner) isAtEnd() bool {
	return len(c.Current) == 0
}

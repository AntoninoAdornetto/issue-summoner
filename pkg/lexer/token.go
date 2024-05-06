package lexer

type TokenType = int

const (
	SINGLE_LINE_COMMENT = iota
	MULTI_LINE_COMMENT
	STRING
	EOF
)

type Token struct {
	TokenType      TokenType
	Lexeme         []byte
	Line           int
	StartByteIndex int
	EndByteIndex   int
}

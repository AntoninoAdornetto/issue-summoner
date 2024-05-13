package lexer

type TokenType = int

const (
	SINGLE_LINE_COMMENT = iota
	MULTI_LINE_COMMENT
	STRING
	EOF
)

var (
	ASTERISK       byte = '*'
	BACK_TICK      byte = '`'
	BACKWARD_SLASH byte = '\\'
	FORWARD_SLASH  byte = '/'
	HASH           byte = '#'
	QUOTE          byte = '\''
	DOUBLE_QUOTE   byte = '"'
	NEWLINE        byte = '\n'
	TAB            byte = '\t'
	WHITESPACE     byte = ' '
)

type Token struct {
	TokenType      TokenType
	Lexeme         []byte
	Line           int
	StartByteIndex int
	EndByteIndex   int
}

package lexer3

type TokenType = int

const (
	TOKEN_SINGLE_LINE_COMMENT_START TokenType = iota
	TOKEN_SINGLE_LINE_COMMENT_END
	TOKEN_MULTI_LINE_COMMENT_START
	TOKEN_MULTI_LINE_COMMENT_END
	TOKEN_COMMENT_ANNOTATION
	TOKEN_COMMENT_TITLE
	TOKEN_COMMENT_DESCRIPTION
	TOKEN_EOF
)

const (
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

// Token manages comment tokens. Most tokens that normal compilers and interpreters
// would concern themselves with are ignored. Our use case is to scan for comments
// and more specifically, comments that contain an issue annotation, e.g. @TEST_TODO
type Token struct {
	Type   TokenType
	Lexeme []byte
	Line   int // Line Number the token was discovered
	Start  int // Byte start index in lexer src code byte slice field
	End    int // Byte end index  in lexer src code byte slice field
}

// NewToken is used when you want the start/current positions to take control of
// creating a new token. Most of the time we will opt into using makeToken instead
// because the starting point of which we begin parsing may contain white space that
// we want to ignore. makeToken provides more granual control over how tokens are created
func NewToken(tokenType TokenType, l *Lexer) Token {
	return Token{
		Type:   tokenType,
		Lexeme: l.Src[l.Start : l.Current+1],
		Line:   l.Line,
		Start:  l.Start,
		End:    l.Current,
	}
}

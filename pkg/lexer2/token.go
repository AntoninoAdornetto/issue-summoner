package lexer2

type TokenType = int

// @TODO this is a series of comment tokens

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

func NewToken(tokenType TokenType, l *Lexer) Token {
	return Token{
		Type:   tokenType,
		Lexeme: l.Src[l.Start:l.Current],
		Line:   l.Line,
		Start:  l.Start,
		End:    l.Current,
	}
}

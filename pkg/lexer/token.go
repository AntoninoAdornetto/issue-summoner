package lexer

type TokenType = uint16

const (
	TOKEN_SINGLE_LINE_COMMENT_START TokenType = 1 << iota
	TOKEN_SINGLE_LINE_COMMENT_END
	TOKEN_MULTI_LINE_COMMENT_START
	TOKEN_MULTI_LINE_COMMENT_END
	TOKEN_COMMENT_ANNOTATION
	TOKEN_ISSUE_NUMBER
	TOKEN_COMMENT_TITLE
	TOKEN_COMMENT_DESCRIPTION
	TOKEN_SINGLE_LINE_COMMENT
	TOKEN_MULTI_LINE_COMMENT
	TOKEN_OPEN_PARAN
	TOKEN_CLOSE_PARAN
	TOKEN_HASH
	TOKEN_UNKNOWN
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
	OPEN_PARAN     byte = '('
	CLOSE_PARAN    byte = ')'
	WHITESPACE     byte = ' '
)

type Token struct {
	Type   TokenType
	Lexeme []byte // token value
	Line   int    // Line number
	Start  int    // Starting byte index of the token in Lexer Src slice
	End    int    // Ending byte index of the token in Lexer Src slice
}

func NewToken(tokenType TokenType, lexeme []byte, lexer *Lexer) Token {
	return Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   lexer.Line,
		Start:  lexer.Start,
		End:    lexer.Current,
	}
}

func newEofToken(lexer *Lexer) Token {
	pos := len(lexer.Src) - 1
	return Token{
		Lexeme: []byte{0},
		Type:   TOKEN_EOF,
		Line:   lexer.Line,
		Start:  pos,
		End:    pos,
	}
}

// newPosToken accepts params for creating a new token but instead of relying on
// positions from a lexer, it allows you to specify index & line number positions.
// The functions that utilize newPosToken have a lot of testing to ensure the locations
// of [Start], [End], and [Line] are correct. The primary use case for this func is
// to handle the issue number tokens for issues that have been reported.
func newPosToken(start, end, line int, lexeme []byte, tokenType TokenType) Token {
	return Token{
		Start:  start,
		End:    end,
		Line:   line,
		Lexeme: lexeme,
		Type:   tokenType,
	}
}

func decodeTokenType(tokenType TokenType) string {
	switch {
	case containsBits(tokenType, TOKEN_SINGLE_LINE_COMMENT_START):
		return "TOKEN_SINGLE_LINE_COMMENT_START"
	case containsBits(tokenType, TOKEN_SINGLE_LINE_COMMENT_END):
		return "TOKEN_SINGLE_LINE_COMMENT_END"
	case containsBits(tokenType, TOKEN_MULTI_LINE_COMMENT_START):
		return "TOKEN_MULTI_LINE_COMMENT_START"
	case containsBits(tokenType, TOKEN_MULTI_LINE_COMMENT_END):
		return "TOKEN_MULTI_LINE_COMMENT_END"
	case containsBits(tokenType, TOKEN_COMMENT_ANNOTATION):
		return "TOKEN_COMMENT_ANNOTATION"
	case containsBits(tokenType, TOKEN_COMMENT_TITLE):
		return "TOKEN_COMMENT_TITLE"
	case containsBits(tokenType, TOKEN_COMMENT_DESCRIPTION):
		return "TOKEN_COMMENT_DESCRIPTION"
	case containsBits(tokenType, TOKEN_SINGLE_LINE_COMMENT):
		return "TOKEN_SINGLE_LINE_COMMENT"
	case containsBits(tokenType, TOKEN_MULTI_LINE_COMMENT):
		return "TOKEN_MULTI_LINE_COMMENT"
	case containsBits(tokenType, TOKEN_EOF):
		return "TOKEN_EOF"
	case containsBits(tokenType, TOKEN_OPEN_PARAN):
		return "TOKEN_OPEN_PARAN"
	case containsBits(tokenType, TOKEN_CLOSE_PARAN):
		return "TOKEN_CLOSE_PARAN"
	case containsBits(tokenType, TOKEN_HASH):
		return "TOKEN_HASH"
	case containsBits(tokenType, TOKEN_ISSUE_NUMBER):
		return "TOKEN_ISSUE_NUMBER"
	default:
		return "TOKEN_UNKNOWN"
	}
}

func containsBits(a, b TokenType) bool {
	return a&b != 0
}

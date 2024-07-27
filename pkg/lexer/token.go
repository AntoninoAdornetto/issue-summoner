package lexer

type TokenType = uint16

const (
	TOKEN_SINGLE_LINE_COMMENT_START TokenType = 1 << iota
	TOKEN_SINGLE_LINE_COMMENT_END
	TOKEN_MULTI_LINE_COMMENT_START
	TOKEN_MULTI_LINE_COMMENT_END
	TOKEN_COMMENT_ANNOTATION
	TOKEN_COMMENT_TITLE
	TOKEN_COMMENT_DESCRIPTION
	TOKEN_SINGLE_LINE_COMMENT
	TOKEN_MULTI_LINE_COMMENT
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
	WHITESPACE     byte = ' '
)

type Token struct {
	Type   TokenType
	Lexeme []byte // token value
	Line   int    // Line number
	Start  int    // Starting byte index of the token in Lexer Src slice
	End    int    // Ending byte index of the token in Lexer Src slice
}

func (t *Token) ParseSingleLineCommentToken(annotation []byte, trim func(r rune) bool) Comment {
	loc := findAnnotationLocations(annotation, t.Lexeme)
	if loc == nil {
		return Comment{}
	}
	end := loc[1]
	title := bytes.TrimFunc(t.Lexeme[end:], trim)
	return Comment{
		Title:  title,
		Source: t.Lexeme,
	}
}

func (t *Token) ParseMultiLineCommentToken(annotation []byte, trim func(r rune) bool) Comment {
	loc := findAnnotationLocations(annotation, t.Lexeme)
	if loc == nil {
		return Comment{}
	}
	end := loc[1]
	content := bytes.TrimFunc(t.Lexeme[end:], trim)
	newLines := bytes.Split(content, []byte("\n"))

	comment := Comment{
		Title:  bytes.TrimFunc(newLines[0], trim),
		Source: t.Lexeme,
	}

	for i := 1; i < len(newLines); i++ {
		comment.Description = append(
			comment.Description,
			bytes.TrimLeftFunc(newLines[i], trim)...)

		if i != len(newLines)-1 {
			comment.Description = append(comment.Description, ' ')
		}
	}

	return comment
}

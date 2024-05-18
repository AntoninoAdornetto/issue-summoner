package lexer

import "bytes"

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

package lexer

import (
	"bytes"
)

type Comment struct {
	Title, Description   string
	TokenStartIndex      int
	TokenAnnotationIndex int
	TokenEndIndex        int
	LineNumber           int
	AnnotationPos        []int
}

// Notation contains the comment syntax for the target
// lexer programming language. Meaning, if lexing comments
// for a c-like language, we would use "//", "/*", "*/"
// for python, "#", "”'", etc...
type Notation struct {
	SingleLineStart [][]byte
	SingleLineEnd   [][]byte
	MultiLineStart  [][]byte
	MultiLineEnd    [][]byte
}

type CommentManager struct {
	Comments []Comment
}

func BuildComments(tokens []Token) CommentManager {
	manager := CommentManager{Comments: make([]Comment, 0, 10)}

	for i := 0; i < len(tokens); i++ {
		cur := tokens[i]
		if cur.Type == TOKEN_SINGLE_LINE_COMMENT_START ||
			cur.Type == TOKEN_MULTI_LINE_COMMENT_START {
			manager.iterCommentEnd(tokens, &i)
		}
	}

	return manager
}

func (m *CommentManager) iterCommentEnd(tokens []Token, index *int) {
	token := tokens[*index]
	comment := Comment{LineNumber: token.Line}
	title, description := make([][]byte, 0), make([][]byte, 0)

	for ; *index < len(tokens); *index++ {
		token = tokens[*index]
		switch token.Type {
		case TOKEN_SINGLE_LINE_COMMENT_START, TOKEN_MULTI_LINE_COMMENT_START:
			comment.TokenStartIndex = *index
		case TOKEN_COMMENT_ANNOTATION:
			comment.TokenAnnotationIndex = *index
			comment.AnnotationPos = []int{token.Start, token.End}
		case TOKEN_COMMENT_TITLE:
			title = append(title, token.Lexeme)
		case TOKEN_COMMENT_DESCRIPTION:
			description = append(description, token.Lexeme)
		case TOKEN_SINGLE_LINE_COMMENT_END, TOKEN_MULTI_LINE_COMMENT_END:
			comment.TokenEndIndex = *index
			comment.Title = string(bytes.Join(title, []byte(" ")))
			comment.Description = string(bytes.Join(description, []byte(" ")))
			m.Comments = append(m.Comments, comment)
			return
		}
	}
}

func NewCommentNotation(slStart, slEnd, mlStart, mlEnd [][]byte) *Notation {
	return &Notation{
		SingleLineStart: slStart,
		SingleLineEnd:   slEnd,
		MultiLineStart:  mlStart,
		MultiLineEnd:    mlEnd,
	}
}

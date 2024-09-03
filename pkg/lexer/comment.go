package lexer

import (
	"bytes"
	"strconv"
)

type Comment struct {
	TokenAnnotationIndex int
	Title, Description   string
	TokenStartIndex      int   // location of the first token
	TokenEndIndex        int   // location of the last token
	AnnotationPos        []int // start/end index of the annotation
	IssueNumber          int   // will contain a non 0 value if the comment has been reported
	LineNumber           int
	NotationStartIndex   int // index of where the comment starts
	NotationEndIndex     int // index of where the comment ends
}

type CommentManager struct {
	Comments []Comment
}

func BuildComments(tokens []Token) (CommentManager, error) {
	var err error
	manager := CommentManager{Comments: make([]Comment, 0, 10)}

	for i := 0; i < len(tokens); i++ {
		cur := tokens[i]
		if cur.Type == TOKEN_SINGLE_LINE_COMMENT_START ||
			cur.Type == TOKEN_MULTI_LINE_COMMENT_START {
			err = manager.iterCommentEnd(tokens, &i)
		}
	}

	return manager, err
}

func (m *CommentManager) iterCommentEnd(tokens []Token, index *int) error {
	token := tokens[*index]
	comment := Comment{LineNumber: token.Line}
	title, description := make([][]byte, 0), make([][]byte, 0)

	for ; *index < len(tokens); *index++ {
		token = tokens[*index]
		switch token.Type {
		case TOKEN_SINGLE_LINE_COMMENT_START, TOKEN_MULTI_LINE_COMMENT_START:
			comment.TokenStartIndex = *index
			comment.NotationStartIndex = token.Start
		case TOKEN_COMMENT_ANNOTATION:
			comment.TokenAnnotationIndex = *index
			comment.AnnotationPos = []int{token.Start, token.End}
		case TOKEN_COMMENT_TITLE:
			title = append(title, token.Lexeme)
		case TOKEN_COMMENT_DESCRIPTION:
			description = append(description, token.Lexeme)
		case TOKEN_SINGLE_LINE_COMMENT_END, TOKEN_MULTI_LINE_COMMENT_END:
			comment.TokenEndIndex = *index
			comment.NotationEndIndex = token.End
			comment.Title = string(bytes.Join(title, []byte(" ")))
			comment.Description = string(bytes.Join(description, []byte(" ")))
			m.Comments = append(m.Comments, comment)
			return nil
		case TOKEN_ISSUE_NUMBER:
			issueNum, err := strconv.Atoi(string(token.Lexeme))
			if err != nil {
				return err
			}
			comment.IssueNumber = issueNum
		}
	}

	return nil
}

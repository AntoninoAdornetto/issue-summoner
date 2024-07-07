package lexer2

type commentType = string

const (
	COMMENT_TYPE_SINGLE commentType = "single"
	COMMENT_TYPE_MULTI  commentType = "mulit"
)

type Comment struct {
	Tokens      []Token
	Title       string
	Description string
	Type        commentType
}

// ConstructComments will take all the comment tokens we have
// parsed from source code files and create a slice of comments
// that our various commands can use. For example:
// issue-summoner scan will inform the caller of all comment annotations
// that exist in their project. And issue-summoner report will take the list
// of comments and help the user report them to platforms such as github, gitlab, ect...
func ConstructComments(tokens []Token) []Comment {
	// @IMPLEMENT the ConstructComments function to build a list of comments using all located comment tokens
	return []Comment{}
}

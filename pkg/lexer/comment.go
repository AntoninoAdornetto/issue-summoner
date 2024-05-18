package lexer

type Comment struct {
	Title          []byte
	Description    []byte
	TokenIndex     int
	Source         []byte
	SourceFileName string
}

func (c *Comment) Prepare(fileName string, index int) {
	c.TokenIndex = index
	c.SourceFileName = fileName
}

func (c *Comment) Validate() bool {
	return len(c.Source) > 0
}

func (c *Comment) Push(comments *[]Comment, fileName string, tokenIndex int) {
	if c.Validate() {
		c.Prepare(fileName, tokenIndex)
		*comments = append(*comments, *c)
	}
}

package lexer3

import (
	"errors"
	"fmt"
	"path/filepath"
)

const (
	byte_null = 0
)

// @DOCS write descriptions about some of the properties in the Lexer struct
type Lexer struct {
	Src        []byte
	FileName   string
	Tokens     []Token
	Start      int
	Current    int
	Line       int
	Annotation []byte
}

// @DOCS write descriptions about some of the properties in the LexicalAnalyzer interface
type LexicalAnalyzer interface {
	AnaylzeToken() error
	String(delim byte) error
	Comment() error
}

// @DOCS write descriptions about the responsibilities of the Base Lexer
func NewBaseLexer(annotation, src []byte, fileName string) *Lexer {
	return &Lexer{
		Src:        src,
		FileName:   fileName,
		Tokens:     make([]Token, 0),
		Start:      0,
		Current:    0,
		Line:       1,
		Annotation: annotation,
	}
}

// @DOCS write about why target lexers have a token slice passed to them
func NewLexicalAnalyzer(base *Lexer) (LexicalAnalyzer, error) {
	ext := filepath.Ext(base.FileName)
	tokens := make([]Token, 0, 100)
	switch {
	case ext == ".c":
		return &CLexer{Base: base, Tokens: tokens}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported file type of %s. please open a feature request if you would like support.",
			ext,
		)
	}
}

func (base *Lexer) AnalyzeTokens(target LexicalAnalyzer) ([]Token, error) {
	for range base.Src {
		base.Start = base.Current
		if err := target.AnaylzeToken(); err != nil {
			return nil, err
		}

		base.next()
	}

	base.Tokens = append(base.Tokens, NewToken(TOKEN_EOF, base))
	return base.Tokens, nil
}

func (base *Lexer) next() byte {
	if base.end() {
		return byte_null
	}
	base.Current++
	return base.Src[base.Current]
}

func (base *Lexer) peek() byte {
	if base.end() {
		return byte_null
	}
	return base.Src[base.Current]
}

func (base *Lexer) peekNext() byte {
	if base.Current+1 >= len(base.Src) {
		return byte_null
	}
	return base.Src[base.Current+1]
}

func (base *Lexer) end() bool {
	return base.Current >= len(base.Src)-1
}

func (base *Lexer) extractRange(offset int) ([]byte, error) {
	startPos := base.Current - offset
	endPos := base.Current

	if startPos < 0 || endPos > len(base.Src) {
		msg := fmt.Sprintf(
			"Failed to extract token: out of range (start position: %d), (end position: %d) with length of %d",
			startPos,
			endPos,
			len(base.Src),
		)
		return nil, errors.New(msg)
	}

	return base.Src[startPos:endPos], nil
}

// @BUG off by 1 error in the makeToken func
func (base *Lexer) makeToken(tokenType TokenType, lexeme []byte) Token {
	return Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   base.Line,
		Start:  base.Current - len(lexeme),
		End:    base.Current - 1,
	}
}

func (base *Lexer) report(msg string) error {
	return fmt.Errorf("[%s line %d]: Error: %s", base.FileName, base.Line, msg)
}

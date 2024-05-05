/*
The goal for our lexer is not to create a compiler or intrepreter. It's primary purpose
is to scan the raw source code as a series of characters and group them into tokens.
The implementation will ignore many tokens that other scanners would not.

The types of tokens that we are concered about:
Single line comment tokens (such as // for languages that adopted from c or # for python)
Multi line comment tokens (such as /* for c adopted languages and ”' """ for python)
String tokens
End of file tokens

We should handle string tokens due to certain edge cases where a string may contain comment
notation that could cause our tokenizer to fail. The original lexer would fail when working
with specific configuration files in the JS eco system. One example is the eslintrc.js file
where a developer can specifiy ignore patterns as strings. The ignore pattern could include the
same notation as a multi line comment (/*) and would thus cause the lexer to incorrectly create
comment tokens.
*/
package lexer

import (
	"fmt"
	"path/filepath"
)

type Lexer struct {
	Source   []byte
	FileName string
	Tokens   []Token
	Start    int
	Current  int
	Line     int
	Manager  LexingManager
}

// @TODO // replace any with Comment type
type LexingManager interface {
	AnalyzeToken(lexer *Lexer) error
	ParseCommentTokens(lexer *Lexer, annotation []byte) ([]Comment, error)
}

func NewLexer(src []byte, fileName string) (*Lexer, error) {
	ext := filepath.Ext(fileName)
	manger, err := NewLexingManager(ext)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		Source:   src,
		FileName: fileName,
		Start:    0,
		Current:  0,
		Line:     1,
		Manager:  manger,
		Tokens:   make([]Token, 0),
	}, nil
}

func NewLexingManager(ext string) (LexingManager, error) {
	switch {
	case IsAdoptedFromC(ext):
		return &CLexer{}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported file type of %s. please open a feature request if you would like support.",
			ext,
		)
	}
}

func (l *Lexer) AnalyzeTokens() ([]Token, error) {
	for range l.Source {
		s := l.Source[l.Start:l.Current]
		fmt.Println(s)
		l.Start = l.Current
		err := l.Manager.AnalyzeToken(l)
		if err != nil {
			return nil, err
		}
		l.next()
	}
	l.Tokens = append(l.Tokens, Token{TokenType: EOF})
	return l.Tokens, nil
}

func (l *Lexer) isEnd() bool {
	return l.Current >= len(l.Source)-1
}

func (l *Lexer) next() byte {
	if l.isEnd() {
		return 0
	}
	l.Current++
	return l.Source[l.Current]
}

func (l *Lexer) peek() byte {
	if l.isEnd() {
		return 0
	}
	return l.Source[l.Current]
}

func (l *Lexer) peekNext() byte {
	if l.Current+1 >= len(l.Source) {
		return 0
	}
	return l.Source[l.Current+1]
}

func (l *Lexer) addToken(tokenType TokenType, value []byte) {
	l.Tokens = append(l.Tokens, Token{
		TokenType:      tokenType,
		Lexeme:         string(value),
		Line:           l.Line,
		StartByteIndex: l.Start,
		EndByteIndex:   l.Current,
	})
}

func (l *Lexer) report(msg string) error {
	return fmt.Errorf("[%s line %d]: Error: %s", l.FileName, l.Line, msg)
}

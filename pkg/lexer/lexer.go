/*
The goal for our lexer is not to create a compiler or intrepreter. It's primary purpose
is to scan the raw source code as a series of characters and group them into tokens.
The implementation will ignore many tokens that other scanner/lexers would not.

The types of tokens that we are concered about:
Single line comment tokens (such as // for languages that have adopted c like comment syntax or # for python)
Multi line comment tokens (such as /* for c adopted languages and ”' & """ for python)
End of file tokens

We should check for string tokens, but we do not need to create the token or store the lexeme.
The reason for checking strings is so we can prevent certain edge cases from happening.
One example could be where a string contains characters that could be denoted as a comment.
For C like languages that could be a string such as "/*" or "//".
We don't want the lexer to create tokens for strings that may contain comment syntax.

Each language that is supported will need to satisfy the LexingManager interface and support tokenizing
methods for Comments and Strings. This will allow each implementation to utilize the comment notation that
is specific to a language.
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

type LexingManager interface {
	AnalyzeToken(lexer *Lexer) error
	String(lexer *Lexer, delim byte) error
	Comment(lexer *Lexer) error
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
	case ext == ".py":
		return &PythonLexer{}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported file type of %s. please open a feature request if you would like support.",
			ext,
		)
	}
}

func (l *Lexer) AnalyzeTokens() ([]Token, error) {
	for range l.Source {
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
		Lexeme:         value,
		Line:           l.Line,
		StartByteIndex: l.Start,
		EndByteIndex:   l.Current,
	})
}

func (l *Lexer) report(msg string) error {
	return fmt.Errorf("[%s line %d]: Error: %s", l.FileName, l.Line, msg)
}

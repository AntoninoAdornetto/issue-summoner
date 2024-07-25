/*
Copyright © 2024 AntoninoAdornetto

The lexer.go file is responsible for creating a `Base` Lexer, consuming and iterating through bytes
of source code, and determining which `Target` Lexer to use for the Tokenization process.

Base Lexer:
The name Lexer may be a bit misleading for the Base Lexer. There is no strict rule set baked into
the receiver methods. However, the `Base` Lexer has a very important role of sharing byte consumption
methods to `Target` Lexers. For example, we don't want to re-write .next(), .peek() or .nextLexeme()
multiple times for Target Lexers since the logic for said methods are not specific to the Target Lexer
and won't change.

Target Lexer:
Simply put, a `Target` Lexer is the Lexer that handles the Tokenization rule set. For this application,
we are only concerned with creating single and multi line comments. More specifically, we are concerned
with single and multi line comments that contain an issue annotation.

`Target` Lexers are created via the `NewTargetLexer` method. The `Base` Lexer is passed to the function,
via dependency injection, as input and is stored within each `Target` Lexer so that targets can access the
shared byte consumption methods. `Target` Lexers must satisfy the methods contained in the `LexicalTokenizer`
interface. I know I mentioned we are only concerned with Comments in source code but you will notice a requirement
for a `String` method in the interface. We must account for strings to combat an edge case. Let me explain, if we
are lexing a python string that contains a hash character "#" (comment notation symbol), our lexer could very well-
explode. Same could be said for c or go strings that contain 1 or more forward slashes "/". String tokens are not
persisted, just consumed until the closing delimiter is located.

Lastly, it's important to mention how `Target` Lexers are created. When instantiating a new `Base` Lexer,
the src code file path is provided. This path is utilized to read the base file extension.
If the file extension is .c, .go, .cpp, .h ect, then we would return a Target Lexer that supports c-like comment
syntax since they all denote single and multi line comments with the same notation. For .py files, we would return
a PythonLexer and so on.
*/
package lexer

import (
	"fmt"
	"path/filepath"
)

type Lexer struct {
	FilePath   string
	FileName   string
	Src        []byte  // source code bytes
	Tokens     []Token // comment tokens after lexical analysis has been complete
	Start      int     // byte index
	Current    int     // byte index, used in conjunction with Start to construct tokens
	Line       int     // Line number
	Annotation []byte  // issue annotation to search for within comments
}

type LexicalTokenizer interface {
	AnalyzeToken() error
	String(delim byte) error
	Comment() error
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

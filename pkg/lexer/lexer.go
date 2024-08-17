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
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"unicode"
)

type U8 uint8

const (
	FLAG_PURGE U8 = 1 << iota
	FLAG_SCAN
)

type Lexer struct {
	FilePath   string
	FileName   string
	Src        []byte         // source code bytes
	Tokens     []Token        // comment tokens after lexical analysis has been complete
	Start      int            // byte index
	Current    int            // byte index, used in conjunction with Start to construct tokens
	Line       int            // Line number
	Annotation []byte         // issue annotation to search for within comments
	re         *regexp.Regexp // primary use is for purging comments
	flags      U8
}

type LexicalTokenizer interface {
	AnalyzeToken() error
	String(delim byte) error
	Comment() error
}

func NewLexer(annotation, src []byte, filePath string, flags U8) *Lexer {
	lex := &Lexer{
		Src:        src,
		FilePath:   filePath,
		FileName:   filepath.Base(filePath),
		Tokens:     make([]Token, 0, 100),
		Start:      0,
		Current:    0,
		Line:       1,
		Annotation: annotation,
		flags:      flags,
	}

	if flags&FLAG_PURGE != 0 {
		lex.re = regexp.MustCompile(string(annotation))
	}

	return lex
}

func NewTargetLexer(base *Lexer) (LexicalTokenizer, error) {
	ext := filepath.Ext(base.FileName)
	tokens := make([]Token, 0, 100)

	switch {
	case derivedFromC(ext):
		return &Clexer{Base: base, DraftTokens: tokens}, nil
	default:
		// @TODO return a list of supported programming languages when an error is returned from invoking NewTargetLexer
		return nil, fmt.Errorf("unsupported file extension (%s)", ext)
	}
}

func (base *Lexer) AnalyzeTokens(target LexicalTokenizer) ([]Token, error) {
	for base.Current < len(base.Src) {
		base.resetStartIndex()
		if err := target.AnalyzeToken(); err != nil {
			return nil, err
		} else {
			base.next()
		}
	}

	base.Tokens = append(base.Tokens, newEofToken(base))
	return base.Tokens, nil
}

func (base *Lexer) next() byte {
	base.Current++
	if base.pastEnd() {
		return 0
	}
	return base.Src[base.Current]
}

func (base *Lexer) pastEnd() bool {
	return base.Current > len(base.Src)-1
}

func (base *Lexer) peek() byte {
	if base.pastEnd() {
		return 0
	}
	return base.Src[base.Current]
}

func (base *Lexer) peekNext() byte {
	if base.Current+1 > len(base.Src)-1 {
		return 0
	}
	return base.Src[base.Current+1]
}

func (base *Lexer) nextLexeme() []byte {
	base.resetStartIndex()
	lexeme := make([]byte, 0, 10)

	for !unicode.IsSpace(rune(base.peek())) {
		lexeme = append(lexeme, base.peek())
		if base.breakLexemeIter() {
			break
		} else {
			base.next()
		}
	}

	return lexeme
}

func (base *Lexer) breakLexemeIter() bool {
	return base.Current+1 > len(base.Src)-1 || unicode.IsSpace(rune(base.peekNext()))
}

func (base *Lexer) matchAnnotation(token *Token) bool {
	if base.re != nil {
		return base.re.Match(token.Lexeme)
	}
	return bytes.Equal(token.Lexeme, base.Annotation)
}

func (base *Lexer) resetStartIndex() {
	base.Start = base.Current
}

func (base *Lexer) reportError(msg string) error {
	return fmt.Errorf("[%s line %d]: Error: %s", base.FilePath, base.Line, msg)
}

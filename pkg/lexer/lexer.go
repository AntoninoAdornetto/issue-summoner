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
	"errors"
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

var (
	errTargetTokenize   = "failed to tokenize (%s). Want token type of TOKEN_SINGLE_LINE_COMMENT or TOKEN_MULTI_LINE_COMMENT, got %s"
	errTargetTokenizeSl = "failed to tokenize (%s). Want token type of TOKEN_SINGLE_LINE_COMMENT, got %s"
	errStringClose      = "failed to locate closing string delim (%c): %s"
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
	ext        string         // file extension
	flags      U8
}

// AnalyzeToken - checks the current byte from [Lexer.peek()] and determines how we should process the proceeding bytes
// String - tokens from the string method are not stored. It's needed to prevent lexing comment notation within a string
// Comment - The bread and butter of our target lexers. Handles processing single & multi line comments
// processLexeme - transforms the lexeme into a token and appends it to the draft tokens contained in the target lexer struct
type LexicalTokenizer interface {
	AnalyzeToken() error
	String(delim byte) error
	Comment() error
	processLexeme(lexeme []byte, commentType TokenType) error
}

func NewLexer(annotation, src []byte, filePath string, flags U8) *Lexer {
	fileName := filepath.Base(filePath)
	lex := &Lexer{
		Src:        src,
		FilePath:   filePath,
		FileName:   fileName,
		Tokens:     make([]Token, 0, 100),
		Start:      0,
		Current:    0,
		Line:       1,
		Annotation: annotation,
		flags:      flags,
		ext:        filepath.Ext(fileName),
	}

	if flags&FLAG_PURGE != 0 {
		lex.re = regexp.MustCompile(string(annotation))
	}

	return lex
}

func NewTargetLexer(base *Lexer) (LexicalTokenizer, error) {
	tokens := make([]Token, 0, 100)

	switch {
	case derivedFromC(base.ext):
		return &Clexer{Base: base, DraftTokens: tokens}, nil
	case isShell(base.ext) || base.FileName == "Makefile":
		return &ShellLexer{Base: base, DraftTokens: tokens}, nil
	default:
		return nil, fmt.Errorf("failed to create target lexer with file extension of: %s", base.ext)
	}
}

func (base *Lexer) AnalyzeTokens(target LexicalTokenizer) ([]Token, error) {
	for base.Current < len(base.Src) {
		base.resetStartIndex()
		if err := target.AnalyzeToken(); err != nil {
			return nil, base.reportError(err.Error())
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

func (base *Lexer) processAnnotation(lexeme []byte, isAnnotated bool) ([]Token, error) {
	if isAnnotated || len(lexeme) == 0 {
		return []Token{}, nil
	}

	tokens := make([]Token, 0, 5)
	switch base.flags {
	case FLAG_SCAN:
		if base.matchAnnotationScan(lexeme) {
			tokens = append(tokens, NewToken(TOKEN_COMMENT_ANNOTATION, lexeme, base))
		}
	case FLAG_PURGE:
		if base.matchAnnotationPurge(lexeme) {
			base.appendReportedTokens(lexeme, &tokens)
		}
	default:
		return tokens, errors.New("failed to create annotation tokens. Want SCAN or PURGE flag")
	}

	return tokens, nil
}

func (base *Lexer) matchAnnotationScan(lexeme []byte) bool {
	return bytes.Equal(lexeme, base.Annotation)
}

func (base *Lexer) matchAnnotationPurge(lexeme []byte) bool {
	return base.re != nil && base.re.Match(lexeme)
}

// appendReportedTokens accepts a lexeme and token slice as input. It is responbile for building tokens that
// correspond to a reported issue on a source code hosting platform. For example, if we have reported an
// issue to github and the issue number is 432. The issue annotation would be written as @YOUR_ANNOTATION(#432)
// after reporting it using issue-summoner in the source code file the [Annotation] was located in.
// Later on, when we want to check the status of the reported issue, the program will need to locate every
// every issue number, such as (432), and check the status of it. appendReportedTokens uses
// the [re] regexp to match the lexeme against a pattern. Only if there is a match will appendReportedTokens be invoked.
func (base *Lexer) appendReportedTokens(lexeme []byte, tokens *[]Token) {
	index := bytes.Index(lexeme, []byte{OPEN_PARAN})
	start, end := base.Start, (base.Start + index)
	annotation := newPosToken(start, end-1, base.Line, lexeme[:index], TOKEN_COMMENT_ANNOTATION)
	*tokens = append(*tokens, annotation)
	base.processIssueNumberTokens(lexeme, tokens, index)
}

func (base *Lexer) processIssueNumberTokens(lexeme []byte, tokens *[]Token, index int) {
	for ; index < len(lexeme); index++ {
		start := base.Start + index
		end := start

		switch lexeme[index] {
		case OPEN_PARAN:
			base.appendPosToken(start, end, lexeme[index], TOKEN_OPEN_PARAN, tokens)
		case HASH:
			index = base.processHashToken(lexeme, tokens, index)
		case CLOSE_PARAN:
			base.appendPosToken(start, end, lexeme[index], TOKEN_CLOSE_PARAN, tokens)
		}
	}
}

func (base *Lexer) processHashToken(lexeme []byte, tokens *[]Token, index int) int {
	start := base.Start + index
	end := start
	base.appendPosToken(start, end, lexeme[index], TOKEN_HASH, tokens)
	index++

	start = base.Start + index
	issueNumLexeme := make([]byte, 0, 5)
	for index < len(lexeme) && lexeme[index] != CLOSE_PARAN {
		issueNumLexeme = append(issueNumLexeme, lexeme[index])
		index++
	}

	end = (base.Start + index) - 1
	issueNum := newPosToken(start, end, base.Line, issueNumLexeme, TOKEN_ISSUE_NUMBER)
	*tokens = append(*tokens, issueNum)
	return index - 1
}

func (base *Lexer) appendPosToken(start, end int, char byte, tokenType TokenType, tokens *[]Token) {
	token := newPosToken(start, end, base.Line, []byte{char}, tokenType)
	*tokens = append(*tokens, token)
}

func (base *Lexer) initTokenization(tokenType TokenType, draftTokens *[]Token) error {
	startToken, err := base.openCommentToken(tokenType)
	if err != nil {
		return err
	}
	*draftTokens = append(*draftTokens, startToken)
	return nil
}

func (base *Lexer) openCommentToken(tokenType TokenType) (Token, error) {
	var token Token

	if !containsBits(tokenType, TOKEN_SINGLE_LINE_COMMENT_START^TOKEN_MULTI_LINE_COMMENT_START) {
		return token, fmt.Errorf(
			"failed to start comment analysis with token type of %s. Want single or multi line comment start token",
			decodeTokenType(tokenType),
		)
	}

	lexeme := base.nextLexeme()
	if len(lexeme) == 0 {
		return token, errors.New(
			"failed to start comment analysis. Want comment start notation lexeme to have a len greater than 0",
		)
	}

	token = NewToken(tokenType, lexeme, base)
	return token, nil
}

func (base *Lexer) resetStartIndex() {
	base.Start = base.Current
}

// promoteTokens is invoked only when an annotation is located during the tokenization
// process of a single or multi line comment. It serves as a form of validation to keep
// the [Tokens] slice free of tokens that do not contain issues.
func (base *Lexer) promoteTokens(draftTokens []Token) {
	base.Tokens = append(base.Tokens, draftTokens...)
}

func (base *Lexer) reportError(msg string) error {
	return fmt.Errorf("[%s line %d]: Error: %s", base.FilePath, base.Line, msg)
}

package lexer_test

import (
	"testing"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/lexer"
	"github.com/stretchr/testify/require"
)

var annotation = []byte("@TEST_TODO")

// should return a valid lexer when using a c file
func TestNewLexerC(t *testing.T) {
	lm, err := lexer.NewLexer([]byte{}, "main.c")
	require.NoError(t, err)
	require.IsType(t, &lexer.Lexer{}, lm)
}

// should return a valid lexer when using a cpp file
func TestNewLexerCPP(t *testing.T) {
	lm, err := lexer.NewLexer([]byte{}, "main.cpp")
	require.NoError(t, err)
	require.IsType(t, &lexer.Lexer{}, lm)
}

// should return a valid lexer when using a java file
func TestNewLexerJava(t *testing.T) {
	lm, err := lexer.NewLexer([]byte{}, "main.java")
	require.NoError(t, err)
	require.IsType(t, &lexer.Lexer{}, lm)
}

// should return a valid lexer when using a go file
func TestNewLexerGo(t *testing.T) {
	lm, err := lexer.NewLexer([]byte{}, "main.go")
	require.NoError(t, err)
	require.IsType(t, &lexer.Lexer{}, lm)
}

// should return an error when an unsupported file is passed in
func TestNewLexerUnsupported(t *testing.T) {
	lm, err := lexer.NewLexer([]byte{}, "main.unsupported")
	require.Error(t, err)
	require.Nil(t, lm)
}

// should return the lexing manager for c file extension
func TestNewLexingManagerC(t *testing.T) {
	lm, err := lexer.NewLexingManager(".c")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for js file extension
func TestNewLexingManagerJS(t *testing.T) {
	lm, err := lexer.NewLexingManager(".js")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for ts file extension
func TestNewLexingManagerTS(t *testing.T) {
	lm, err := lexer.NewLexingManager(".ts")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for cpp file extension
func TestNewLexingManagerCPP(t *testing.T) {
	lm, err := lexer.NewLexingManager(".cpp")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for java file extension
func TestNewLexingManagerJava(t *testing.T) {
	lm, err := lexer.NewLexingManager(".java")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return the c lexing manager for go file extension
func TestNewLexingManagerGo(t *testing.T) {
	lm, err := lexer.NewLexingManager(".go")
	require.NoError(t, err)
	require.IsType(t, &lexer.CLexer{}, lm)
}

// should return an error when an unsupported file extension is provided
func TestNewLexingManagerUnsupported(t *testing.T) {
	lm, err := lexer.NewLexingManager(".unsupported")
	require.Error(t, err)
	require.Nil(t, lm)
}

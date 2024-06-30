package ui

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	warn_level    = "[WARNING %s]"
	fatal_level   = "[ERROR %s]"
	success_level = "[SUCCESS %s]"
	hint_level    = "[HINT %s]"
)

type Logger struct {
	errorStyle     lipgloss.Style
	successStyle   lipgloss.Style
	warningStyle   lipgloss.Style
	hintStyle      lipgloss.Style
	standardStyle  lipgloss.Style
	debugIndicator bool
}

func NewLogger(debugIndicator bool) *Logger {
	return &Logger{
		errorStyle:     ErrorTextStyle,
		successStyle:   SuccessTextStyle,
		warningStyle:   NoteTextStyle,
		hintStyle:      SecondaryTextStyle,
		standardStyle:  PrimaryTextStyle,
		debugIndicator: debugIndicator,
	}
}

func (l *Logger) Fatal(message string) {
	level := l.errorStyle.Render(fmt.Sprintf(fatal_level, getTimeStamp()))
	fmt.Fprintf(os.Stderr, level, message)

	if l.debugIndicator {
		fmt.Printf("\n%s\n", string(debug.Stack()))
	}

	os.Exit(1)
}

func (l *Logger) Success(message string) {
	level := l.successStyle.Render(fmt.Sprintf(success_level, getTimeStamp()))
	fmt.Printf("%s %s\n", level, message)
}

func (l *Logger) Warning(message string) {
	level := l.warningStyle.Render(fmt.Sprintf(warn_level, getTimeStamp()))
	fmt.Printf("%s %s\n", level, message)
}

func (l *Logger) Hint(message string) {
	level := l.hintStyle.Render(fmt.Sprintf(hint_level, getTimeStamp()))
	fmt.Printf("%s %s\n", level, message)
}

func (l *Logger) Print(message string) {
	fmt.Println(message)
}

func (l *Logger) PrintStdout(message string) {
	fmt.Print(message)
}

func getTimeStamp() string {
	currentTime := time.Now()
	return currentTime.Format("01-02-2006 15:04:05")
}

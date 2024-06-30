package ui

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Logger struct {
	errorStyle     lipgloss.Style
	successStyle   lipgloss.Style
	warningStyle   lipgloss.Style
	tipStyle       lipgloss.Style
	debugIndicator bool
}

func NewLogger(debugIndicator bool) *Logger {
	return &Logger{
		errorStyle:     ErrorTextStyle,
		successStyle:   SuccessTextStyle,
		warningStyle:   NoteTextStyle,
		debugIndicator: debugIndicator,
		tipStyle:       SecondaryTextStyle,
	}
}

func (l *Logger) LogFatal(message string) {
	ts := getTimeStamp()
	errLog := l.errorStyle.Render(fmt.Sprintf("[ERROR %s]", ts))
	fmt.Fprintf(os.Stderr, errLog, message)

	if l.debugIndicator {
		fmt.Printf("\n%s\n", string(debug.Stack()))
	}

	os.Exit(1)
}

func (l *Logger) LogSuccess(message string) {
	fmt.Println(l.successStyle.Render(message))
}

func (l *Logger) LogWarning(message string) {
	fmt.Println(l.warningStyle.Render(message))
}

func (l *Logger) LogTip(message string) {
	fmt.Println(l.tipStyle.Render(message))
}

func getTimeStamp() string {
	currentTime := time.Now()
	return currentTime.Format("01-02-2006 15:04:05")
}

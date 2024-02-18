package ui

import (
	"fmt"
	"os"
)

func LogFatal(message string) {
	fmt.Fprintln(os.Stderr, ErrorTextStyle.Render(message))
	os.Exit(1)
}

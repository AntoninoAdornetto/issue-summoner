package issue

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func PrintIssueDetails(issues []Issue, keyStyle, valStyle lipgloss.Style) {
	for _, issue := range issues {
		fmt.Printf("\n\n")
		paths := strings.Split(issue.FilePath, "/")
		fmt.Println(
			keyStyle.Render("Filename: "),
			valStyle.Render(paths[len(paths)-1]),
		)
		fmt.Println(keyStyle.Render("Title: "), valStyle.Render(issue.Title))
		fmt.Println(
			keyStyle.Render("Description: "),
			valStyle.Render(issue.Description),
		)
		fmt.Println(
			keyStyle.Render("Line number: "),
			valStyle.Render(fmt.Sprintf("%d", issue.LineNumber)),
		)
	}
}

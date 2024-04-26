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
			keyStyle.Render("Start Line number: "),
			valStyle.Render(fmt.Sprintf("%d", issue.StartLineNumber)),
		)

		fmt.Println(
			keyStyle.Render("End Line number: "),
			valStyle.Render(fmt.Sprintf("%d", issue.EndLineNumber)),
		)

		fmt.Println(
			keyStyle.Render("Annotation Line number: "),
			valStyle.Render(fmt.Sprintf("%d", issue.AnnotationLineNumber)),
		)
	}
}

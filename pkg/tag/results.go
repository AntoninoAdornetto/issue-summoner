package tag

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func PrintTagResults(tags []Tag, keyStyle, valStyle lipgloss.Style) {
	for _, tag := range tags {
		fmt.Printf("\n\n")
		fmt.Println(
			keyStyle.Render("Filename: "),
			valStyle.Render(tag.FileInfo.Name()),
		)
		fmt.Println(keyStyle.Render("Title: "), valStyle.Render(tag.Title))
		fmt.Println(
			keyStyle.Render("Description: "),
			valStyle.Render(tag.Description),
		)
		fmt.Println(
			keyStyle.Render("Start Line number: "),
			valStyle.Render(fmt.Sprintf("%d", tag.StartLineNumber)),
		)

		fmt.Println(
			keyStyle.Render("End Line number: "),
			valStyle.Render(fmt.Sprintf("%d", tag.EndLineNumber)),
		)

		fmt.Println(
			keyStyle.Render("Annotation Line number: "),
			valStyle.Render(fmt.Sprintf("%d", tag.AnnotationLineNum)),
		)

		fmt.Println(
			keyStyle.Render("Multi line comment: "),
			valStyle.Render(fmt.Sprintf("%t", tag.IsMultiLine)),
		)

		fmt.Println(
			keyStyle.Render("Single line comment: "),
			valStyle.Render(fmt.Sprintf("%t", tag.IsSingleLine)),
		)
	}
}

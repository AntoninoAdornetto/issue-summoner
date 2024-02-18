/*
Copyright © 2024 Antonino Adornetto

The scan command processes each source file individually and searches for specific tags (actionable comments) that the user specifies.
It respects the `.gitignore` settings and ensures that any files designated as ignored are not scanned.
Finally, a detailed report is presented to the user about the tags that were found during the scan.
*/
package scan

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	tagKeyTextStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color("170")).
			Bold(true)
	tagValTextStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color("190")).
			Italic(true)
)

type ScanManager struct{}

func (ScanManager) Open(fileName string) (*os.File, error) {
	return os.Open(fileName)
}

func (ScanManager) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans source code file(s) and searches for actionable comments",
	Long:  `Scans a local git respository for Tags (actionable comments) and reports findings to the console.`,
	Run: func(cmd *cobra.Command, _ []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			log.Fatalf("Failed to read 'path' flag: %s", err)
		}

		gitIgnorePath, err := cmd.Flags().GetString("gitignorePath")
		if err != nil {
			log.Fatalf("Failed to read 'gitignorePath' flag\n%s", err)
		}

		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			log.Fatalf("Failed to read 'mode' flag\n%s", err)
		}

		tagName, err := cmd.Flags().GetString("tag")
		if err != nil {
			log.Fatalf("Failed to read 'tag' flag\n%s", err)
		}

		if path == "" {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Failed to get working directory\n%s", err)
			}
			path = wd
		}

		if gitIgnorePath == "" {
			gitIgnorePath = filepath.Join(path, tag.GitIgnoreFile)
		}

		scanManager := ScanManager{}
		ignorePatterns, err := tag.ProcessIgnorePatterns(gitIgnorePath, scanManager)
		if err != nil {
			log.Fatal(err)
		}

		tagManager := tag.TagManager{
			TagName: tagName,
			Mode:    mode,
		}

		if err := tagManager.ValidateMode(mode); err != nil {
			log.Fatalf("Unsupported mode %s provided\n%s", mode, err)
		}

		if mode == tag.PendingMode {
			pendingTagManager := tag.PendedTagManager{TagManager: tagManager}
			tags, err := tag.Walk(tag.WalkParams{
				Root:           path,
				TagManager:     &pendingTagManager,
				FileOperator:   scanManager,
				IgnorePatterns: ignorePatterns,
			})

			if err != nil {
				log.Fatalf("Failed to scan your project.\n%s", err)
			}

			fmt.Println(tagValTextStyle.Render(fmt.Sprintf("%d Tags Found", len(tags))))

			for _, t := range tags {
				fmt.Printf("\n\n")
				fmt.Println(
					tagKeyTextStyle.Render("Filename: "),
					tagValTextStyle.Render(t.FileInfo.Name()),
				)
				fmt.Println(tagKeyTextStyle.Render("Title: "), tagValTextStyle.Render(t.Title))
				fmt.Println(
					tagKeyTextStyle.Render("Description: "),
					tagValTextStyle.Render(t.Description),
				)
				fmt.Println(
					tagKeyTextStyle.Render("Start Line number: "),
					tagValTextStyle.Render(fmt.Sprintf("%d", t.StartLineNumber)),
				)

				fmt.Println(
					tagKeyTextStyle.Render("End Line number: "),
					tagValTextStyle.Render(fmt.Sprintf("%d", t.EndLineNumber)),
				)

				fmt.Println(
					tagKeyTextStyle.Render("Annotation Line number: "),
					tagValTextStyle.Render(fmt.Sprintf("%d", t.AnnotationLineNum)),
				)

				fmt.Println(
					tagKeyTextStyle.Render("Multi line comment: "),
					tagValTextStyle.Render(fmt.Sprintf("%t", t.IsMultiLine)),
				)

				fmt.Println(
					tagKeyTextStyle.Render("Single line comment: "),
					tagValTextStyle.Render(fmt.Sprintf("%t", t.IsSingleLine)),
				)
			}
		}
	},
}

func init() {
	ScanCmd.Flags().StringP("path", "p", "", "Path to your local git project.")
	ScanCmd.Flags().StringP("tag", "t", "@TODO", "Actionable comment tag to search for.")
	ScanCmd.Flags().StringP("mode", "m", "P", "Mode: 'I' (Issued) or 'P' (Pending).")
	ScanCmd.Flags().StringP("gitignorePath", "g", "", "Path to .gitignore file.")
}

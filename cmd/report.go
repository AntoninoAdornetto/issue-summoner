/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const (
	UNAUTHORIZED_ERROR = "Please run `issue-summoner authorize` and complete the authorization process. This will allow us to submit issues on your behalf."
)

func init() {
	ReportCmd.Flags().StringP("path", "p", "", "Path to your local git project.")
	ReportCmd.Flags().StringP("tag", "t", "@TODO", "Actionable comment tag to search for.")
	ReportCmd.Flags().StringP("gitignorePath", "g", "", "Path to .gitignore file.")
	ReportCmd.Flags().
		StringP("scm", "s", scm.GH, "What service do you use for managing source code? GitHub, GitLab...")
}

type ReportManager struct{}

func (ReportManager) Open(fileName string) (*os.File, error) {
	return os.Open(fileName)
}

func (ReportManager) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "report issues to a source code management system",
	Long: `Scans source code files for Tag annotations and reports them, as trackable issues,
	to a source code managment system of your choosing.`,
	Run: func(cmd *cobra.Command, _ []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'path' flag\n%s", err).Error())
		}

		gitIgnorePath, err := cmd.Flags().GetString("gitignorePath")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'gitignorePath' flag\n%s", err).Error())
		}

		tagName, err := cmd.Flags().GetString("tag")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'tag' flag\n%s", err).Error())
		}

		sourceCodeManager, err := cmd.Flags().GetString("scm")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'scm' flag\n%s", err).Error())
		}

		isAuthorized, err := scm.CheckForAccess(sourceCodeManager)
		if err != nil {
			if os.IsNotExist(err) {
				ui.LogFatal(UNAUTHORIZED_ERROR)
			} else {
				ui.LogFatal(err.Error())
			}
		}

		if !isAuthorized {
			ui.LogFatal(UNAUTHORIZED_ERROR)
		}

		if path == "" {
			wd, err := os.Getwd()
			if err != nil {
				ui.LogFatal(fmt.Errorf("Failed to get working directory\n%s", err).Error())
			}
			path = wd
		}

		if gitIgnorePath == "" {
			gitIgnorePath = filepath.Join(path, tag.GitIgnoreFile)
		}

		scanManager := ReportManager{}
		ignorePatterns, err := tag.ProcessIgnorePatterns(gitIgnorePath, scanManager)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		tagManager := tag.TagManager{TagName: tagName}

		pendingTagManager := tag.PendedTagManager{TagManager: tagManager}
		tags, err := tag.Walk(tag.WalkParams{
			Root:           path,
			TagManager:     &pendingTagManager,
			FileOperator:   scanManager,
			IgnorePatterns: ignorePatterns,
		})

		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to scan your project.\n%s", err).Error())
		}

		selection := ui.Selection{
			Options: make(map[string]bool),
		}

		tagOptions := make([]ui.Item, len(tags))
		for i, t := range tags {
			tagOptions[i] = ui.Item{
				Title: t.Title,
				Desc:  t.Description,
				ID:    t.Title,
			}
		}

		// @TODO  Review ui.InitialModelMultiSelect params
		// Do we really need to pass a bolean pointer to the model? Look into this.
		ok := err != nil

		teaProgram := tea.NewProgram(
			ui.InitialModelMultiSelect(
				tagOptions,
				&selection,
				"Select all items that you want to report as issues.",
				&ok,
			),
		)

		if _, err := teaProgram.Run(); err != nil {
			cobra.CheckErr(ui.ErrorTextStyle.Render(err.Error()))
		}

		gm := scm.GetGitConfig(sourceCodeManager)
		toReport := make([]scm.Issue, 0)

		for _, t := range tags {
			if selection.Options[t.Title] {
				toReport = append(toReport, scm.Issue{Title: t.Title, Description: t.Description})
			}
		}

		err = gm.Report(toReport, sourceCodeManager)
		if err != nil {
			ui.LogFatal(err.Error())
		}
	},
}

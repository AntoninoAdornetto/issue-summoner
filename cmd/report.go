/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"bytes"
	"runtime"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/AntoninoAdornetto/issue-summoner/templates"
	"github.com/AntoninoAdornetto/issue-summoner/v2/git"
	"github.com/AntoninoAdornetto/issue-summoner/v2/issue"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report pending issues to a source code management platform",
	Long: `Report will scan your git project for comments that include issue annotations.
Issue annotations can be as simple as @TODO or any other value that you seefit. 
The only requirement is that the annotation resides in a single or multi line comment. 
Once issue annotations are discovered, you will be presented with a list of all the issues 
that were located and you can select which ones you would like to report to a source code management
platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		annotation, path := handleCommonFlags(cmd)

		sourceCodeManager, err := cmd.Flags().GetString(flag_scm)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		repo, err := git.NewRepository(path)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		iMan := issue.NewIssueManager(annotation)
		if err := iMan.Walk(repo.WorkTree); err != nil {
			ui.LogFatal(err.Error())
		}

		_, err = git.NewGitManager(sourceCodeManager, repo)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		options := make([]ui.Item, len(iMan.Issues))
		for i, iss := range iMan.Issues {
			options[i] = ui.Item{
				Title: iss.Title,
				Desc:  iss.Description,
				ID:    iss.ID,
			}
		}

		selections := ui.Selection{
			Options: make(map[string]bool),
		}

		var quit bool
		teaProgram := tea.NewProgram(
			ui.InitialModelMultiSelect(options, &selections, select_issues, &quit),
		)

		if _, err := teaProgram.Run(); err != nil {
			ui.LogFatal(err.Error())
		}

		// @TODO remove embedded template, no real need for it. Just create an inline template in the v2 issue package
		tmpl, err := templates.LoadIssueTemplate()
		if err != nil {
			ui.LogFatal(err.Error())
		}

		env := runtime.GOOS
		queue := make([]git.ScmIssue, 0, 5)
		for i, iss := range iMan.Issues {
			if selections.Options[iss.ID] {
				buf := bytes.Buffer{}
				iss.Environment = env
				if err := tmpl.Execute(&buf, iss); err != nil {
					ui.LogFatal(err.Error())
				}
				queue = append(queue, git.ScmIssue{Title: iss.Title, Body: buf.String(), Index: i})
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	reportCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
	reportCmd.Flags().StringP(flag_scm, shortflag_scm, scm.GITHUB, flag_desc_scm)
}

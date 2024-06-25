/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"fmt"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
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

		iMan, err := issue.NewIssueManager(annotation, true)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		if err := iMan.Walk(repo.WorkTree); err != nil {
			ui.LogFatal(err.Error())
		}

		gitManager, err := git.NewGitManager(sourceCodeManager, repo)
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

		queue := make([]issue.Issue, 0, len(iMan.Issues))
		for _, codeIssue := range iMan.Issues {
			if selections.Options[codeIssue.ID] {
				queue = append(queue, codeIssue)
			}
		}

		reportedChan := make(chan git.ReportedIssue)
		for i, codeIssue := range queue {
			go func(item issue.Issue, index int) {
				toReport := git.CodeIssue{Title: item.Title, Body: item.Body, Index: index}
				res, err := gitManager.Report(toReport)
				if err != nil {
					ui.ErrorTextStyle.Render(
						fmt.Sprintf("Error: failed to report issue (%s)", item.Title),
					)
					return
				}

				reportedChan <- res
			}(codeIssue, i)
		}

		done := make(chan bool)
		go func() {
			for range queue {
				rp := <-reportedChan
				fmt.Println(rp.Index)
				fmt.Println(rp.ID)
			}
			done <- true
		}()

		<-done
		fmt.Println(
			ui.SuccessTextStyle.Render(
				fmt.Sprintf(
					"Reported %d issues to %s successfully\n",
					len(queue),
					sourceCodeManager,
				),
			),
		)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	reportCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
	reportCmd.Flags().StringP(flag_scm, shortflag_scm, scm.GITHUB, flag_desc_scm)
}

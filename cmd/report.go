/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/AntoninoAdornetto/issue-summoner/templates"
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

		_, err = scm.ReadAccessToken(sourceCodeManager)
		if err != nil {
			if os.IsNotExist(err) {
				ui.LogFatal(
					"configuration file does not exist. please run <issue-summoner authorize> or see <issue-summoner authorize --help>",
				)
			} else {
				ui.LogFatal(err.Error())
			}
		}

		issueManager, err := issue.NewIssueManager(issue.PENDING_ISSUE, annotation)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		_, err = issueManager.Walk(path)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		issues := issueManager.GetIssues()
		if len(issues) == 0 {
			fmt.Println(ui.ErrorTextStyle.Render(no_issues))
			return
		}

		selections := ui.Selection{
			Options: make(map[string]bool),
		}

		options := make([]ui.Item, len(issues))
		for i, is := range issues {
			options[i] = ui.Item{
				Title: is.Title,
				Desc:  is.Description,
				ID:    is.ID,
			}
		}

		var quit bool
		teaProgram := tea.NewProgram(
			ui.InitialModelMultiSelect(
				options,
				&selections,
				select_issues,
				&quit,
			),
		)

		if _, err := teaProgram.Run(); err != nil {
			ui.LogFatal(err.Error())
		}

		tmpl, err := templates.LoadIssueTemplate()
		if err != nil {
			ui.LogFatal(err.Error())
		}

		reportQueue := make([]scm.GitIssue, 0)
		for i, is := range issues {
			if selections.Options[is.ID] {
				md, err := is.ExecuteIssueTemplate(tmpl)
				if err != nil {
					ui.LogFatal(err.Error())
				}
				reportQueue = append(
					reportQueue,
					scm.GitIssue{Title: is.Title, Body: string(md), QueueIndex: i},
				)
			}
		}

		out := bytes.Buffer{}
		remoteCmd := exec.Command("git", "remote", "-v")
		remoteCmd.Stdout = &out
		if err := remoteCmd.Run(); err != nil {
			ui.LogFatal(err.Error())
		}

		userName, repoName, err := scm.ExtractUserRepoName(out.Bytes())
		if err != nil {
			ui.LogFatal(err.Error())
		}

		gitManager, err := scm.NewGitManager(sourceCodeManager, userName, repoName)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		reported := gitManager.Report(reportQueue)
		for ch := range reported {
			if err := issueManager.WriteIssueID(ch.ID, ch.QueueIndex); err != nil {
				ui.LogFatal(err.Error())
			}
		}

		if len(reportQueue) == 0 {
			return
		}

		fmt.Println(
			ui.SuccessTextStyle.Render(
				fmt.Sprintf(
					"Success! Uploaded %d issue(s) to %s",
					len(reportQueue),
					sourceCodeManager,
				),
			),
		)

		fmt.Println(
			ui.SecondaryTextStyle.Render("make sure to commit and push the annotation updates!"),
		)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	reportCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
	reportCmd.Flags().StringP(flag_scm, shortflag_scm, scm.GITHUB, flag_desc_scm)
}

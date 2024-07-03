/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/git"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	multiselect "github.com/AntoninoAdornetto/issue-summoner/pkg/ui/multi_select"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui/spinner"
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
		annotation, path := getCommonFlags(cmd)
		logger := getLogger(cmd)

		srcCodeManager, err := cmd.Flags().GetString("scm")
		if err != nil {
			logger.Fatal(err.Error())
		}

		repo, err := git.NewRepository(path)
		if err != nil {
			logger.Fatal(err.Error())
		}

		manager, err := issue.NewIssueManager(annotation, issue.ISSUE_MODE_PEND, true)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if err := manager.Walk(repo.WorkTree); err != nil {
			logger.Fatal(err.Error())
		}

		if len(manager.Issues) == 0 {
			logger.Info(no_issues + annotation)
			return
		}

		gitManager, err := git.NewGitManager(srcCodeManager, repo)
		if err != nil {
			logger.Fatal(err.Error())
		}

		options := make([]multiselect.Item, len(manager.Issues))
		for i, toReport := range manager.Issues {
			options[i] = multiselect.Item{
				Title: toReport.Title,
				Desc:  toReport.Description,
				ID:    toReport.Index,
			}
		}

		selections := multiselect.Selection{
			Options: make(map[int]bool),
		}

		var quit bool
		teaProgram := tea.NewProgram(
			multiselect.InitialModelMultiSelect(options, &selections, select_issues, &quit),
		)

		if _, err := teaProgram.Run(); err != nil {
			logger.Fatal(err.Error())
		}

		selectedCount := 0
		for _, selected := range selections.Options {
			if selected {
				selectedCount++
			}
		}

		switch selectedCount {
		case 0:
			logger.Info("No issues selected")
			return
		case 1:
			logger.PrintStdout(
				fmt.Sprintf("\nReport %d issue to %s? (y/n): ", selectedCount, srcCodeManager),
			)
		default:
			logger.PrintStdout(
				fmt.Sprintf("\nReport %d issues to %s? (y/n): ", selectedCount, srcCodeManager),
			)
		}

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			logger.Fatal(scanner.Err().Error())
		}

		if scanner.Text() != "y" {
			logger.Info("Aborting report request")
			return
		}

		spin := tea.NewProgram(
			spinner.InitialModelNew(fmt.Sprintf("Reporting to %s", srcCodeManager)),
		)

		go func() {
			if _, err := spin.Run(); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		defer func() {
			if r := recover(); err != nil {
				logger.Fatal(fmt.Sprintf("Failed to record with unexpected error: %s", r))
			}

			if err := spin.ReleaseTerminal(); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		reported := make(chan git.ReportResponse, selectedCount)
		for index := range selections.Options {
			toReport := manager.Issues[index]
			request := git.ReportRequest{Title: toReport.Title, Body: toReport.Body, Index: index}
			go gitManager.Report(request, reported)
		}

		// @TODO create type for this
		failed := make([]struct {
			err   error
			index int
		}, 0)

		for range selectedCount {
			result := <-reported
			if result.Err != nil {
				failed = append(failed, struct {
					err   error
					index int
				}{err: result.Err, index: result.Index})
			} else {
				currentIssue := manager.Issues[result.Index]
				manager.UpdateMapVal(currentIssue.FilePath, result.Index, result.ID)
			}
		}

		for _, e := range failed {
			logger.Warning(
				fmt.Sprintf(
					"Failed to process request for issue (%s)\tError: %s",
					manager.Issues[e.index].Title,
					e.err.Error(),
				),
			)
		}

		resultCount := selectedCount - len(failed)
		if resultCount == 0 {
			if err := spin.ReleaseTerminal(); err != nil {
				logger.Fatal(err.Error())
			}

			logger.Fatal(
				"All selected issues have failed to report. Try running <issue-summoner authorize> before reporting again to refresh access token & repo scope",
			)
		}

		logger.Info("Updating src code comments with published issue ids")

		// @TODO implement the bulk write operations based on issues that are in the same file
		for filePath, issues := range manager.ReportMap {
			fmt.Println(filePath)
			for _, is := range issues {
				fmt.Println(is.Title)
			}
		}

		logger.Info(
			fmt.Sprintf(
				"%d issues(s) remaining in your queue. run <issue-summoner scan -v> to see what's left",
				len(manager.Issues),
			),
		)

		logger.Success(fmt.Sprintf("%d issue(s) reported to %s", resultCount, srcCodeManager))

	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	reportCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
	reportCmd.Flags().StringP(flag_scm, shortflag_scm, git.GITHUB, flag_desc_scm)
	reportCmd.Flags().BoolP(flag_debug, shortflag_debug, false, flag_desc_debug)
}

/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/git"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type writeIssueResult struct {
	Err     error
	PathKey string
}

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

		srcCodeManager, err := cmd.Flags().GetString(flag_scm)
		if err != nil {
			logger.Fatal(err.Error())
		}

		repo, err := git.NewRepository(path)
		if err != nil {
			logger.Fatal(err.Error())
		}

		manager, err := issue.NewIssueManager([]byte(annotation), issue.IssueModeReport)
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

		options := make([]ui.Item, len(manager.Issues))
		for i, toReport := range manager.Issues {
			options[i] = ui.Item{
				Title: toReport.Title,
				Desc:  toReport.Description,
				ID:    i,
			}
		}

		var quit bool
		selections := ui.Selection{Options: make(map[int]bool)}
		multiSelect := tea.NewProgram(
			ui.InitMultiSelect(options, &selections, select_issues, &quit),
		)

		if _, err := multiSelect.Run(); err != nil {
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
			logger.Fatal(err.Error())
		}

		switch scanner.Text() {
		case "y", "yes", "return":
			break
		default:
			logger.Info("Aborting report request")
			return
		}

		spinner := tea.NewProgram(ui.InitSpinner(fmt.Sprintf("Reporting to %s", srcCodeManager)))
		go func() {
			if _, err := spinner.Run(); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		defer func() {
			if r := recover(); err != nil {
				logger.Fatal(fmt.Sprintf("Failed to recover program with unexpected error: %s", r))
			}

			if err := spinner.ReleaseTerminal(); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		wg := sync.WaitGroup{}
		wg.Add(selectedCount)

		reportedChan := make(chan git.ReportResponse, selectedCount)
		for index := range selections.Options {
			toReport := manager.Issues[index]
			req := git.ReportRequest{Title: toReport.Title, Body: toReport.Body, Index: index}
			go func(request git.ReportRequest) {
				defer wg.Done()
				gitManager.Report(request, reportedChan)
			}(req)
		}

		wg.Wait()
		close(reportedChan)

		for r := range reportedChan {
			if r.Err != nil {
				logger.Warning(r.Err.Error())
			} else {
				manager.Group(r.Index, r.ID)
			}
		}

		reportCount := 0
		for range manager.IssueMap {
			reportCount++
		}

		wg.Add(reportCount)
		writeChan := make(chan writeIssueResult, reportCount)
		for path := range manager.IssueMap {
			go func(filePath string) {
				result := writeIssueResult{PathKey: filePath}
				defer wg.Done()
				if err := manager.WriteIssues(filePath); err != nil {
					result.Err = err
				}
				writeChan <- result
			}(path)
		}

		wg.Wait()
		close(writeChan)
		for r := range writeChan {
			messages, err := manager.Results(r.PathKey, srcCodeManager, r.Err != nil)
			if err != nil {
				logger.Warning(err.Error())
			}

			for _, msg := range messages {
				if r.Err != nil {
					logger.Warning(msg)
				} else {
					logger.Success(msg)
					fmt.Printf("\n")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	reportCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
	reportCmd.Flags().StringP(flag_scm, shortflag_scm, git.Github, flag_desc_scm)
	reportCmd.Flags().BoolP(flag_debug, shortflag_debug, false, flag_desc_debug)
}

/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"fmt"
	"sync"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/git"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scans source code for Issue annotations",
	Long: `Scan provides functionality for managing and reviewing both reported
and un-reported issues that reside in your codebase. It serves as an aid to the report
command through two primary modes, scan and purge mode. These modes help you manage and 
track issues directly within your codebase using custom annotations. The default scan mode
will inform you how many issues are in your codebase that have not been reported. Purge mode
will inform you how many reported issues are in your codebase that are still open. Additonally,
purge mode will check the status of each reported issue and remove the corresponding comment if 
the source code hosting platform indicates it is in a resolved state. Both modes can be used to 
print details about the issues, such as the description and location of the issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		annotation, path := getCommonFlags(cmd)
		logger := getLogger(cmd)

		verbose, err := cmd.Flags().GetBool(flag_verbose)
		if err != nil {
			logger.Fatal(err.Error())
		}

		mode, err := cmd.Flags().GetString(flag_mode)
		if err != nil {
			logger.Fatal(err.Error())
		}

		repo, err := git.NewRepository(path)
		if err != nil {
			logger.Fatal(err.Error())
		}

		manager, err := issue.NewIssueManager([]byte(annotation), mode)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if err := manager.Walk(repo.WorkTree); err != nil {
			logger.Fatal(err.Error())
		}

		if len(manager.Issues) == 0 {
			logger.Success(fmt.Sprintf("Scan finished: %s %s", no_issues, annotation))
			return
		}

		var msg string

		issueCount, purgeCount := len(manager.Issues), 0
		if mode == issue.IssueModePurge {
			wg := sync.WaitGroup{}
			wg.Add(issueCount)

			sourceCodeHost, err := cmd.Flags().GetString(flag_sch)
			if err != nil {
				logger.Fatal(err.Error())
			}

			logger.Info("Checking statuses of reported issues on " + sourceCodeHost)

			gitManager, err := git.NewGitManager(sourceCodeHost, repo)
			if err != nil {
				logger.Fatal(err.Error())
			}

			statusChan := make(chan git.StatusResponse, issueCount)
			for _, iss := range manager.Issues {
				go func(toCheck issue.Issue) {
					defer wg.Done()
					gitManager.GetStatus(iss.Comment.IssueNumber, toCheck.Index, statusChan)
				}(iss)
			}

			wg.Wait()
			close(statusChan)

			for c := range statusChan {
				currentIssue := manager.Issues[c.Index]
				switch {
				case c.Err != nil:
					logger.Warning(
						fmt.Sprintf(
							"failed to get issue status for <%s> with error: %s",
							currentIssue.Title,
							c.Err.Error(),
						),
					)
				case c.Resolved:
					if err := manager.Group(c.Index, currentIssue.Comment.IssueNumber); err != nil {
						logger.Warning("Failed to group")
					}
				case c.Err != nil && !c.Resolved:
					logger.Warning(
						fmt.Sprintf("Issue <%s> has not been resolved yet", currentIssue.Title),
					)
				}
			}

			for pathKey, entries := range manager.IssueMap {
				// @TODO utilize go routines for large purge requests
				if err := manager.Purge(pathKey); err != nil {
					logger.Fatal(err.Error())
				} else {
					purgeCount += len(entries)
				}
			}

			if purgeCount > 0 {
				msg += fmt.Sprintf("Purged %d issue(s). ", purgeCount)
			}

			if issueCount > 0 {
				msg += fmt.Sprintf(
					"%d issue(s) remain open and awaiting resolution ",
					issueCount-purgeCount,
				)
			}
		} else {
			msg = fmt.Sprintf("Found %d issue annotations using %s", len(manager.Issues), annotation)
		}

		if verbose {
			for _, iss := range manager.Issues {
				fmt.Printf("\n\n")

				fmt.Println(
					ui.AccentTextStyle.Render("File name: "),
					ui.PrimaryTextStyle.Render(iss.FileName),
				)

				fmt.Println(
					ui.AccentTextStyle.Render("Title: "),
					ui.PrimaryTextStyle.Render(iss.Title),
				)

				fmt.Println(
					ui.AccentTextStyle.Render("Description: "),
					ui.PrimaryTextStyle.Render(iss.Description),
				)

				fmt.Println(
					ui.AccentTextStyle.Render("Line number: "),
					ui.PrimaryTextStyle.Render(fmt.Sprintf("%d", iss.LineNumber)),
				)

				if mode == issue.IssueModePurge && iss.Comment.IssueNumber != 0 {
					fmt.Println(
						ui.AccentTextStyle.Render("Issue number: "),
						ui.PrimaryTextStyle.Render(fmt.Sprintf("%d", iss.Comment.IssueNumber)),
					)
				}
			}
		}

		logger.Success(msg)

		if !verbose {
			logger.Hint(tip_verbose)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
	scanCmd.Flags().BoolP(flag_debug, shortflag_debug, false, flag_desc_debug)
	scanCmd.Flags().StringP(flag_mode, shortflag_mode, issue.IssueModeScan, flag_desc_mode)
	scanCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	scanCmd.Flags().StringP(flag_sch, shortflag_sch, git.Github, flag_desc_sch)
	scanCmd.Flags().BoolP(flag_verbose, shortflag_verbose, false, flag_desc_verbose)
}

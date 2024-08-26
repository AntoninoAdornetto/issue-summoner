/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"fmt"

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

		msg := fmt.Sprintf("Found %d issue annotations using %s", len(manager.Issues), annotation)

		if verbose {
			for _, issue := range manager.Issues {
				fmt.Printf("\n\n")

				fmt.Println(
					ui.AccentTextStyle.Render("File name: "),
					ui.PrimaryTextStyle.Render(issue.FileName),
				)

				fmt.Println(
					ui.AccentTextStyle.Render("Title: "),
					ui.PrimaryTextStyle.Render(issue.Title),
				)

				fmt.Println(
					ui.AccentTextStyle.Render("Description: "),
					ui.PrimaryTextStyle.Render(issue.Description),
				)

				fmt.Println(
					ui.AccentTextStyle.Render("Line number: "),
					ui.PrimaryTextStyle.Render(fmt.Sprintf("%d", issue.LineNumber)),
				)
			}
		}

		fmt.Printf("\n")
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
	scanCmd.Flags().BoolP(flag_verbose, shortflag_verbose, false, flag_desc_verbose)
}

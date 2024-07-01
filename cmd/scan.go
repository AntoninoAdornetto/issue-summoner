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
	Short: "scans source code for Issue annotations (actionable comments)",
	Long: `scans your git project for comments that include issue annotations.
Issue annotations can be as simple as @TODO or any annotation that you see
fit. The only requirement is that the annotation resides in a single or multi
line comment. Once found, you can see details about the located comment using
the verbose flag. The scan command is a preliminary command to the report 
command. Report will actually publish the located comments to your favorite
source code management platform. Scan is for reviewing the issue annotations
that reside in your code base.`,
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

		logger.Log(fmt.Sprintf("Scanning for comments with %s annotation", annotation))

		repo, err := git.NewRepository(path)
		if err != nil {
			logger.Fatal(err.Error())
		}

		manager, err := issue.NewIssueManager(annotation, mode, false)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if err := manager.Walk(repo.WorkTree); err != nil {
			logger.Fatal(err.Error())
		}

		if len(manager.Issues) == 0 {
			fmt.Println(ui.SecondaryTextStyle.Render(fmt.Sprintf("%s %s", no_issues, annotation)))
			return
		}

		msg := fmt.Sprintf(
			"Found %d issue annoations using %s",
			len(manager.Issues),
			annotation,
		)

		if verbose {
			manager.Print(ui.DimTextStyle, ui.PrimaryTextStyle)
			logger.Success(msg)
		} else {
			logger.Success(msg)
			logger.Hint(tip_verbose)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	scanCmd.Flags().StringP(flag_mode, shortflag_mode, issue.ISSUE_MODE_PEND, flag_desc_mode)
	scanCmd.Flags().BoolP(flag_verbose, shortflag_verbose, false, flag_desc_verbose)
	scanCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
	scanCmd.Flags().BoolP(flag_debug, shortflag_debug, false, flag_desc_debug)
}

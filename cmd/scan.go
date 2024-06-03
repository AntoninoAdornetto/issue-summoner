/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"fmt"

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
		annotation, path := handleCommonFlags(cmd)

		verbose, err := cmd.Flags().GetBool(flag_verbose)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		mode, err := cmd.Flags().GetString(flag_mode)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		issueManager, err := issue.NewIssueManager(mode, annotation)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		_, err = issueManager.Walk(path)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		issues := issueManager.GetIssues()
		if len(issues) > 0 {
			success := fmt.Sprintf("Found %d issue annotations using %s", len(issues), annotation)
			fmt.Println(ui.SuccessTextStyle.Render(success))
		} else {
			fmt.Println(ui.SecondaryTextStyle.Render(fmt.Sprintf("%s %s", no_issues, annotation)))
			return
		}

		if verbose {
			issue.PrintIssueDetails(issues, ui.DimTextStyle, ui.PrimaryTextStyle)
		} else {
			fmt.Println(ui.SecondaryTextStyle.Render(tip_verbose))
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP(flag_path, shortflag_path, "", flag_desc_path)
	scanCmd.Flags().StringP(flag_gignore, shortflag_gignore, "", flag_desc_gignore)
	scanCmd.Flags().StringP(flag_mode, shortflag_mode, issue.PENDING_ISSUE, flag_desc_mode)
	scanCmd.Flags().BoolP(flag_verbose, shortflag_verbose, false, flag_desc_verbose)
	scanCmd.Flags().StringP(flag_annotation, shortflag_annotation, "@TODO", flag_desc_annotation)
}

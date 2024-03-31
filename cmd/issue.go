/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

func init() {
	IssueCmd.Flags().StringP("path", "p", "", "Path to your local git project.")
	IssueCmd.Flags().StringP("tag", "t", "@TODO", "Actionable comment tag to search for.")
	IssueCmd.Flags().StringP("gitignorePath", "g", "", "Path to .gitignore file.")

	IssueCmd.Flags().
		StringP("mode", "m", issue.PENDING_ISSUE, "Mode: 'I' (Issued) or 'P' (Pending).")

	IssueCmd.Flags().
		BoolP("verbose", "v", false, "Displays all information for each Tag annotation that is located")
}

// issueCmd represents the issue command
var IssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'path' flag\n%s", err).Error())
		}

		gitIgnorePath, err := cmd.Flags().GetString("gitignorePath")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'gitignorePath' flag\n%s", err).Error())
		}

		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'mode' flag\n%s", err).Error())
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'verbose' flag\n%s", err).Error())
		}

		tagName, err := cmd.Flags().GetString("tag")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'tag' flag\n%s", err).Error())
		}

		if path == "" {
			wd, err := os.Getwd()
			if err != nil {
				ui.LogFatal(fmt.Errorf("Failed to get working directory\n%s", err).Error())
			}
			path = wd
		}

		issueManager, err := issue.GetIssueManager(mode)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		fmt.Println(issueManager, gitIgnorePath, mode, verbose, tagName)
	},
}

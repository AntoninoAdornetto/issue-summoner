/*
Copyright © 2024 Antonino Adornetto
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/issue"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scans source code for Issue annotations (actionable comments)",
	Long: `scans your git project for comments that include issue annotations.
		The comment is used for reporting purposes to see what actionable comments your
		project contains. This can give you an idea of all the issues in your code base
		prior to uploading them to a source code management platform.
		`,
	Run: func(cmd *cobra.Command, args []string) {
		annotation, err := cmd.Flags().GetString("annotation")
		if err != nil {
			ui.LogFatal(err.Error())
		}

		ignorePath, err := cmd.Flags().GetString("gitignore")
		if err != nil {
			ui.LogFatal(err.Error())
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			ui.LogFatal(err.Error())
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			ui.LogFatal(err.Error())
		}

		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			ui.LogFatal(err.Error())
		}

		if path == "" {
			wd, err := os.Getwd()
			if err != nil {
				ui.LogFatal(err.Error())
			}
			path = wd
		}

		if ignorePath == "" {
			ignorePath = filepath.Join(path, ".gitignore")
		}

		gitIgnore, err := os.Open(ignorePath)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		ignorePatterns, err := scm.ParseIgnorePatterns(gitIgnore)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		im, err := issue.NewIssueManager(mode, annotation)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		_, err = im.Walk(path, ignorePatterns)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		issues := im.GetIssues()
		if len(issues) > 0 {
			fmt.Println(
				ui.SuccessTextStyle.Render(
					fmt.Sprintf(
						"\nFound %d (%s) issue annotations in your project.",
						len(issues),
						annotation,
					),
				),
			)
		} else {
			fmt.Println(ui.SecondaryTextStyle.Render(fmt.Sprintf("\nNo tags were located in your project using the annotation %s", annotation)))
		}

		if verbose {
			issue.PrintTagResults(issues, ui.DimTextStyle, ui.PrimaryTextStyle)
		} else {
			fmt.Println(
				ui.SecondaryTextStyle.Render(
					"Tip: run issue-summoner scan -v (verbose) for more details about the tag annotations that were found",
				),
			)
		}
	},
}

func init() {
	ScanCmd.Flags().StringP("path", "p", "", "path to local git repo")
	ScanCmd.Flags().StringP("gitignore", "g", "", "gitignore file path")
	ScanCmd.Flags().StringP("mode", "m", "pending", "i = issued | issues or p = pending")
	ScanCmd.Flags().BoolP("verbose", "v", false, "log information about each issue found")
	ScanCmd.Flags().StringP("annotation", "a", "@TODO", "Issue Annotation program will search for")
}

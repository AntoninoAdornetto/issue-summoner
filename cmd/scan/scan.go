/*
Copyright © 2024 Antonino Adornetto

The scan command processes each source file individually and searches for specific tags (actionable comments) that the user specifies.
It respects the `.gitignore` settings and ensures that any files designated as ignored are not scanned.
Finally, a detailed report is presented to the user about the tags that were found during the scan.
*/
package scan

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

type ScanManager struct{}

func (ScanManager) Open(fileName string) (*os.File, error) {
	return os.Open(fileName)
}

func (ScanManager) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans source code file(s) and searches for actionable comments",
	Long:  `Scans a local git respository for Tags (actionable comments) and reports findings to the console.`,
	Run: func(cmd *cobra.Command, _ []string) {
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

		if gitIgnorePath == "" {
			gitIgnorePath = filepath.Join(path, tag.GitIgnoreFile)
		}

		scanManager := ScanManager{}
		ignorePatterns, err := tag.ProcessIgnorePatterns(gitIgnorePath, scanManager)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		tagManager := tag.TagManager{
			TagName: tagName,
			Mode:    mode,
		}

		if err := tagManager.ValidateMode(mode); err != nil {
			ui.LogFatal(err.Error())
		}

		if mode == tag.PendingMode {
			pendingTagManager := tag.PendedTagManager{TagManager: tagManager}
			tags, err := tag.Walk(tag.WalkParams{
				Root:           path,
				TagManager:     &pendingTagManager,
				FileOperator:   scanManager,
				IgnorePatterns: ignorePatterns,
			})

			if err != nil {
				ui.LogFatal(fmt.Errorf("Failed to scan your project.\n%s", err).Error())
			}

			if verbose {
				tag.PrintTagResults(tags, ui.DimTextStyle, ui.PrimaryTextStyle)
			}

			if len(tags) > 0 {
				fmt.Println(
					ui.SuccessTextStyle.Render(
						fmt.Sprintf(
							"\nScan complete, found %d tags in your project that are using the annotation %s. See report command for submiting as issues to SCM \n",
							len(tags),
							tagName,
						),
					),
				)
				if !verbose {
					fmt.Println(
						ui.SecondaryTextStyle.Render(
							"Pass -v (verbose) flag for more details about the annotations found",
						),
					)
				}
			} else {
				fmt.Println(ui.SecondaryTextStyle.Render(fmt.Sprintf("\nNo tags were located in your project using the annotation %s", tagName)))
				fmt.Println(ui.NoteTextStyle.Render("\nRun issue-summoner scan --help to see usuage"))
			}
		}
	},
}

func init() {
	ScanCmd.Flags().StringP("path", "p", "", "Path to your local git project.")
	ScanCmd.Flags().StringP("tag", "t", "@TODO", "Actionable comment tag to search for.")
	ScanCmd.Flags().StringP("mode", "m", "P", "Mode: 'I' (Issued) or 'P' (Pending).")
	ScanCmd.Flags().StringP("gitignorePath", "g", "", "Path to .gitignore file.")
	ScanCmd.Flags().
		BoolP("verbose", "v", false, "Displays all information for each Tag annotation that is located")
}

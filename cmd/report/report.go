/*
Copyright © 2024 AntoninoAdornetto
*/
package report

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/tag"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

type ReportManager struct{}

func (ReportManager) Open(fileName string) (*os.File, error) {
	return os.Open(fileName)
}

func (ReportManager) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

// reportCmd represents the report command
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "report issues to a source code management system",
	Long: `Scans source code files for Tag annotations and reports them
	to a source code managment system of your choosing`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'path' flag\n%s", err).Error())
		}

		gitIgnorePath, err := cmd.Flags().GetString("gitignorePath")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'gitignorePath' flag\n%s", err).Error())
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

		scanManager := ReportManager{}
		ignorePatterns, err := tag.ProcessIgnorePatterns(gitIgnorePath, scanManager)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		tagManager := tag.TagManager{TagName: tagName}

		pendingTagManager := tag.PendedTagManager{TagManager: tagManager}
		_, err = tag.Walk(tag.WalkParams{
			Root:           path,
			TagManager:     &pendingTagManager,
			FileOperator:   scanManager,
			IgnorePatterns: ignorePatterns,
		})

		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to scan your project.\n%s", err).Error())
		}

		gc := scm.GitConfig{}

		err = gc.User()
		if err != nil {
			ui.LogFatal(
				fmt.Errorf("Failed to retrieve user.name from your global git config. See `git config global --help`", err).
					Error(),
			)
		}

		err = gc.RepoName()
		if err != nil {
			ui.LogFatal(
				fmt.Errorf("Failed to retrieve git remote origin url. %s", err).
					Error(),
			)
		}
	},
}

func init() {
	ReportCmd.Flags().StringP("path", "p", "", "Path to your local git project.")
	ReportCmd.Flags().StringP("tag", "t", "@TODO", "Actionable comment tag to search for.")
	ReportCmd.Flags().StringP("gitignorePath", "g", "", "Path to .gitignore file.")
}

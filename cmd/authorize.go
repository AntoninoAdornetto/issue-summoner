/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/git"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Create access tokens for the source code management platform you want to use for issue creation",
	Long: `Authorize command is used to create tokens that will allow issue summoner to report issues on your
	behalf.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger(cmd)

		srcCodeManager, err := cmd.Flags().GetString("scm")
		if err != nil {
			logger.Fatal(err.Error())
		}

		wd, err := os.Getwd()
		if err != nil {
			logger.Fatal(err.Error())
		}

		repo, err := git.NewRepository(wd)
		if err != nil {
			logger.Fatal(err.Error())
		}

		gitManager, err := git.NewGitManager(srcCodeManager, repo)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if gitManager.IsAuthorized() {
			logger.Warning(fmt.Sprintf("You are authorized for %s already", srcCodeManager))
			logger.PrintStdout("Do you want to create a new access token? (y/n): ")

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			proceed := scanner.Text()
			fmt.Printf("\n\n")

			if proceed != "y" {
				logger.Fatal("Authorization process aborted")
			}
		}

		spinner := tea.NewProgram(
			ui.InitSpinner(fmt.Sprintf("Pending authorization for %s", srcCodeManager)),
		)

		go func() {
			if _, err := spinner.Run(); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		if err := gitManager.Authorize(); err != nil {
			if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
				logger.Fatal(releaseErr.Error())
			}
			logger.Fatal(err.Error())
		}

		logger.Success(fmt.Sprintf("Authorization for %s succeeded!", srcCodeManager))
		if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
			logger.Fatal(releaseErr.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().StringP(flag_scm, shortflag_scm, git.GITHUB, flag_desc_scm)
	authorizeCmd.Flags().BoolP(flag_debug, shortflag_debug, false, flag_desc_debug)
}

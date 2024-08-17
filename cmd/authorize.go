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
	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// We only support GitHub, at the moment, but eventually I want to support all that are contained
// in the `allowedPlatforms` slice.
var allowedPlatforms = []string{scm.GITHUB, scm.GITLAB, scm.BITBUCKET}

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Create access tokens for the source code management platform you want to use for issue creation",
	Long: `Access tokens can be created for multiple source code management platforms (github, gitlab, bitbucket). This allows
Issue Summoner to submit issues to a specified source code management platform on your behalf`,
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

		if gitManager.Authenticated() {
			logger.Warning(fmt.Sprintf("You are authorized for %s already", srcCodeManager))
			logger.PrintStdout("Do you want to create a new access token? (y/n): ")

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			proceed := scanner.Text()

			fmt.Printf("\n")
			switch proceed {
			case "y", "yes":
				break
			default:
				logger.Fatal("Authorization process aborted")
			}
		}

		spinner := tea.NewProgram(
			spinner.InitialModelNew(fmt.Sprintf("Pending %s authorization", srcCodeManager)),
		)

		defer func() {
			if err := spinner.ReleaseTerminal(); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := spinner.Run(); err != nil {
				logger.Fatal(err.Error())
			}
		}()

		if err := gitManager.Authorize(); err != nil {
			if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
				logger.Fatal(err.Error())
			}
		}

		logger.Success(fmt.Sprintf("Authorization for %s succeeded!", srcCodeManager))
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().StringP(flag_scm, shortflag_scm, git.Github, flag_desc_scm)
	authorizeCmd.Flags().BoolP(flag_debug, shortflag_debug, false, flag_desc_debug)
}

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
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Create access tokens for the source code hosting platform you want to use for issue creation",
	Long: `Access tokens can be created for multiple source code hosting platforms (github, gitlab, bitbucket). This allows
Issue Summoner to submit issues to a specified source code hosting platform on your behalf`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := getLogger(cmd)

		srcCodeHost, err := cmd.Flags().GetString(flag_sch)
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

		gitManager, err := git.NewGitManager(srcCodeHost, repo)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if gitManager.Authenticated() {
			logger.Warning(fmt.Sprintf("You are authorized for %s already", srcCodeHost))
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
			ui.InitSpinner(fmt.Sprintf("Pending %s authorization", srcCodeHost)),
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

		logger.Success(fmt.Sprintf("Authorization for %s succeeded!", srcCodeHost))
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().StringP(flag_sch, shortflag_sch, git.Github, flag_desc_sch)
	authorizeCmd.Flags().BoolP(flag_debug, shortflag_debug, false, flag_desc_debug)
}

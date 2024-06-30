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
		srcCodeManager, err := cmd.Flags().GetString("scm")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'scm' flag\n%s", err).Error())
		}

		wd, err := os.Getwd()
		if err != nil {
			ui.LogFatal(err.Error())
		}

		repo, err := git.NewRepository(wd)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		gitManager, err := git.NewGitManager(srcCodeManager, repo)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		if gitManager.IsAuthorized() {
			fmt.Println(
				ui.PrimaryTextStyle.Render(
					fmt.Sprintf(
						"you are authorized for %s for already. Do you want to create a new access token?",
						srcCodeManager,
					),
				),
			)
			fmt.Print(
				ui.PrimaryTextStyle.Italic(true).
					Render("Type y to create a new token or type n to cancel the request: "),
			)

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			proceed := scanner.Text()
			fmt.Printf("\n\n")
			if proceed != "y" {
				ui.LogFatal("Authorization process aborted")
			}
		}

		spinner := &tea.Program{}
		go func() {
			spinner = tea.NewProgram(ui.InitialModelNew("Pending Authorization..."))
			if _, err := spinner.Run(); err != nil {
				ui.LogFatal(err.Error())
			}
		}()

		if err := gitManager.Authorize(); err != nil {
			ui.LogFatal(err.Error())
		}

		fmt.Println(
			ui.SuccessTextStyle.Render(
				fmt.Sprintf("Authorization for %s succeeded!", srcCodeManager),
			),
		)
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().StringP(flag_scm, shortflag_scm, git.GITHUB, flag_desc_scm)
}

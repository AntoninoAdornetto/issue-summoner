/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/AntoninoAdornetto/issue-summoner/v2/git"
	"github.com/spf13/cobra"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Create access tokens for the source code management platform you want to use for issue creation",
	Long: `Authorize will help the program create issues for both public and private
	repositories. Depending on the source code management platform you would like to authorize,
	we will need to verify your device with scopes that give the program access to opening new 
	issues. For example, when you Authorize with GitHub, we will need to create an access token with 
	repo scopes to grant read/write access to code, and issues. 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		sourceCodeManager, err := cmd.Flags().GetString("scm")
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

		gitManager, err := git.NewGitManager(sourceCodeManager, repo)
		if err != nil {
			ui.LogFatal(err.Error())
		}

		if gitManager.IsAuthorized() {
			fmt.Println(
				ui.PrimaryTextStyle.Render(
					fmt.Sprintf(
						"you are authorized for %s for already. Do you want to create a new access token?",
						sourceCodeManager,
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

		if err := gitManager.Authorize(); err != nil {
			ui.LogFatal(err.Error())
		}

		fmt.Println(
			ui.SuccessTextStyle.Render(
				fmt.Sprintf("Authorization for %s succeeded!", sourceCodeManager),
			),
		)
	},
}

func init() {
	rootCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().StringP(flag_scm, shortflag_scm, scm.GITHUB, flag_desc_scm)
}

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
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// We only support GitHub, at the moment, but eventually I want to support all that are contained
// in the `allowedPlatforms` slice.
var allowedPlatforms = []string{scm.GH, scm.GL, scm.BB}

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
		sourceCodeManager, err := cmd.Flags().GetString("scm")
		if err != nil {
			ui.LogFatal(fmt.Errorf("Failed to read 'scm' flag\n%s", err).Error())
		}

		found := false
		for _, v := range allowedPlatforms {
			if found {
				break
			}
			found = v == sourceCodeManager
		}

		if !found {
			ui.LogFatal(
				fmt.Sprintf(
					"%s is an unsupported source code management platform.\n",
					sourceCodeManager,
				),
			)
		}

		hasAccess, err := scm.CheckForAccess(sourceCodeManager)
		if err != nil && !os.IsNotExist(err) {
			ui.LogFatal(err.Error())
		}

		if hasAccess {
			fmt.Println(
				ui.PrimaryTextStyle.Render(
					fmt.Sprintf(
						"Looks like you are authorized for %s's platform already. Do you want to create a new access token?",
						sourceCodeManager,
					),
				),
			)

			fmt.Print(
				ui.PrimaryTextStyle.Italic(true).
					Render("Type 'y' to continue or type 'n' to cancel the request: "),
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

		gitManager := scm.GetGitConfig(sourceCodeManager)
		err = gitManager.Authorize()
		if err != nil {
			if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
				ui.ErrorTextStyle.Render("Error releasing terminal\n%s", releaseErr.Error())
			}
			ui.LogFatal(fmt.Errorf("Authorization failed.\n%s", err).Error())
		}

		err = spinner.ReleaseTerminal()
		if err != nil {
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
	authorizeCmd.Flags().StringP(flag_scm, shortflag_scm, scm.GH, flag_desc_scm)
}

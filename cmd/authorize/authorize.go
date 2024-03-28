/*
Copyright © 2024 AntoninoAdornetto
*/
package authorize

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

// We only support GitHub, at the moment, but eventually I want to support all that are contained
// in the below slice.
var allowedPlatforms = []string{scm.GH, scm.GL, scm.BB}

func init() {
	AuthorizeCmd.Flags().
		StringP("scm", "s", scm.GH, "What source code manager platform would you like to Authorize?")
}

// authorizeCmd represents the authorize command
var AuthorizeCmd = &cobra.Command{
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

		gitManager := scm.GetGitConfig(sourceCodeManager)

		hasAccess, err := scm.CheckForAccess(sourceCodeManager)
		if err != nil && !os.IsNotExist(err) {
			ui.LogFatal(err.Error())
		}

		if hasAccess {
			fmt.Println(
				ui.NoteTextStyle.Render(
					fmt.Sprintf(
						"Looks like you are authorized for %s's platform already. Do you want to create a new access token?",
						sourceCodeManager,
					),
				),
			)

			fmt.Print(
				ui.NoteTextStyle.Italic(true).Render("Type y to continue or n to cancel: "),
			)

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			proceed := scanner.Text()
			if proceed != "y" {
				ui.LogFatal("Aborted")
			}
		}

		fmt.Println(
			ui.SecondaryTextStyle.Render(
				fmt.Sprintf(
					"\nYou will be prompted to complete a few steps to authorize Issue Summoner for %s's platform.\nThis will allow us to open issues on your behalf",
					sourceCodeManager,
				),
			),
		)

		time.Sleep(time.Second * 2)

		err = gitManager.Authorize()
		if err != nil {
			ui.LogFatal(fmt.Errorf("Authorization failed.\n%s", err).Error())
		}

		fmt.Println(
			ui.SuccessTextStyle.Render(
				fmt.Sprintf("Authorization for %s succeeded!\n", sourceCodeManager),
			),
		)
	},
}

/*
Copyright © 2024 Antonino Adornetto
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

const Logo = `
.___                                _________                                                 
|   | ______ ________ __   ____    /   _____/__ __  _____   _____   ____   ____   ___________ 
|   |/  ___//  ___/  |  \_/ __ \   \_____  \|  |  \/     \ /     \ /  _ \ /    \_/ __ \_  __ \
|   |\___ \ \___ \|  |  /\  ___/   /        \  |  /  Y Y  \  Y Y  (  <_> )   |  \  ___/|  | \/
|___/____  >____  >____/  \___  > /_______  /____/|__|_|  /__|_|  /\____/|___|  /\___  >__|   
         \/     \/            \/          \/            \/      \/            \/     \/       
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "issue-summoner",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func subCommands() {
	rootCmd.AddCommand(AuthorizeCmd)
	rootCmd.AddCommand(ScanCmd)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.issue-summoner.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	subCommands()
	fmt.Println(ui.AccentTextStyle.Render(Logo))
}

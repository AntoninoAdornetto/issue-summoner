/*
Copyright © 2024 AntoninoAdornetto
*/
package cmd

import (
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/AntoninoAdornetto/issue-summoner/v2/git"
	"github.com/spf13/cobra"
)

// scan2Cmd represents the scan2 command
var scan2Cmd = &cobra.Command{
	Use:   "scan2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			ui.LogFatal(err.Error())
		}

		if path == "" {
			path, err = os.Getwd()
			if err != nil {
				ui.LogFatal(err.Error())
			}
		}

		_, err = git.NewRepository(path)
		if err != nil {
			ui.LogFatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(scan2Cmd)
	scan2Cmd.Flags().StringP("path", "p", "", "path to your working directory")
}

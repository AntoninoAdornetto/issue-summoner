/*
Copyright © 2024 Antonino Adornetto
*/

package scan

import (
	"github.com/spf13/cobra"
)

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Processes each source file individually and searches for specific tags",
	Long: `Processes each source file individually and searches for specific tags. 

It respects the .gitignore settings and ensures that any files designated as ignored are not scanned. 
Finally, a detailed report is presented to the user about the tags that were found during the scan.`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

func init() {}

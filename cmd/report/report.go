/*
Copyright © 2024 AntoninoAdornetto
*/
package report

import (
	"fmt"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/scm"
	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "report issues to a source code management system",
	Long: `Scans source code files for Tag annotations and reports them
	to a source code managment system of your choosing`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
}

package cmd

import (
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/common"
	"github.com/spf13/cobra"
)

const (
	err_unauthorized     = "Please run `issue-summoner authorize` and complete the authorization process. This will allow us to submit issues on your behalf."
	flag_annotation      = "annotation"
	flag_debug           = "debug"
	flag_desc_annotation = "The annotation to search for (@TODO:, @FIXME, etc)"
	flag_desc_debug      = "Log the stack trace when errors occur"
	flag_desc_mode       = "scan: searches for annotations denoted with the --annotation flag. purge: checks status of reported issues and removes comments"
	flag_desc_path       = "the path to your local git repository"
	flag_desc_sch        = "The source code hosting platform you would like to use. Such as, github, gitlab, or bitbucket"
	flag_desc_verbose    = "log detailed information about each issue annotation that is located during the scan"
	flag_mode            = "mode"
	flag_path            = "path"
	flag_sch             = "sch"
	flag_verbose         = "verbose"
	found_issues         = "Number of issues found: "
	issue_template_path  = "./templates/issue.tmpl"
	no_issues            = "No issues were found in your project using the annotation: "
	select_issues        = "Select the issues you wish to report"
	shortflag_annotation = "a"
	shortflag_debug      = "d"
	shortflag_mode       = "m"
	shortflag_path       = "p"
	shortflag_sch        = "s"
	shortflag_verbose    = "v"
	tip_verbose          = "run issue-summoner scan -v (verbose) for more details about the tag annotations that were found"
)

func getCommonFlags(cmd *cobra.Command) (annotation string, path string) {
	var err error
	annotation, err = cmd.Flags().GetString(flag_annotation)
	if err != nil {
		cobra.CheckErr(err)
	}

	path, err = cmd.Flags().GetString(flag_path)
	if err != nil {
		cobra.CheckErr(err)
	}

	if path == "" {
		wd, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(err)
		}
		path = wd
	}

	return annotation, path
}

func getLogger(cmd *cobra.Command) *common.Logger {
	debugIndicator, err := cmd.Flags().GetBool(flag_debug)
	if err != nil {
		cobra.CheckErr(err)
	}

	return common.NewLogger(debugIndicator)
}

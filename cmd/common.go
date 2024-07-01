package cmd

import (
	"os"

	"github.com/AntoninoAdornetto/issue-summoner/pkg/ui"
	"github.com/spf13/cobra"
)

const (
	err_unauthorized     = "Please run `issue-summoner authorize` and complete the authorization process. This will allow us to submit issues on your behalf."
	no_issues            = "No issues were found in your project using the annotation: "
	found_issues         = "Number of issues found: "
	select_issues        = "Select the issues you wish to report"
	issue_template_path  = "./templates/issue.tmpl"
	tip_verbose          = "run issue-summoner scan -v (verbose) for more details about the tag annotations that were found"
	flag_path            = "path"
	flag_mode            = "mode"
	flag_scm             = "scm"
	flag_debug           = "debug"
	flag_verbose         = "verbose"
	flag_annotation      = "annotation"
	shortflag_path       = "p"
	shortflag_scm        = "s"
	shortflag_mode       = "m"
	shortflag_verbose    = "v"
	shortflag_annotation = "a"
	shortflag_debug      = "d"
	flag_desc_path       = "the path to your local git repository"
	flag_desc_scm        = "The source code management platform you would like to use. Such as, github, gitlab, or bitbucket"
	flag_desc_mode       = "'processed' is for issues that have already been pushed to a scm. 'pending' is for issues that have not yet been published"
	flag_desc_verbose    = "log detailed information about each issue annotation that is located during the scan"
	flag_desc_annotation = "The issue annotation to search for. Example: @TODO:"
	flag_desc_debug      = "Log the stack trace when errors occur"
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

func getLogger(cmd *cobra.Command) (logger *ui.Logger) {
	debugIndicator, err := cmd.Flags().GetBool(flag_debug)
	if err != nil {
		cobra.CheckErr(err)
	}

	return ui.NewLogger(debugIndicator)
}

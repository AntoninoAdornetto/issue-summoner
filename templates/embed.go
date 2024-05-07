/*
Using embed will allow us to directly embed the contents of the issue.tmpl
template file into a variable (issueTemplate). When the program is compiled,
go can read the contents of the file. This process is exactly what we need
to publish git issues from any directory in the file system. Additionally,
the template file is small (4KB) and we only expect to have the 1 file.
*/
package templates

import (
	"embed"
	"text/template"
)

var (
	//go:embed issue.tmpl
	issueTemplate embed.FS
)

func LoadIssueTemplate() (*template.Template, error) {
	tmpl, err := template.New("issue.tmpl").ParseFS(issueTemplate, "issue.tmpl")
	return tmpl, err
}

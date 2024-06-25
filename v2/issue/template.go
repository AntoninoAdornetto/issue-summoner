/*
The template in this file is used when executing the 'report' command.
For each code issue that is selected, the below template is executed against
the given issue. The result is a formatted markdown issue that is published to
the source code managmenet platform flag that is passed into the report command.
*/
package issue

import "text/template"

var (
	issue_template_markdown = `### Description
{{ .Description }}

### Location

***File name:*** {{ .FileName }} ***Line number:*** {{ .LineNumber }}

### Environment

{{ .Environment }}

### Generated with :heart:

created by [issue-summoner](https://github.com/AntoninoAdornetto/issue-summoner)
	`
)

func generateIssueTemplate() (*template.Template, error) {
	return template.New("").Parse(issue_template_markdown)
}

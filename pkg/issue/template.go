/*
THE TEMPLATE IN THIS FILE IS USED WHEN EXECUTING THE 'REPORT' COMMAND.
FOR EACH CODE ISSUE THAT IS SELECTED, THE BELOW TEMPLATE IS EXECUTED AGAINST
THE GIVEN ISSUE. THE RESULT IS A FORMATTED MARKDOWN ISSUE THAT IS PUBLISHED TO
THE SOURCE CODE MANAGMENET PLATFORM FLAG THAT IS PASSED INTO THE REPORT COMMAND.
*/
package issue

import "text/template"

var (
	issue_template_markdown = `### Description
{{ .Description }}

### Location

***File name:*** {{ .FileName }} ***Line number:*** {{ .LineNumber }}

### Environment

{{ .OS }}

### Generated with :heart:

created by [issue-summoner](https://github.com/AntoninoAdornetto/issue-summoner)
	`
)

func generateIssueTemplate() (*template.Template, error) {
	return template.New("").Parse(issue_template_markdown)
}

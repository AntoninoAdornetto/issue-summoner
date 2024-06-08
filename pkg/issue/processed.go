package issue

/*
* @TODO implement the processed issues logic
* The functions implementations should satisfy the IssueManager interface and allow the
* user of the program to walk the project, and locate comments/issues that have
* been reported. This will be helpful in later functionality where we can remove
* comments once issues have been marked as resolved.
 */

type ProcessedIssue struct {
	Annotation string
	Issues     []Issue
}

func (pi *ProcessedIssue) Walk(root string) (int, error) {
	n := 0
	return n, nil
}

func (pi *ProcessedIssue) Scan(src []byte, path string) error {
	return nil
}

func (pi *ProcessedIssue) GetIssues() []Issue {
	return pi.Issues
}

func (pi *ProcessedIssue) WriteIssueID(id int64, issueIndex int) error {
	return nil
}

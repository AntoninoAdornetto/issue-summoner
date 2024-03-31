package issue

type ProcessedIssue struct {
	Issues []Issue
}

func (pi *ProcessedIssue) Scan() ([]Issue, error) {
	return nil, nil
}

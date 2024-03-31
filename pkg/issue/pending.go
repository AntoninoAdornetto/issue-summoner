package issue

type PendingIssue struct {
	Issues []Issue
}

func (pi *PendingIssue) Scan() ([]Issue, error) {
	return nil, nil
}

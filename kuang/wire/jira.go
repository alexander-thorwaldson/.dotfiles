package wire

// JiraIssueListRequest is the input for listing Jira issues.
type JiraIssueListRequest struct {
	Project  string `json:"project,omitempty"`
	Query    string `json:"query,omitempty"`
	Assignee string `json:"assignee,omitempty"`
	Status   string `json:"status,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// JiraIssueViewRequest is the input for viewing a Jira issue.
type JiraIssueViewRequest struct {
	Issue string `json:"issue"`
}

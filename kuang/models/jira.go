package models

// -- Jira Issues --

// JiraUser represents a Jira user.
type JiraUser struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// JiraStatus represents a Jira issue status.
type JiraStatus struct {
	Name string `json:"name"`
}

// JiraPriority represents a Jira issue priority.
type JiraPriority struct {
	Name string `json:"name"`
}

// JiraIssueType represents a Jira issue type.
type JiraIssueType struct {
	Name string `json:"name"`
}

// JiraIssueSummary is a compact representation of a Jira issue.
type JiraIssueSummary struct {
	Key      string `json:"key"`
	Summary  string `json:"summary"`
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
	Priority string `json:"priority"`
	Type     string `json:"type"`
}

// JiraIssueListResult wraps a list of Jira issue summaries.
type JiraIssueListResult struct {
	Items []JiraIssueSummary `json:"items"`
}

// JiraComment represents a comment on a Jira issue.
type JiraComment struct {
	Author  string `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

// JiraIssueDetail is the full representation of a Jira issue.
type JiraIssueDetail struct {
	Key         string        `json:"key"`
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Assignee    string        `json:"assignee"`
	Reporter    string        `json:"reporter"`
	Priority    string        `json:"priority"`
	Type        string        `json:"type"`
	Labels      []string      `json:"labels"`
	Comments    []JiraComment `json:"comments"`
}

// -- Jira Projects --

// JiraProject represents a Jira project.
type JiraProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// JiraProjectListResult wraps a list of Jira projects.
type JiraProjectListResult struct {
	Items []JiraProject `json:"items"`
}

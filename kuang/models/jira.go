package models

// -- Jira Issues --

type JiraUser struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type JiraStatus struct {
	Name string `json:"name"`
}

type JiraPriority struct {
	Name string `json:"name"`
}

type JiraIssueType struct {
	Name string `json:"name"`
}

type JiraIssueSummary struct {
	Key      string        `json:"key"`
	Summary  string        `json:"summary"`
	Status   string        `json:"status"`
	Assignee string        `json:"assignee"`
	Priority string        `json:"priority"`
	Type     string        `json:"type"`
}

type JiraIssueListResult struct {
	Items []JiraIssueSummary `json:"items"`
}

type JiraComment struct {
	Author  string `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

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

type JiraProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type JiraProjectListResult struct {
	Items []JiraProject `json:"items"`
}

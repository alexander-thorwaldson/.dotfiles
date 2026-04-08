package models

// -- Pull Requests --

type PRAuthor struct {
	Login string `json:"login"`
}

type PRSummary struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	State  string   `json:"state"`
	Author PRAuthor `json:"author"`
	URL    string   `json:"url"`
}

type PRListResult struct {
	Items []PRSummary `json:"items"`
}

type PRReview struct {
	Author PRAuthor `json:"author"`
	State  string   `json:"state"`
	Body   string   `json:"body"`
}

type PRComment struct {
	Author PRAuthor `json:"author"`
	Body   string   `json:"body"`
}

type PRDetail struct {
	Number   int         `json:"number"`
	Title    string      `json:"title"`
	State    string      `json:"state"`
	Body     string      `json:"body"`
	Author   PRAuthor    `json:"author"`
	URL      string      `json:"url"`
	Reviews  []PRReview  `json:"reviews"`
	Comments []PRComment `json:"comments"`
}

// -- Issues --

type IssueSummary struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	State  string   `json:"state"`
	Author PRAuthor `json:"author"`
	URL    string   `json:"url"`
}

type IssueListResult struct {
	Items []IssueSummary `json:"items"`
}

type IssueComment struct {
	Author PRAuthor `json:"author"`
	Body   string   `json:"body"`
}

type IssueDetail struct {
	Number   int            `json:"number"`
	Title    string         `json:"title"`
	State    string         `json:"state"`
	Body     string         `json:"body"`
	Author   PRAuthor       `json:"author"`
	URL      string         `json:"url"`
	Comments []IssueComment `json:"comments"`
}

// -- Repos --

type BranchRef struct {
	Name string `json:"name"`
}

type RepoCount struct {
	TotalCount int `json:"totalCount"`
}

type RepoDetail struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	URL              string            `json:"url"`
	DefaultBranchRef BranchRef         `json:"defaultBranchRef"`
	Languages        []RepoLanguage    `json:"languages"`
	Issues           RepoCount         `json:"issues"`
	PullRequests     RepoCount         `json:"pullRequests"`
}

type RepoLanguage struct {
	Node struct {
		Name string `json:"name"`
	} `json:"node"`
}

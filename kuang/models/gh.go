// Package models defines the data structures returned by tool commands.
package models

// -- Pull Requests --

// PRAuthor represents the author of a pull request.
type PRAuthor struct {
	Login string `json:"login"`
}

// PRSummary is a compact representation of a pull request.
type PRSummary struct {
	Title  string   `json:"title"`
	State  string   `json:"state"`
	Author PRAuthor `json:"author"`
	URL    string   `json:"url"`
	Number int      `json:"number"`
}

// PRListResult wraps a list of pull request summaries.
type PRListResult struct {
	Items []PRSummary `json:"items"`
}

// PRReview represents a review left on a pull request.
type PRReview struct {
	Author PRAuthor `json:"author"`
	State  string   `json:"state"`
	Body   string   `json:"body"`
}

// PRComment represents a comment on a pull request.
type PRComment struct {
	Author PRAuthor `json:"author"`
	Body   string   `json:"body"`
}

// PRDetail is the full representation of a pull request.
type PRDetail struct {
	Title    string      `json:"title"`
	State    string      `json:"state"`
	Body     string      `json:"body"`
	Author   PRAuthor    `json:"author"`
	URL      string      `json:"url"`
	Reviews  []PRReview  `json:"reviews"`
	Comments []PRComment `json:"comments"`
	Number   int         `json:"number"`
}

// -- Issues --

// IssueSummary is a compact representation of an issue.
type IssueSummary struct {
	Title  string   `json:"title"`
	State  string   `json:"state"`
	Author PRAuthor `json:"author"`
	URL    string   `json:"url"`
	Number int      `json:"number"`
}

// IssueListResult wraps a list of issue summaries.
type IssueListResult struct {
	Items []IssueSummary `json:"items"`
}

// IssueComment represents a comment on an issue.
type IssueComment struct {
	Author PRAuthor `json:"author"`
	Body   string   `json:"body"`
}

// IssueDetail is the full representation of an issue.
type IssueDetail struct {
	Title    string         `json:"title"`
	State    string         `json:"state"`
	Body     string         `json:"body"`
	Author   PRAuthor       `json:"author"`
	URL      string         `json:"url"`
	Comments []IssueComment `json:"comments"`
	Number   int            `json:"number"`
}

// -- Repos --

// BranchRef holds a branch name reference.
type BranchRef struct {
	Name string `json:"name"`
}

// RepoCount holds a total count for a repository resource.
type RepoCount struct {
	TotalCount int `json:"totalCount"`
}

// RepoDetail is the full representation of a repository.
type RepoDetail struct {
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	URL              string         `json:"url"`
	DefaultBranchRef BranchRef      `json:"defaultBranchRef"`
	Languages        []RepoLanguage `json:"languages"`
	Issues           RepoCount      `json:"issues"`
	PullRequests     RepoCount      `json:"pullRequests"`
}

// RepoLanguage represents a language used in a repository.
type RepoLanguage struct {
	Node struct {
		Name string `json:"name"`
	} `json:"node"`
}

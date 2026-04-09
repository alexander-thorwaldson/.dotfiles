package wire

// -- Pull Requests --

// PRListRequest is the input for listing pull requests.
type PRListRequest struct {
	Repo  string `json:"repo"`
	State string `json:"state,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

// PRViewRequest is the input for viewing a pull request.
type PRViewRequest struct {
	Repo   string `json:"repo"`
	Number int    `json:"number"`
}

// PRCreateRequest is the input for creating a pull request.
type PRCreateRequest struct {
	Repo     string   `json:"repo"`
	Title    string   `json:"title"`
	Head     string   `json:"head"`
	Body     string   `json:"body,omitempty"`
	Base     string   `json:"base,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	Reviewer []string `json:"reviewers,omitempty"`
	Draft    bool     `json:"draft,omitempty"`
}

// -- Issues --

// IssueListRequest is the input for listing issues.
type IssueListRequest struct {
	Repo  string `json:"repo"`
	State string `json:"state,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

// IssueViewRequest is the input for viewing an issue.
type IssueViewRequest struct {
	Repo   string `json:"repo"`
	Number int    `json:"number"`
}

// -- Repo --

// RepoViewRequest is the input for viewing a repository.
type RepoViewRequest struct {
	Repo string `json:"repo"`
}

// -- Actions --

// RunListRequest is the input for listing workflow runs.
type RunListRequest struct {
	Repo     string `json:"repo"`
	Branch   string `json:"branch,omitempty"`
	Status   string `json:"status,omitempty"`
	Workflow string `json:"workflow,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// RunViewRequest is the input for viewing a workflow run.
type RunViewRequest struct {
	Repo  string `json:"repo"`
	RunID int    `json:"runId"`
}

// RunLogRequest is the input for viewing failed logs from a workflow run.
type RunLogRequest struct {
	Repo  string `json:"repo"`
	RunID int    `json:"runId"`
}

// PRChecksRequest is the input for viewing checks on a pull request.
type PRChecksRequest struct {
	Repo   string `json:"repo"`
	Number int    `json:"number"`
}

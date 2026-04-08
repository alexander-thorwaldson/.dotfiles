package models

import (
	"encoding/json"
	"testing"
)

func TestPRSummary_Unmarshal(t *testing.T) {
	raw := `[{"number":1,"title":"Fix bug","state":"OPEN","author":{"login":"alice"},"url":"https://github.com/o/r/pull/1"}]`
	var items []PRSummary
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	pr := items[0]
	if pr.Number != 1 {
		t.Errorf("expected number 1, got %d", pr.Number)
	}
	if pr.Title != "Fix bug" {
		t.Errorf("expected title %q, got %q", "Fix bug", pr.Title)
	}
	if pr.Author.Login != "alice" {
		t.Errorf("expected author %q, got %q", "alice", pr.Author.Login)
	}
}

func TestPRDetail_Unmarshal(t *testing.T) {
	raw := `{
		"number": 42,
		"title": "Add feature",
		"state": "OPEN",
		"body": "This adds a new feature",
		"author": {"login": "bob"},
		"url": "https://github.com/o/r/pull/42",
		"reviews": [{"author": {"login": "carol"}, "state": "APPROVED", "body": "LGTM"}],
		"comments": [{"author": {"login": "dave"}, "body": "Nice work"}]
	}`
	var pr PRDetail
	if err := json.Unmarshal([]byte(raw), &pr); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if pr.Number != 42 {
		t.Errorf("expected number 42, got %d", pr.Number)
	}
	if pr.Body != "This adds a new feature" {
		t.Errorf("unexpected body: %q", pr.Body)
	}
	if len(pr.Reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(pr.Reviews))
	}
	if pr.Reviews[0].State != "APPROVED" {
		t.Errorf("expected review state APPROVED, got %q", pr.Reviews[0].State)
	}
	if len(pr.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(pr.Comments))
	}
	if pr.Comments[0].Body != "Nice work" {
		t.Errorf("unexpected comment body: %q", pr.Comments[0].Body)
	}
}

func TestPRListResult_Unmarshal(t *testing.T) {
	raw := `{"items": [{"number": 1, "title": "PR", "state": "OPEN", "author": {"login": "x"}, "url": "u"}]}`
	var result PRListResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Items))
	}
}

func TestIssueSummary_Unmarshal(t *testing.T) {
	raw := `[{"number":10,"title":"Bug","state":"OPEN","author":{"login":"eve"},"url":"https://github.com/o/r/issues/10"}]`
	var items []IssueSummary
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Number != 10 {
		t.Errorf("expected number 10, got %d", items[0].Number)
	}
}

func TestIssueDetail_Unmarshal(t *testing.T) {
	raw := `{
		"number": 10,
		"title": "Bug",
		"state": "OPEN",
		"body": "Steps to reproduce",
		"author": {"login": "eve"},
		"url": "https://github.com/o/r/issues/10",
		"comments": [{"author": {"login": "frank"}, "body": "Can confirm"}]
	}`
	var issue IssueDetail
	if err := json.Unmarshal([]byte(raw), &issue); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if issue.Body != "Steps to reproduce" {
		t.Errorf("unexpected body: %q", issue.Body)
	}
	if len(issue.Comments) != 1 || issue.Comments[0].Author.Login != "frank" {
		t.Errorf("unexpected comments: %+v", issue.Comments)
	}
}

func TestRepoDetail_Unmarshal(t *testing.T) {
	raw := `{
		"name": "myrepo",
		"description": "A repo",
		"url": "https://github.com/o/myrepo",
		"defaultBranchRef": {"name": "main"},
		"languages": [{"node": {"name": "Go"}}],
		"issues": {"totalCount": 5},
		"pullRequests": {"totalCount": 3}
	}`
	var repo RepoDetail
	if err := json.Unmarshal([]byte(raw), &repo); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if repo.Name != "myrepo" {
		t.Errorf("expected name %q, got %q", "myrepo", repo.Name)
	}
	if repo.DefaultBranchRef.Name != "main" {
		t.Errorf("expected branch %q, got %q", "main", repo.DefaultBranchRef.Name)
	}
	if repo.Issues.TotalCount != 5 {
		t.Errorf("expected 5 issues, got %d", repo.Issues.TotalCount)
	}
	if repo.PullRequests.TotalCount != 3 {
		t.Errorf("expected 3 PRs, got %d", repo.PullRequests.TotalCount)
	}
	if len(repo.Languages) != 1 || repo.Languages[0].Node.Name != "Go" {
		t.Errorf("unexpected languages: %+v", repo.Languages)
	}
}

func TestPRSummary_Marshal(t *testing.T) {
	pr := PRSummary{
		Number: 1,
		Title:  "Test",
		State:  "OPEN",
		Author: PRAuthor{Login: "alice"},
		URL:    "https://example.com",
	}
	data, err := json.Marshal(pr)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var roundtrip PRSummary
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("roundtrip unmarshal error: %v", err)
	}
	if roundtrip != pr {
		t.Errorf("roundtrip mismatch: got %+v", roundtrip)
	}
}

func TestRunSummary_Unmarshal(t *testing.T) {
	raw := `[{"number":100,"databaseId":9999,"displayTitle":"CI","headBranch":"main","event":"push","status":"completed","conclusion":"success","workflowName":"CI","url":"https://github.com/o/r/actions/runs/9999","createdAt":"2026-04-07T10:00:00Z"}]`
	var items []RunSummary
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Conclusion != "success" {
		t.Errorf("expected conclusion success, got %s", items[0].Conclusion)
	}
	if items[0].DatabaseID != 9999 {
		t.Errorf("expected databaseId 9999, got %d", items[0].DatabaseID)
	}
}

func TestRunDetail_Unmarshal(t *testing.T) {
	raw := `{
		"number": 100,
		"attempt": 1,
		"displayTitle": "CI",
		"headBranch": "main",
		"headSha": "abc123",
		"event": "push",
		"status": "completed",
		"conclusion": "failure",
		"workflowName": "CI",
		"url": "https://github.com/o/r/actions/runs/100",
		"createdAt": "2026-04-07T10:00:00Z",
		"jobs": [{"name": "build", "status": "completed", "conclusion": "failure", "steps": [{"name": "checkout", "status": "completed", "conclusion": "success", "number": 1}, {"name": "test", "status": "completed", "conclusion": "failure", "number": 2}]}]
	}`
	var run RunDetail
	if err := json.Unmarshal([]byte(raw), &run); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if run.Conclusion != "failure" {
		t.Errorf("expected conclusion failure, got %s", run.Conclusion)
	}
	if len(run.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(run.Jobs))
	}
	if len(run.Jobs[0].Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(run.Jobs[0].Steps))
	}
	if run.Jobs[0].Steps[1].Conclusion != "failure" {
		t.Errorf("expected step 2 failure, got %s", run.Jobs[0].Steps[1].Conclusion)
	}
}

func TestPRCheck_Unmarshal(t *testing.T) {
	raw := `[{"name":"CI","state":"SUCCESS","bucket":"pass","description":"All checks passed","workflow":"CI","link":"https://github.com/o/r/actions/runs/1","event":"pull_request","startedAt":"2026-04-07T10:00:00Z","completedAt":"2026-04-07T10:05:00Z"}]`
	var items []PRCheck
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 check, got %d", len(items))
	}
	if items[0].Bucket != "pass" {
		t.Errorf("expected bucket pass, got %s", items[0].Bucket)
	}
	if items[0].State != "SUCCESS" {
		t.Errorf("expected state SUCCESS, got %s", items[0].State)
	}
}

func TestEmptyCollections(t *testing.T) {
	raw := `{"items": []}`
	var result PRListResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result.Items == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
}

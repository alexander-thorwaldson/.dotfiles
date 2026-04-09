package models

import (
	"encoding/json"
	"testing"
)

func TestJiraIssueSummary_Unmarshal(t *testing.T) {
	raw := `[{"key":"PROJ-1","summary":"Fix bug","status":"Open","assignee":"alice","priority":"High","type":"Bug"}]`
	var items []JiraIssueSummary
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Key != "PROJ-1" {
		t.Errorf("expected key PROJ-1, got %s", items[0].Key)
	}
	if items[0].Priority != "High" {
		t.Errorf("expected priority High, got %s", items[0].Priority)
	}
}

func TestJiraIssueDetail_Unmarshal(t *testing.T) {
	raw := `{
		"key": "PROJ-42",
		"summary": "Fix auth",
		"description": "Auth flow is broken",
		"status": "In Progress",
		"assignee": "bob",
		"reporter": "carol",
		"priority": "Critical",
		"type": "Bug",
		"labels": ["security", "urgent"],
		"comments": [
			{"author": "dave", "body": "On it", "created": "2026-04-01T10:00:00Z"}
		]
	}`
	var detail JiraIssueDetail
	if err := json.Unmarshal([]byte(raw), &detail); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if detail.Key != "PROJ-42" {
		t.Errorf("expected key PROJ-42, got %s", detail.Key)
	}
	if detail.Reporter != "carol" {
		t.Errorf("expected reporter carol, got %s", detail.Reporter)
	}
	if len(detail.Labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(detail.Labels))
	}
	if detail.Labels[0] != "security" {
		t.Errorf("expected first label security, got %s", detail.Labels[0])
	}
	if len(detail.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(detail.Comments))
	}
	if detail.Comments[0].Created != "2026-04-01T10:00:00Z" {
		t.Errorf("expected created timestamp, got %s", detail.Comments[0].Created)
	}
}

func TestJiraProject_Unmarshal(t *testing.T) {
	raw := `{"key": "PROJ", "name": "My Project"}`
	var project JiraProject
	if err := json.Unmarshal([]byte(raw), &project); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if project.Key != "PROJ" {
		t.Errorf("expected key PROJ, got %s", project.Key)
	}
	if project.Name != "My Project" {
		t.Errorf("expected name %q, got %q", "My Project", project.Name)
	}
}

func TestJiraIssueListResult_Empty(t *testing.T) {
	raw := `{"items": []}`
	var result JiraIssueListResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
}

func TestJiraIssueSummary_Marshal(t *testing.T) {
	issue := JiraIssueSummary{
		Key:      "TEST-1",
		Summary:  "Test issue",
		Status:   "Open",
		Assignee: "alice",
		Priority: "Medium",
		Type:     "Task",
	}
	data, err := json.Marshal(issue)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var roundtrip JiraIssueSummary
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("roundtrip error: %v", err)
	}
	if roundtrip != issue {
		t.Errorf("roundtrip mismatch: got %+v", roundtrip)
	}
}

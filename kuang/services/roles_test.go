package services

import (
	"testing"
)

func TestExpandRole_Engineer(t *testing.T) {
	perms := ExpandRole(RoleEngineer)
	if len(perms) == 0 {
		t.Fatal("expected permissions for engineer role")
	}

	expected := []string{
		"gh_pr_list", "gh_pr_view", "gh_pr_create",
		"gh_pr_checks", "gh_run_list",
		"dev_start", "dev_stop", "dev_status", "dev_log",
		"pnpm_run",
	}
	for _, e := range expected {
		if !containsString(perms, e) {
			t.Errorf("engineer missing permission %q", e)
		}
	}
}

func TestExpandRole_Architect(t *testing.T) {
	perms := ExpandRole(RoleArchitect)

	// Architect should have read access.
	if !containsString(perms, "gh_pr_view") {
		t.Error("architect should have gh_pr_view")
	}
	if !containsString(perms, "jira_issue_list") {
		t.Error("architect should have jira_issue_list")
	}

	// Architect should NOT have write access.
	if containsString(perms, "gh_pr_create") {
		t.Error("architect should not have gh_pr_create")
	}
	if containsString(perms, "dev_start") {
		t.Error("architect should not have dev_start")
	}
	if containsString(perms, "pnpm_run") {
		t.Error("architect should not have pnpm_run")
	}
}

func TestExpandRole_PM(t *testing.T) {
	perms := ExpandRole(RolePM)

	if !containsString(perms, "jira_issue_list") {
		t.Error("pm should have jira_issue_list")
	}
	if !containsString(perms, "dev_status") {
		t.Error("pm should have dev_status")
	}

	// PM should NOT have dev management beyond status.
	if containsString(perms, "dev_start") {
		t.Error("pm should not have dev_start")
	}
	if containsString(perms, "gh_pr_create") {
		t.Error("pm should not have gh_pr_create")
	}
}

func TestExpandRole_Admin(t *testing.T) {
	perms := ExpandRole(RoleAdmin)
	all := AllToolPermissions()

	for _, tool := range all {
		if !containsString(perms, tool) {
			t.Errorf("admin missing permission %q", tool)
		}
	}
}

func TestExpandRole_Unknown(t *testing.T) {
	perms := ExpandRole(Role("hacker"))
	if perms != nil {
		t.Errorf("expected nil for unknown role, got %v", perms)
	}
}

func TestExpandRole_NoDuplicates(t *testing.T) {
	for _, role := range []Role{RoleEngineer, RoleDesigner, RoleArchitect, RolePM, RoleSecurity, RoleExecutive, RoleAdmin} {
		perms := ExpandRole(role)
		seen := make(map[string]bool)
		for _, p := range perms {
			if seen[p] {
				t.Errorf("role %q has duplicate permission %q", role, p)
			}
			seen[p] = true
		}
	}
}

func TestRoleHasPermission(t *testing.T) {
	if !RoleHasPermission(RoleEngineer, "gh_pr_create") {
		t.Error("engineer should have gh_pr_create")
	}
	if RoleHasPermission(RoleArchitect, "gh_pr_create") {
		t.Error("architect should not have gh_pr_create")
	}
	if !RoleHasPermission(RoleSecurity, "gh_run_log") {
		t.Error("security should have gh_run_log")
	}
	if RoleHasPermission(Role("nobody"), "gh_pr_list") {
		t.Error("unknown role should have no permissions")
	}
}

func TestAllToolPermissions(t *testing.T) {
	all := AllToolPermissions()
	if len(all) == 0 {
		t.Fatal("expected tool permissions")
	}

	// Every tool registered in groupTools should appear.
	for group, tools := range groupTools {
		for _, tool := range tools {
			if !containsString(all, tool) {
				t.Errorf("AllToolPermissions missing %q from group %q", tool, group)
			}
		}
	}
}

func TestAllToolPermissions_NoDuplicates(t *testing.T) {
	all := AllToolPermissions()
	seen := make(map[string]bool)
	for _, p := range all {
		if seen[p] {
			t.Errorf("AllToolPermissions has duplicate %q", p)
		}
		seen[p] = true
	}
}

func TestDevStatus_InBothGroups(t *testing.T) {
	// dev_status is in both DevManage and DevStatus groups.
	// PM gets it via DevStatus, Engineer gets it via DevManage.
	if !RoleHasPermission(RolePM, "dev_status") {
		t.Error("pm should have dev_status via DevStatus group")
	}
	if !RoleHasPermission(RoleEngineer, "dev_status") {
		t.Error("engineer should have dev_status via DevManage group")
	}
}

func containsString(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

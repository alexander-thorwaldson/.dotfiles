package handlers

// Role is a named authorization role carried in a certificate's OU field.
type Role string

// Supported roles.
const (
	RoleEngineer  Role = "engineer"
	RoleDesigner  Role = "designer"
	RoleArchitect Role = "architect"
	RolePM        Role = "pm"
	RoleSecurity  Role = "security"
	RoleExecutive Role = "executive"
	RoleAdmin     Role = "admin"
)

// PermissionGroup is a named set of related tool permissions.
type PermissionGroup string

// Supported permission groups.
const (
	GroupReadGH    PermissionGroup = "read:gh"
	GroupWriteGH   PermissionGroup = "write:gh"
	GroupCIGH      PermissionGroup = "ci:gh"
	GroupReadJira  PermissionGroup = "read:jira"
	GroupDevManage PermissionGroup = "dev:manage"
	GroupDevStatus PermissionGroup = "dev:status"
	GroupDevRun    PermissionGroup = "dev:run"
)

// groupTools maps each permission group to the individual tool permissions it grants.
var groupTools = map[PermissionGroup][]string{
	GroupReadGH:    {"gh_pr_list", "gh_pr_view", "gh_issue_list", "gh_issue_view", "gh_repo_view"},
	GroupWriteGH:   {"gh_pr_create"},
	GroupCIGH:      {"gh_pr_checks", "gh_run_list", "gh_run_view", "gh_run_log"},
	GroupReadJira:  {"jira_issue_list", "jira_issue_view"},
	GroupDevManage: {"dev_start", "dev_stop", "dev_status", "dev_log"},
	GroupDevStatus: {"dev_status"},
	GroupDevRun:    {"pnpm_run"},
}

// roleGroups maps each role to its granted permission groups.
var roleGroups = map[Role][]PermissionGroup{
	RoleEngineer:  {GroupReadGH, GroupWriteGH, GroupCIGH, GroupDevManage, GroupDevRun},
	RoleDesigner:  {GroupReadGH, GroupCIGH, GroupReadJira, GroupDevRun},
	RoleArchitect: {GroupReadGH, GroupCIGH, GroupReadJira},
	RolePM:        {GroupReadGH, GroupReadJira, GroupDevStatus},
	RoleSecurity:  {GroupReadGH, GroupCIGH, GroupDevRun},
	RoleExecutive: {GroupReadGH, GroupCIGH, GroupReadJira, GroupDevStatus},
	RoleAdmin:     {GroupReadGH, GroupWriteGH, GroupCIGH, GroupReadJira, GroupDevManage, GroupDevRun},
}

// ExpandRole resolves a role into its full set of individual tool permissions.
func ExpandRole(role Role) []string {
	groups, ok := roleGroups[role]
	if !ok {
		return nil
	}

	seen := make(map[string]bool)
	var perms []string
	for _, g := range groups {
		for _, tool := range groupTools[g] {
			if !seen[tool] {
				seen[tool] = true
				perms = append(perms, tool)
			}
		}
	}
	return perms
}

// AllToolPermissions returns every tool permission across all groups.
func AllToolPermissions() []string {
	seen := make(map[string]bool)
	var perms []string
	for _, tools := range groupTools {
		for _, tool := range tools {
			if !seen[tool] {
				seen[tool] = true
				perms = append(perms, tool)
			}
		}
	}
	return perms
}

// RoleHasPermission checks if a role grants a specific tool permission.
func RoleHasPermission(role Role, permission string) bool {
	for _, perm := range ExpandRole(role) {
		if perm == permission {
			return true
		}
	}
	return false
}

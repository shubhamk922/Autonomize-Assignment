package jira

import (
	"fmt"
	"strings"

	"example.com/team-monitoring/domain"
)

func BuildJQL(q domain.JiraQuery) string {
	var clauses []string

	if q.Project != "" {
		clauses = append(clauses, fmt.Sprintf(`project = "%s"`, q.Project))
	}
	if q.Assignee != "" {
		clauses = append(clauses, fmt.Sprintf(`assignee = %s`, q.Assignee))
	} else {
		clauses = append(clauses, "assignee = currentUser()")
	}
	if q.Status != "" {
		clauses = append(clauses, fmt.Sprintf(`status = "%s"`, q.Status))
	}

	if len(clauses) == 0 {
		return "ORDER BY created DESC"
	}

	return strings.Join(clauses, " AND ")
}

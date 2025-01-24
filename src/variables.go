package main

import "database/sql"

var (
	cmd            string
	desc           string
	proj           string = ""
	projPassed     bool   = false
	priority       uint64 = 0
	id             uint64
	tagsAdd        []string     = []string{}
	tagsRemove     []string     = []string{}
	dueDate        sql.NullTime = sql.NullTime{}
	dueDatePassed  bool         = false
	doneDate       sql.NullTime = sql.NullTime{}
	doneDatePassed bool         = false

	verbs = []string{
		"add",
		"complete",
		"delete",
		"done",
		"list",
		"modify",
		"remove",
		"show",
		"update",
		"version",
		"purge",
		"shell",
	}

	formats = []string{
		"2006-01-02", // "2023-01-23"
		"2006/01/02", // "2023/01/23"
		"02/01/2006", // "23/01/2023"
		"01/02/2006", // "01/23/2023"
	}

	tableHeader = []string{"ID", "Project", "Task", "Priority", "Created", "Due", "Done", "Tags"}
)

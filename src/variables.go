package main

import "database/sql"

const (
	version = "0.9.0"
)

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

	tableHeader = []string{"ID", "Project", "Task", "Priority", "Created", "Due", "Done", "Tags"}
)

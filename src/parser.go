package main

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/wbabcock/TaskNinja/src/utils"
)

func parseInput(args []string) {
	// if only 2 args passed handle them (e.g. taskninja 1 done, taskninja done 1)
	if len(args) == 2 && args[0][:1] != "+" && args[1][:1] != "+" {
		if utils.IsNumeric(args[0]) {
			cmd = args[1]
			id, _ = strconv.ParseUint(args[0], 10, 64)
		}

		if utils.IsNumeric(args[1]) {
			cmd = args[0]
			id, _ = strconv.ParseUint(args[1], 10, 64)
		}
	} else {
		for i, v := range args {
			// look for the id in the first two args and set it
			if i == 0 && utils.IsNumeric(v) {
				id, _ = strconv.ParseUint(v, 10, 64)
			}

			if i == 1 && utils.IsNumeric(v) {
				id, _ = strconv.ParseUint(v, 10, 64)
			}

			// Parse the rest of the args
			if !utils.IsNumeric(v) {
				// Check for verb keyword
				if len(cmd) == 0 {
					if utils.SliceContains(verbs, v) {
						cmd = v
					}
				} else if v[:1] == "+" {
					// Add Tags
					tagsAdd = append(tagsAdd, v[1:])
				} else if v[:1] == "-" {
					// Remove Tags
					tagsRemove = append(tagsRemove, v[1:])
				} else if len(v) >= 8 && v[:8] == "project:" {
					// Set project
					projPassed = true
					proj = v[8:]
				} else if len(v) >= 4 && v[:4] == "due:" {
					// Set the due date
					dueDatePassed = true
					dueDate = parseDueDate(v[4:])
				} else if len(v) >= 9 && v[:9] == "priority:" {
					// Set priority
					switch v[9:] {
					case "H", "h":
						priority = 4
					case "M", "m":
						priority = 3
					case "L", "l":
						priority = 1
					default:
						priority = 2
					}
				} else {
					// everything else is the task description
					desc = append(desc, v)
				}
			}
		}
	}
}

func parseDueDate(d string) sql.NullTime {
	switch {
	case strings.Contains(d, "-"):
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, d)
		if err != nil {
			color.Red("bad due date format supplied")
		}
		return sql.NullTime{
			Time:  parsedDate,
			Valid: true,
		}
	case strings.ToLower(d) == "today":
		t := time.Now()
		y, m, d := t.Date()
		return sql.NullTime{
			Time:  time.Date(y, m, d, 23, 59, 59, 0, time.Local),
			Valid: true,
		}
	case strings.ToLower(d) == "tomorrow":
		t := time.Now().AddDate(0, 0, 1)
		y, m, d := t.Date()
		return sql.NullTime{
			Time:  time.Date(y, m, d, 23, 59, 59, 0, time.Local),
			Valid: true,
		}
	case strings.ToLower(d) == "eow":
		t := time.Now()
		daysUntilSunday := int(time.Saturday - t.Weekday())
		eow := t.AddDate(0, 0, daysUntilSunday)
		y, m, d := eow.Date()
		return sql.NullTime{
			Time:  time.Date(y, m, d, 23, 59, 59, 0, time.Local),
			Valid: true,
		}
	case strings.ToLower(d) == "eoww":
		t := time.Now()
		daysUntilSaturday := int(time.Friday - t.Weekday())
		eow := t.AddDate(0, 0, daysUntilSaturday)
		y, m, d := eow.Date()
		return sql.NullTime{
			Time:  time.Date(y, m, d, 23, 59, 59, 0, time.Local),
			Valid: true,
		}
	case strings.ToLower(d) == "eom":
		t := time.Now()
		firstOfNextMonth := t.AddDate(0, 1, -t.Day()+1)
		eom := firstOfNextMonth.AddDate(0, 0, -1)
		y, m, d := eom.Date()
		return sql.NullTime{
			Time:  time.Date(y, m, d, 23, 59, 59, 0, time.Local),
			Valid: true,
		}
	case strings.ToLower(d) == "eoy":
		t := time.Now()
		y, _, _ := t.Date()
		return sql.NullTime{
			Time:  time.Date(y, 12, 31, 23, 59, 59, 0, time.Local),
			Valid: true,
		}
	}

	return sql.NullTime{}
}

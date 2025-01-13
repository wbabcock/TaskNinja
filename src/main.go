package main

import (
	"cmp"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/wbabcock/TaskNinja/src/db"
	"github.com/wbabcock/TaskNinja/src/utils"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

const (
	version = "0.8.1"
)

var (
	cmd           string
	desc          []string
	proj          string
	projPassed    bool   = false
	priority      uint64 = 0
	id            uint64
	tagsAdd       []string
	tagsRemove    []string
	dueDate       sql.NullTime = sql.NullTime{}
	dueDatePassed bool         = false

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
	}
)

func init() {
	// default location for the database
	homeDir, _ := os.UserHomeDir()
	dbPath := homeDir + "/taskninja_data.sqlite"

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		err := db.Setup_databae(dbPath)
		if err != nil {
			fmt.Println("Error setting up database:", err)
		}
	} else if err != nil {
		fmt.Println("Error checking file existence:", err)
		return
	} else {
		err := db.Connect_databae(dbPath)
		if err != nil {
			fmt.Println("Error opening database:", err)
		}
	}
}

func main() {
	parseInput(os.Args[1:])
	switch cmd {
	case "version":
		fmt.Printf("taskninja version %s\n", version)
	case "add":
		if priority == 0 {
			priority = 2 // default value
		}
		task := db.Task{
			Description: strings.Join(desc, " "),
			Project:     utils.ToNullString(proj),
			Priority:    priority,
			DueDTM:      dueDate,
			CreatedDTM:  time.Now(),
		}
		task.Save()

		// Manage Tags
		for _, v := range tagsAdd {
			tag := db.Tag{
				TaskId: task.Id,
				Name:   v,
			}
			tag.Save()
		}

		for _, v := range tagsRemove {
			db.DeleteTagByName(uint64(task.Id), v)
		}
	case "modify":
		task, err := db.GetTaskById(id)
		if err != nil {
			fmt.Println(err)
			db.Disconnect_database()
			os.Exit(1)
		}

		if dueDatePassed {
			task.DueDTM = dueDate
		}

		task.Description = cmp.Or(strings.Join(desc, " "), task.Description)

		if projPassed {
			task.Project = utils.ToNullString(proj)
		}

		if priority > 0 {
			task.Priority = priority
		}
		task.Update()

		// Manage Tags
		for _, v := range tagsAdd {
			tag := db.Tag{
				TaskId: task.Id,
				Name:   v,
			}
			tag.Save()
		}

		for _, v := range tagsRemove {
			db.DeleteTagByName(uint64(task.Id), v)
		}
	case "complete", "done":
		task, err := db.GetTaskById(id)
		if err != nil {
			fmt.Println(err)
			db.Disconnect_database()
			os.Exit(1)
		}
		task.CompletedDTM = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		task.Update()
	case "purge":
		err := db.PurgeCompletedTasks()
		if err != nil {
			fmt.Println(err)
			db.Disconnect_database()
			os.Exit(1)
		}
		c := color.New(color.FgGreen)
		c.Printf("Completed tasks have been remove\n")
	case "remove", "delete":
		err := db.DeleteTaskById(id)
		if err != nil {
			fmt.Println(err)
			db.Disconnect_database()
			os.Exit(1)
		}
		c := color.New(color.FgGreen)
		c.Printf("Task %d has been remove\n", id)
	case "list", "show":
		listTasks()
	}

	// Disconnect db
	db.Disconnect_database()
}

func listTasks() {
	tasks, err := db.ListTasks(tagsAdd)
	if err != nil {
		fmt.Println(err)
		db.Disconnect_database()
		os.Exit(1)
	}

	if len(tasks) == 0 {
		color.Green("you have nothing to do!")
		db.Disconnect_database()
		os.Exit(1)
	}
	fmt.Println()
	tbl := tablewriter.NewWriter(os.Stdout)
	tbl.SetHeader([]string{"ID", "Project", "Task", "Priority", "Created", "Due", "Done", "Tags"})
	tbl.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	//tbl.SetHeaderLine(false)
	//tbl.SetBorder(false)
	tbl.SetColumnSeparator("|")
	tbl.SetCenterSeparator("+")
	tbl.SetRowSeparator("-")
	tbl.SetAutoWrapText(false)
	tbl.SetTablePadding("\t")
	//tbl.SetNoWhiteSpace(true)
	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, t := range tasks {
		due := t.DueDTM.Time.Format("01/02/2006")
		if !t.DueDTM.Valid {
			due = ""
		}
		comp := t.CompletedDTM.Time.Format("01/02/2006")
		if !t.CompletedDTM.Valid {
			comp = ""
		}

		p := ""
		switch t.Priority {
		case 1:
			p = "L"
		case 2:
			p = ""
		case 3:
			p = "M"
		case 4:
			p = "H"
		}

		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		today := time.Now()
		dueCheckDay := t.DueDTM.Time.Add(-24 * time.Hour)
		dueCheckThreeDay := t.DueDTM.Time.Add(-72 * time.Hour)

		if today.After(dueCheckDay) {
			due = red(due)
		} else if today.After(dueCheckThreeDay) {
			due = yellow(due)
		}

		row := []string{
			fmt.Sprintf("%d", t.Id),
			t.Project.String,
			t.Description,
			p,
			t.CreatedDTM.Format("01/02/2006"),
			due,
			comp,
			t.Tags,
		}

		tbl.Append(row)
	}

	tbl.Render()

	fmt.Println()
}

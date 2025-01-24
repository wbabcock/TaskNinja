package main

import (
	"cmp"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/wbabcock/TaskNinja/src/db"
	"github.com/wbabcock/TaskNinja/src/utils"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

const (
	version = "0.9.0"
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
		createTask()
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

		task.Description = cmp.Or(desc, task.Description)

		if projPassed {
			task.Project = utils.ToNullString(proj)
		}

		if priority > 0 {
			task.Priority = priority
		}
		task.Update()
	case "complete", "done":
		_, err := completeTask(id)
		if err != nil {
			db.Disconnect_database()
			os.Exit(1)
		}
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
	case "shell":
		loadShell()
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
	tbl.SetHeader(tableHeader)
	tbl.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	tbl.SetColumnSeparator("|")
	tbl.SetCenterSeparator("+")
	tbl.SetRowSeparator("-")
	tbl.SetAutoWrapText(false)
	tbl.SetTablePadding("\t")
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
			getPriorityValue(t.Priority),
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

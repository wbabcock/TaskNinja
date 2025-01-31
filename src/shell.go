package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wbabcock/TaskNinja/src/db"
	"github.com/wbabcock/TaskNinja/src/internals/priority"
)

var (
	// these are only applicable here
	app   *tview.Application
	table *tview.Table

	txtID          = tview.NewInputField().SetLabel("ID: ")
	txtProject     = tview.NewInputField().SetLabel("PROJECT: ")
	txtDescription = tview.NewTextArea().SetLabel("TASK: ")
	cboPriority    = tview.NewDropDown().
			SetLabel("PRIORITY: ").
			SetOptions([]string{" Low ", " Normal ", " Medium ", " High "}, func(text string, index int) {}).
			SetCurrentOption(1)
	txtCreated = tview.NewInputField().SetLabel("CREATED: ")
	txtDue     = tview.NewInputField().SetLabel("DUE: ").
			SetAcceptanceFunc(tview.InputFieldMaxLength(10))
	txtDone = tview.NewInputField().SetLabel("DONE: ").
		SetAcceptanceFunc(tview.InputFieldMaxLength(10))
	txtTags = tview.NewInputField().SetLabel("TAGS: ")

	currentRow = 1
	colorMain  = tview.Styles.PrimaryTextColor
	colorDue   = tview.Styles.PrimaryTextColor

	message = tview.NewTextView()

	tasks []db.Task
)

func msgbox(msg string) {
	go func() {
		message.SetText(msg)
		time.Sleep(3 * time.Second)
		message.SetText("")
		app.Draw()
	}()
}

func loadTasks() error {
	var err error
	tasks, err = db.ListTasks(tagsAdd)
	if err != nil {
		return err
	}
	return nil
}

func initFormControls() {
	txtID.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyTab && event.Key() != tcell.KeyBacktab {
			return nil
		}
		return event
	})

	txtCreated.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyTab && event.Key() != tcell.KeyBacktab {
			return nil
		}
		return event
	})
}

func loadTable() {
	if len(tasks) == 0 {
		table.SetCell(0, 0,
			tview.NewTableCell("There are not tasks to show").
				SetTextColor(colorMain))
	} else {
		// setup the header
		for i, col := range tableHeader {
			table.SetCell(0, i,
				tview.NewTableCell(col).
					SetTextColor(tview.Styles.SecondaryTextColor).
					SetSelectable(false))
		}
	}

	// table content
	for row, task := range tasks {
		due := ""
		if task.DueDTM.Valid {
			due = task.DueDTM.Time.Format("01/02/2006")

			today := time.Now()
			dueCheckDay := task.DueDTM.Time.Add(-24 * time.Hour)
			dueCheckThreeDay := task.DueDTM.Time.Add(-72 * time.Hour)

			if today.After(dueCheckDay) {
				colorDue = tcell.ColorRed
			} else if today.After(dueCheckThreeDay) {
				colorDue = tcell.ColorYellow
			}
		} else {
			// reset the color to default
			colorDue = tview.Styles.TitleColor
		}
		done := ""
		if task.CompletedDTM.Valid {
			done = task.CompletedDTM.Time.Format("01/02/2006")
		}

		table.SetCell(row+1, 0,
			tview.NewTableCell(fmt.Sprintf("%d", task.Id)).
				SetTextColor(colorMain).
				SetAlign(tview.AlignLeft))
		table.SetCell(row+1, 1,
			tview.NewTableCell(task.Project.String).
				SetTextColor(colorMain).
				SetAlign(tview.AlignLeft))
		table.SetCell(row+1, 2,
			tview.NewTableCell(task.Description).
				SetTextColor(colorMain).
				SetAlign(tview.AlignLeft))
		table.SetCell(row+1, 3,
			tview.NewTableCell(task.Priority.String()).
				SetTextColor(colorMain).
				SetAlign(tview.AlignCenter))
		table.SetCell(row+1, 4,
			tview.NewTableCell(task.CreatedDTM.Format("01/02/2006")).
				SetTextColor(colorMain).
				SetAlign(tview.AlignLeft))
		table.SetCell(row+1, 5,
			tview.NewTableCell(due).
				SetTextColor(colorDue).
				SetAlign(tview.AlignLeft))
		table.SetCell(row+1, 6,
			tview.NewTableCell(done).
				SetTextColor(colorMain).
				SetAlign(tview.AlignLeft))
		table.SetCell(row+1, 7,
			tview.NewTableCell(task.Tags).
				SetTextColor(colorMain).
				SetAlign(tview.AlignLeft))

	}
}

func refreshTable() {
	table.Clear()
	loadTasks()
	loadTable()
	table.Select(1, 0)
}

func refreshForm() {
	txtID.SetText("")
	txtProject.SetText("")
	txtDescription.SetText("", true)
	cboPriority.SetCurrentOption(1)
	txtCreated.SetText("")
	txtDue.SetText("")
	txtDone.SetText("")
	txtTags.SetText("")
}

func loadFormDetails(table *tview.Table) {
	id, _ = strconv.ParseUint(table.GetCell(currentRow, 0).Text, 10, 64)
	proj = table.GetCell(currentRow, 1).Text
	desc = table.GetCell(currentRow, 2).Text
	taskPriority = priority.GetIndex(table.GetCell(currentRow, 3).Text)

	tagsRemove = strings.Split(strings.ReplaceAll(strings.ReplaceAll(table.GetCell(currentRow, 7).Text, ", ", ","), " ", "_"), ",")

	txtID.SetText(fmt.Sprintf("%d", id))
	txtProject.SetText(proj)
	txtDescription.SetText(desc, true)
	cboPriority.SetCurrentOption(int(taskPriority) - 1)

	txtCreated.SetText(table.GetCell(currentRow, 4).Text)

	txtDue.SetText(table.GetCell(currentRow, 5).Text)
	txtDone.SetText(table.GetCell(currentRow, 6).Text)
	txtTags.SetText(table.GetCell(currentRow, 7).Text)
}

func loadShell() {
	app = tview.NewApplication()

	// load tasks
	err := loadTasks()
	if err != nil {
		log.Fatal(err)
		db.Disconnect_database()
		os.Exit(1)
	}

	// setup of the controls for the form
	initFormControls()

	// initialize the table
	table = tview.NewTable().Select(currentRow, 0).SetSelectable(true, false)
	loadTable()

	// create the form
	form := tview.NewForm().
		AddFormItem(txtID.SetDisabled(true)).
		AddFormItem(txtProject).
		AddFormItem(txtDescription).
		AddFormItem(cboPriority).
		AddFormItem(txtCreated.SetDisabled(true)).
		AddFormItem(txtDue).
		AddFormItem(txtDone).
		AddFormItem(txtTags).
		AddButton("Save", func() {
			idStr := txtID.GetText()
			id, _ := strconv.ParseUint(idStr, 10, 64)
			proj = strings.ReplaceAll(txtProject.GetText(), " ", "_")
			projPassed = true
			desc = txtDescription.GetText()
			p, _ := cboPriority.GetCurrentOption()
			taskPriority = priority.New(uint8(p + 1))
			tagsAdd = strings.Split(strings.ReplaceAll(strings.ReplaceAll(txtTags.GetText(), ", ", ","), " ", "_"), ",")

			if txtDue.GetText() == "" {
				dueDatePassed = true
				dueDate = sql.NullTime{}
			} else {
				parsedDate, err := parseDate(txtDue.GetText())
				if err != nil {
					msgbox("ERROR: " + err.Error())
					dueDatePassed = false
				}
				dueDatePassed = true
				dueDate = sql.NullTime{
					Time:  parsedDate,
					Valid: true,
				}
			}

			if txtDone.GetText() == "" {
				doneDatePassed = true
				doneDate = sql.NullTime{}
			} else {
				parsedDate, err := parseDate(txtDone.GetText())
				if err != nil {
					msgbox("ERROR: " + err.Error())
					doneDatePassed = false
				}
				doneDatePassed = true
				doneDate = sql.NullTime{
					Time:  parsedDate,
					Valid: true,
				}
			}

			if idStr == "" {
				createTask()
			} else {
				updateTask(id)
			}
			tagsAdd = []string{}
			refreshForm()
			refreshTable()
			app.SetFocus(table)
		}).
		AddButton("Clear", func() {
			refreshForm()
		}).
		AddButton("Complete", func() {
			idStr := txtID.GetText()
			if idStr == "" {
				msgbox("ERROR: You must select a task to complete")
				return
			}
			id, _ := strconv.ParseUint(idStr, 10, 64)
			completeTask(id)
			refreshForm()
			refreshTable()
			app.SetFocus(table)
		}).
		AddButton("Remove", func() {
			idStr := txtID.GetText()
			if idStr == "" {
				return
			}
			id, _ := strconv.ParseUint(idStr, 10, 64)
			err := db.DeleteTaskById(id)
			if err != nil {
				log.Fatal(err)
				db.Disconnect_database()
				os.Exit(1)
			}
			refreshForm()
			refreshTable()
			app.SetFocus(table)
		})
	form.SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)

	table.SetSelectedFunc(func(row int, column int) {
		currentRow = row
		loadFormDetails(table)
		app.SetFocus(form)
	})
	table.SetBackgroundColor(tcell.ColorDefault).
		SetBorder(true)

	//-----------------------------------------------------------------------------
	//		SETUP APP LAYOUT AND KEYBINDINGS
	//-----------------------------------------------------------------------------

	// App Title at top
	appVersion := fmt.Sprintf("version: %s", version)
	appTitle := tview.NewFlex().
		AddItem(tview.NewTextView().
			SetText("Task Ninja - Shell").
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetDynamicColors(true).
			SetTextAlign(tview.AlignCenter), 0, 1, false).
		AddItem(tview.NewTextView().
			SetText(appVersion).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetDynamicColors(true).
			SetTextAlign(tview.AlignRight), len(appVersion), 0, false)

	// Key menu at the bottom
	keyMenu := tview.NewTextView().
		SetText("ESC: Quit | ENTER: Select | TAB: Next | CTRL+N: Change Pane").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetTextColor(tview.Styles.SecondaryTextColor)

	message = tview.NewTextView().
		SetText("").
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorRed)

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(appTitle, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(table, 0, 2, true).
			AddItem(form, 0, 1, false), 0, 2, true).
		AddItem(message, 1, 0, false).
		AddItem(keyMenu, 1, 0, false)

	// Keybindings
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			// Toggle focus between the two items
			if app.GetFocus() == table {
				app.SetFocus(form)
			} else {
				refreshForm()
				app.SetFocus(table)
			}
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			// Exit the application on Escape
			db.Disconnect_database()
			app.Stop()
			return nil
		}
		return event
	})

	// Run program
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

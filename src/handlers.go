package main

import (
	"cmp"
	"database/sql"
	"time"

	"github.com/wbabcock/TaskNinja/src/db"
	"github.com/wbabcock/TaskNinja/src/utils"
)

func createTask() (db.Task, error) {
	if priority == 0 {
		priority = 2 // default value
	}
	task := db.Task{
		Description:  desc,
		Project:      utils.ToNullString(proj),
		Priority:     priority,
		DueDTM:       dueDate,
		CreatedDTM:   time.Now(),
		CompletedDTM: doneDate,
	}
	err := task.Save()
	if err != nil {
		return db.Task{}, err
	}

	// Manage Tags
	for _, v := range tagsAdd {
		tag := db.Tag{
			TaskId: task.Id,
			Name:   v,
		}
		tag.Save()
	}

	return task, nil
}

func updateTask(taskId uint64) (db.Task, error) {
	task, err := db.GetTaskById(taskId)
	if err != nil {
		return db.Task{}, err
	}

	if dueDatePassed {
		task.DueDTM = dueDate
	}

	if doneDatePassed {
		task.CompletedDTM = doneDate
	}

	task.Description = cmp.Or(desc, task.Description)

	if projPassed {
		task.Project = utils.ToNullString(proj)
	}

	if priority > 0 {
		task.Priority = priority
	}

	err = task.Update()
	if err != nil {
		return db.Task{}, err
	}

	// Manage Tags
	for _, v := range tagsRemove {
		db.DeleteTagByName(uint64(task.Id), v)
	}

	for _, v := range tagsAdd {
		tag := db.Tag{
			TaskId: task.Id,
			Name:   v,
		}
		tag.Save()
	}

	return task, nil
}

func completeTask(taskId uint64) (db.Task, error) {
	task, err := db.GetTaskById(taskId)
	if err != nil {
		return db.Task{}, err
	}
	task.CompletedDTM = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	err = task.Update()
	if err != nil {
		return db.Task{}, err
	}

	return task, nil
}

func parseDate(input string) (time.Time, error) {
	var err error
	for _, format := range formats {
		parsedDate, err := time.Parse(format, input)
		if err == nil {
			return parsedDate, nil
		}
	}
	return time.Now(), err
}

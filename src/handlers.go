package main

import (
	"cmp"
	"database/sql"
	"time"

	"github.com/wbabcock/TaskNinja/src/db"
	"github.com/wbabcock/TaskNinja/src/internals/priority"
	"github.com/wbabcock/TaskNinja/src/utils"
)

func createTask() (db.Task, error) {
	if taskPriority == 0 {
		taskPriority = priority.Normal // default value
	}
	task := db.Task{
		Description:  desc,
		Project:      utils.ToNullString(proj),
		Priority:     taskPriority,
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

	if taskPriority > 0 {
		task.Priority = taskPriority
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
	format, err := utils.GetLocaleDateFormat()
	if err != nil {
		return time.Now(), err
	}

	parsedDate, err := time.Parse(format, input)
	if err != nil {
		return time.Now(), err
	}

	return parsedDate, err
}

package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	Id           uint64
	Description  string
	Project      sql.NullString
	Priority     uint64
	CreatedDTM   time.Time
	DueDTM       sql.NullTime
	CompletedDTM sql.NullTime
	Tags         string
}

func ListTasks(tags []string) ([]Task, error) {

	var tagFilter string
	if len(tags) > 0 {
		tagFilter = " where "

		for _, v := range tags {
			tagFilter += "tags like '%" + v + "%' or "
		}

		tagFilter = tagFilter[:len(tagFilter)-4]
	}

	stmt, err := db.Prepare(fmt.Sprintf(`
		select tasks.id, project, description, priority, created_dtm, due_dtm, completed_dtm, coalesce(tags_view.tags, '') as tags 
		from tasks left join tags_view on tasks.id = tags_view.task_id
		%s  
		order by due_dtm desc, priority desc;
	`, tagFilter))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Project, &task.Description, &task.Priority, &task.CreatedDTM, &task.DueDTM, &task.CompletedDTM, &task.Tags)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil

}

func GetTaskById(id uint64) (Task, error) {
	todo := Task{}
	stmt, err := db.Prepare(`
		select id, project, description, priority, created_dtm, due_dtm, completed_dtm from tasks where id = ?;
	`)
	if err != nil {
		return todo, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	err = row.Scan(&todo.Id, &todo.Project, &todo.Description, &todo.Priority, &todo.CreatedDTM, &todo.DueDTM, &todo.CompletedDTM)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (t *Task) Save() error {
	stmt, err := db.Prepare(`
		insert into tasks (description, project, priority, created_dtm, due_dtm, completed_dtm)
		values (?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(t.Description, t.Project, t.Priority, t.CreatedDTM, t.DueDTM, t.CompletedDTM)
	id, _ := result.LastInsertId()
	t.Id = uint64(id)
	return err
}

func (t *Task) Update() error {
	stmt, err := db.Prepare(`
		update tasks 
			set description=?, 
				project=?, 
				priority=?,
				due_dtm=?, 
				completed_dtm=?
		where id = ?;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Description, t.Project, t.Priority, t.DueDTM, t.CompletedDTM, t.Id)
	return err
}

func DeleteTaskById(id uint64) error {
	stmt, err := db.Prepare(`
		delete from tasks where id = ?;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	ra, _ := r.RowsAffected()
	if ra == 0 {
		return fmt.Errorf("task doesn't exist for id")
	}
	return nil
}

func PurgeCompletedTasks() error {
	stmt, err := db.Prepare(`
		delete from tasks where completed_dtm is not null;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

package db

import "fmt"

type Tag struct {
	Id     uint64
	TaskId uint64
	Name   string
}

func (t *Tag) Save() error {
	stmt, err := db.Prepare(`
		insert into tags (task_id, name)
		values (?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(t.TaskId, t.Name)
	id, _ := result.LastInsertId()
	t.Id = uint64(id)
	return err
}

func DeleteTagByName(todoId uint64, name string) error {
	stmt, err := db.Prepare(`
		delete from tags where task_id = ? and name = ?;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(todoId, name)
	if err != nil {
		return err
	}

	ra, _ := r.RowsAffected()
	if ra == 0 {
		return fmt.Errorf("tag doesn't exist for task")
	}
	return nil
}

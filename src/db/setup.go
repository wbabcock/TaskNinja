package db

func Setup_databae(dbPath string) error {
	// Create a table in the database
	Connect_databae(dbPath)
	err := create_task_table()
	if err != nil {
		return err
	}

	err = create_tag_table()
	if err != nil {
		return err
	}

	err = create_tag_view()
	if err != nil {
		return err
	}
	return nil
}

func create_task_table() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY, 
			description VARCHAR(150),
			project VARCHAR(50),
			priority INTERGER,
			created_dtm DATETIME,
			due_dtm DATETIME,
			completed_dtm DATETIME,
    		CONSTRAINT description_length CHECK (LENGTH(description) <= 150),
    		CONSTRAINT project_length CHECK (LENGTH(project) <= 50)
		)`)
	if err != nil {
		return err
	}

	return nil
}

func create_tag_table() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY, 
			task_id INTEGER,
			name VARCHAR(50),
    		CONSTRAINT name_length CHECK (LENGTH(name) <= 50)
		)`)
	if err != nil {
		return err
	}

	return nil
}

func create_tag_view() error {
	_, err := db.Exec(`
		CREATE VIEW tags_view AS
			SELECT task_id, GROUP_CONCAT(name, ',') AS tags
			FROM tags
			GROUP BY task_id;
	`)
	if err != nil {
		return err
	}

	return nil
}

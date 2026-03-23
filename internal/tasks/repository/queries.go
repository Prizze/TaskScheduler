package repository

const (
	getTagByID = `
		SELECT id, name
		FROM tags
		WHERE id = $1 AND user_id = $2;
	`

	createTask = `
		INSERT INTO tasks (user_id, title, description, status_id, priority_id, due_date)
		VALUES (
			$1,
			$2,
			$3,
			(SELECT id FROM task_status WHERE name = $4),
			(SELECT id FROM task_priority WHERE name = $5),
			$6
		)
		RETURNING id, user_id, title, description, $4, $5, due_date, created_at, updated_at;
	`

	createTaskTag = `
		INSERT INTO task_tags (task_id, tag_id)
		VALUES ($1, $2);
	`
)

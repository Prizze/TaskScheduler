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

	getTaskByID = `
		SELECT
			t.id,
			t.user_id,
			t.title,
			t.description,
			ts.name,
			tp.name,
			t.due_date,
			t.created_at,
			t.updated_at
		FROM tasks t
		JOIN task_status ts ON ts.id = t.status_id
		JOIN task_priority tp ON tp.id = t.priority_id
		WHERE t.id = $1 AND t.user_id = $2;
	`

	getTaskTags = `
		SELECT tg.id, tg.name
		FROM task_tags tt
		JOIN tags tg ON tg.id = tt.tag_id
		WHERE tt.task_id = $1
		ORDER BY tg.id;
	`

	updateTask = `
		UPDATE tasks
		SET
			title = $3,
			description = $4,
			status_id = (SELECT id FROM task_status WHERE name = $5),
			priority_id = (SELECT id FROM task_priority WHERE name = $6),
			due_date = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, title, description, $5, $6, due_date, created_at, updated_at;
	`

	deleteTaskTags = `
		DELETE FROM task_tags
		WHERE task_id = $1;
	`

	createTaskTag = `
		INSERT INTO task_tags (task_id, tag_id)
		VALUES ($1, $2);
	`

	deleteTask = `
		DELETE FROM tasks
		WHERE id = $1 AND user_id = $2;
	`
)

package repository

const (
	createTag = `
		INSERT INTO tags (user_id, name)
		VALUES ($1, $2)
		RETURNING id, user_id, name, created_at;
	`

	getTags = `
		SELECT id, user_id, name, created_at
		FROM tags
		WHERE user_id = $1
		ORDER BY id;
	`

	deleteTag = `
		DELETE FROM tags
		WHERE id = $1 AND user_id = $2;
	`
)

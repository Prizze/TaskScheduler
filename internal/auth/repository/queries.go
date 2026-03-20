package repository

const (
	getUserByID = `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	getUserByEmail = `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	checkEmail = `
		SELECT id FROM users
		WHERE email = $1;
	`

	createUser = `
		INSERT INTO users (email, password_hash, name)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, name, created_at, updated_at;
	`
)

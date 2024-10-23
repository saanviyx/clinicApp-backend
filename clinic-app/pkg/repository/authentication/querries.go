package authentication

const (
	// Register a new user
	RegisterUserQuery = `
		INSERT INTO Users (
			username, 
			name, 
			email, 
			password, 
			role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id;
	`
	// Login
	LoginUserQuery = `
		SELECT 
			user_id, 
			username, 
			password, 
			role
		FROM Users
		WHERE username = $1 and password = $2;
	`
)

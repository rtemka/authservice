package domain

// User - user's entity
type User struct {
	Login string `json:"login"`
	Password
	CreatedAt  int64 `json:"created_at"`
	IsDisabled bool  `json:"is_disabled"`
}

// Password is the hashed password
type Password struct {
	Hash        string `json:"-"`
	GeneratedAt int64  `json:"generated_at"`
	IsActive    bool   `json:"is_active"`
}

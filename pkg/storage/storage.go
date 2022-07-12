package storage

type Storage interface {
	AddUser(login string)
}

type User struct {
	Login
	Password
	IsDisabled bool
}

type Login struct {
	Id         int64  `json:"login_id"`
	Name       string `json:"login"`
	PasswordID int64  `json:"-"`
}

type Password struct {
	Id       int64  `json:"password_id"`
	Hash     string `json:"hash"`
	Salt     string `json:"-"`
	Created  int64  `json:"created"`
	IsActive bool   `json:"is_active"`
}

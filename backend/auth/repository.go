package auth

import (
	_ "database/sql"
	_ "fmt"
	_ "github.com/lib/pq"
	"spaced-ace-backend/utils"
)

type DBUser struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type Session struct {
	Id         string `db:"id"`
	UserId     string `db:"user_id"`
	ValidUntil string `db:"valid_until"`
}

var schema = `
CREATE EXTENSION IF NOT EXISTS pg_cron;
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	name TEXT,
	email TEXT,
	password TEXT
);
CREATE INDEX IF NOT EXISTS users_email ON users(email);
CREATE UNLOGGED TABLE IF NOT EXISTS sessions (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id),
	valid_until TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS sessions_id ON sessions(id);
CREATE INDEX IF NOT EXISTS sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS sessions_valid_until ON sessions(valid_until);
SELECT cron.schedule('10 * * * *', $$DELETE FROM sessions WHERE valid_until < now()$$);`

func init() {
	utils.DB.MustExec(schema)
}

func GetUserByEmail(email string) (*DBUser, error) {
	user := DBUser{}
	err := utils.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	return &user, err
}

func GetUserById(id string) (*DBUser, error) {
	user := DBUser{}
	err := utils.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return &user, err
}

func CreateUser(user *DBUser) error {
	_, err := utils.DB.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(user *DBUser) error {
	_, err := utils.DB.Exec("UPDATE users SET name=$2, email=$3, password=$4 WHERE id=$1", user.Id, user.Name, user.Email, user.Password)
	return err
}

func DeleteUser(id string) error {
	_, err := utils.DB.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

func GetUserIdBySession(sessionId string) (string, error) {
	var id string
	err := utils.DB.Get(&id, "SELECT user_id FROM sessions WHERE id=$1", sessionId)
	return id, err
}

func CreateSession(userId string) (string, error) {
	var id string
	err := utils.DB.Get(&id, "INSERT INTO sessions (id, user_id, valid_until) VALUES (gen_random_uuid(), $1, now() + interval '1 hour') RETURNING id", userId)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteSession(id string) error {
	_, err := utils.DB.Exec("DELETE FROM sessions WHERE id=$1", id)
	return err
}

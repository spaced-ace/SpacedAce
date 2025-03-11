package auth

import (
	_ "database/sql"
	_ "fmt"
	_ "github.com/lib/pq"
	"spaced-ace-backend/utils"
)

type DBUser struct {
	Id                string  `db:"id"`
	Name              string  `db:"name"`
	Email             string  `db:"email"`
	Password          string  `db:"password"`
	EmailVerified     bool    `db:"email_verified"`
	VerificationToken *string `db:"verification_token"`
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
	password TEXT,
	email_verified BOOLEAN DEFAULT FALSE,
	verification_token TEXT
);
CREATE INDEX IF NOT EXISTS users_email ON users(email);
CREATE INDEX IF NOT EXISTS users_verification_token ON users(verification_token);
CREATE UNLOGGED TABLE IF NOT EXISTS sessions (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	valid_until TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS sessions_id ON sessions(id);
CREATE INDEX IF NOT EXISTS sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS sessions_valid_until ON sessions(valid_until);
SELECT cron.schedule('del_exp_sessions', '10 * * * *', $$DELETE FROM sessions WHERE valid_until < now()$$);`

func InitDb() {
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
	_, err := utils.DB.Exec("INSERT INTO users (id, name, email, password, email_verified, verification_token) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Id, user.Name, user.Email, user.Password, user.EmailVerified, user.VerificationToken)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(user *DBUser) error {
	_, err := utils.DB.Exec("UPDATE users SET name=$2, email=$3, password=$4, email_verified=$5, verification_token=$6 WHERE id=$1",
		user.Id, user.Name, user.Email, user.Password, user.EmailVerified, user.VerificationToken)
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

func GetUserByVerificationToken(token string) (*DBUser, error) {
	user := DBUser{}
	err := utils.DB.Get(&user, "SELECT * FROM users WHERE verification_token=$1", token)
	return &user, err
}

func VerifyEmail(id string) error {
	_, err := utils.DB.Exec("UPDATE users SET email_verified=true, verification_token=NULL WHERE id=$1", id)
	return err
}

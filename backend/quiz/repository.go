package quiz

import (
	"database/sql"
	_ "fmt"
	"spaced-ace-backend/utils"
)

var schema = `
	CREATE TABLE IF NOT EXISTS quizzes(
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT
	)
	CREATE TABLE IF NOT EXISTS quiz_accesses(
		userid UUID REFERENCES users(id),
		quizid UUID REFERENCES quizzes(id),
		roleid SMALLINT NOT NULL, --1 = owner, 2 = viewer
		PRIMARY KEY(userid, quizid, roleid),
		UNIQUE(userid, quizid)
	)
	`
var (
	QUIZ_OWNER_ACCESS_ID  = 1
	QUIZ_VIEWER_ACCESS_ID = 2
)

func init() {
	utils.DB.MustExec(schema)
}

type DBQuiz struct {
	Id          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}

func CreateQuiz(ownerid string, name string, description string) (*DBQuiz, error) {
	quiz := DBQuiz{}
	err := utils.DB.Get(&quiz, "INSERT INTO quizzes (id, name, description) VALUES (gen_random_uuid(), $1, $2) RETURNING *", name, description)
	if err != nil {
		return nil, err
	}
	err = CreateQuizAccess(ownerid, quiz.Id, QUIZ_OWNER_ACCESS_ID)
	if err != nil {
		return nil, err
	}
	return &quiz, err
}

func CreateQuizAccess(userid string, quizid string, roleid int) error {
	_, err := utils.DB.Exec("INSERT INTO quiz_accesses(userid, quizid, roleid) VALUES ($1, $2, $3)", userid, quizid, roleid)
	return err
}

func GetQuizById(id string) (*DBQuiz, error) {
	quiz := DBQuiz{}
	err := utils.DB.Get(&quiz, "SELECT * FROM quizzes WHERE id = $id", id)
	return &quiz, err
}

func GetQuizAccess(userid string, quizid string) (int, error) {
	var access = 0
	err := utils.DB.Select(&access, "SELECT roleid FROM quiz_accesses WHERE userid = $1 AND quizid = $2", userid, quizid)
	if err != nil {
		return access, err
	}
	return access, nil
}

func GetQuizAccessesOfUser(userid string) (*[]string, error) {
	var accesses []string
	err := utils.DB.Select(&accesses, "SELECT roleid FROM quiz_accesses WHERE userid = $1", userid)
	return &accesses, err
}

func GetQuizAccesses(quizid string) (*[]string, error) {
	var accesses []string
	err := utils.DB.Select(&accesses, "SELECT roleid FROM quiz_accesses WHERE quizid = $1", quizid)
	return &accesses, err
}

func UpdateQuizAccess(userid string, quizid string, roleid int) error {
	_, err := utils.DB.Exec("UPDATE quiz_accesses SET roleid=$3 WHERE userid = $1 AND quizid = $2", userid, quizid, roleid)
	return err
}

func UpdateQuiz(userid string, quizid string, name string, description string) error {
	if name == "" && description != "" {
		_, err := utils.DB.Exec("UPDATE quizzes SET description=$2 WHERE id=$1", quizid, description)
		return err
	}
	if name != "" && description == "" {
		_, err := utils.DB.Exec("UPDATE quizzes SET name=$2 WHERE id=$1", quizid, name)
		return err
	}
	_, err := utils.DB.Exec("UPDATE quizzes SET name=$2, description=$3 WHERE id=$1", quizid, name, description)
	return err
}

func DeleteQuizAccess(userid string, quizid string) error {
	_, err := utils.DB.Exec("DELETE FROM quiz_accesses WHERE userid = $1 and quizid = $2", userid, quizid)
	return err
}

func DeleteQuiz(id string) error {
	_, err := utils.DB.Exec("DELETE FROM quizzes WHERE id = $id", id)
	return err
}

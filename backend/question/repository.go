package question

import (
	_ "fmt"
	"github.com/lib/pq"
	"spaced-ace-backend/utils"
)

var schema = `
CREATE TABLE IF NOT EXISTS multiple_choice_questions (
	uuid UUID PRIMARY KEY,
	quizid UUID REFERENCES quizzes(id) ON DELETE CASCADE,
	question TEXT,
	answers TEXT[4],
	correct_answers CHAR[]
);
CREATE TABLE IF NOT EXISTS single_choice_questions (
	uuid UUID PRIMARY KEY,
	quizid UUID REFERENCES quizzes(id) ON DELETE CASCADE,
	question TEXT,
	answers TEXT[4],
	correct_answer CHAR
);
CREATE TABLE IF NOT EXISTS true_or_false_questions (
	uuid UUID PRIMARY KEY,
	quizid UUID REFERENCES quizzes(id) ON DELETE CASCADE,
	question TEXT,
	correct_answer BOOLEAN
);
`

type DBMultipleChoiceQuestion struct {
	UUID           string         `db:"uuid"`
	QuizID         string         `db:"quizid"`
	Question       string         `db:"question"`
	Answers        pq.StringArray `db:"answers"`
	CorrectAnswers pq.StringArray `db:"correct_answers"`
}
type DBSingleChoiceQuestion struct {
	UUID          string         `db:"uuid"`
	QuizID        string         `db:"quizid"`
	Question      string         `db:"question"`
	Answers       pq.StringArray `db:"answers"`
	CorrectAnswer string         `db:"correct_answer"`
}
type DBTrueOrFalseQuestion struct {
	UUID          string `db:"uuid"`
	QuizID        string `db:"quizid"`
	Question      string `db:"question"`
	CorrectAnswer bool   `db:"correct_answer"`
}

func InitDb() {
	utils.DB.MustExec(schema)
}

func CreateMultipleChoiceQuestion(question *DBMultipleChoiceQuestion) error {
	_, err := utils.DB.Exec(
		"INSERT INTO multiple_choice_questions (uuid, quizid, question, answers, correct_answers) VALUES ($1,$2,$3,$4,$5)",
		question.UUID, question.QuizID, question.Question, question.Answers, question.CorrectAnswers,
	)
	return err
}
func CreateSingleChoiceQuestion(question *DBSingleChoiceQuestion) error {
	_, err := utils.DB.Exec(
		"INSERT INTO single_choice_questions (uuid, quizid, question, answers, correct_answer) VALUES ($1,$2,$3,$4,$5)",
		question.UUID, question.QuizID, question.Question, question.Answers, question.CorrectAnswer,
	)
	return err
}
func CreateTrueOrFalseQuestion(question *DBTrueOrFalseQuestion) error {
	_, err := utils.DB.Exec(
		"INSERT INTO true_or_false_questions (uuid, quizid, question, correct_answer) VALUES ($1,$2,$3,$4)",
		question.UUID, question.QuizID, question.Question, question.CorrectAnswer,
	)
	return err
}

func GetMultipleChoiceQuestions(quizID string) ([]DBMultipleChoiceQuestion, error) {
	questions := []DBMultipleChoiceQuestion{}
	err := utils.DB.Select(&questions, "SELECT * FROM multiple_choice_questions WHERE quizid=$1", quizID)
	return questions, err
}

func GetMultipleChoiceQuestion(uuid string) (DBMultipleChoiceQuestion, error) {
	question := DBMultipleChoiceQuestion{}
	err := utils.DB.Get(&question, "SELECT * FROM multiple_choice_questions WHERE uuid=$1", uuid)
	return question, err
}

func GetSingleChoiceQuestions(quizID string) ([]DBSingleChoiceQuestion, error) {
	questions := []DBSingleChoiceQuestion{}
	err := utils.DB.Select(&questions, "SELECT * FROM single_choice_questions WHERE quizid=$1", quizID)
	return questions, err
}

func GetSingleChoiceQuestion(id string) (DBSingleChoiceQuestion, error) {
	question := DBSingleChoiceQuestion{}
	err := utils.DB.Get(&question, "SELECT * FROM single_choice_questions WHERE uuid=$1", id)
	return question, err
}

func GetTrueOrFalseQuestions(quizID string) ([]DBTrueOrFalseQuestion, error) {
	questions := []DBTrueOrFalseQuestion{}
	err := utils.DB.Select(&questions, "SELECT * FROM true_or_false_questions WHERE quizid=$1", quizID)
	return questions, err
}

func GetTrueOrFalseQuestion(id string) (DBTrueOrFalseQuestion, error) {
	question := DBTrueOrFalseQuestion{}
	err := utils.DB.Get(&question, "SELECT * FROM true_or_false_questions WHERE uuid=$1", id)
	return question, err
}

func DeleteMultipleChoiceQuestion(uuid string) error {
	_, err := utils.DB.Exec("DELETE FROM multiple_choice_questions WHERE uuid=$1", uuid)
	return err
}

func DeleteSingleChoiceQuestion(uuid string) error {
	_, err := utils.DB.Exec("DELETE FROM single_choice_questions WHERE uuid=$1", uuid)
	return err
}

func DeleteTrueOrFalseQuestion(uuid string) error {
	_, err := utils.DB.Exec("DELETE FROM true_or_false_questions WHERE uuid=$1", uuid)
	return err
}

func UpdateMultipleChoiceQuestion(question *DBMultipleChoiceQuestion) error {
	_, err := utils.DB.Exec(
		"UPDATE multiple_choice_questions SET question=$1, answers=$2, correct_answers=$3 WHERE uuid=$4",
		question.Question, question.Answers, question.CorrectAnswers, question.UUID,
	)
	return err
}

func UpdateSingleChoiceQuestion(question *DBSingleChoiceQuestion) error {
	_, err := utils.DB.Exec(
		"UPDATE single_choice_questions SET question=$1, answers=$2, correct_answer=$3 WHERE uuid=$4",
		question.Question, question.Answers, question.CorrectAnswer, question.UUID,
	)
	return err
}

func UpdateTrueOrFalseQuestion(question *DBTrueOrFalseQuestion) error {
	_, err := utils.DB.Exec(
		"UPDATE true_or_false_questions SET question=$1, correct_answer=$2 WHERE uuid=$3",
		question.Question, question.CorrectAnswer, question.UUID,
	)
	return err
}

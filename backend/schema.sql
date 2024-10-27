CREATE TABLE IF NOT EXISTS quiz_sessions(
    id   UUID PRIMARY KEY NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE NOT NULL,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    closes_at TIMESTAMP
);
CREATE INDEX idx_quiz_sessions_user_id ON quiz_sessions(user_id);
CREATE INDEX idx_quiz_sessions_quiz_id ON quiz_sessions(quiz_id);

CREATE TABLE IF NOT EXISTS single_choice_answers(
    id   UUID PRIMARY KEY NOT NULL,
    session_id UUID REFERENCES quiz_sessions(id) NOT NULL,
    question_id UUID REFERENCES single_choice_questions(uuid) ON DELETE CASCADE NOT NULL,
    answer TEXT[1]
);
CREATE INDEX idx_single_choice_answers_session_id ON single_choice_answers(session_id);
ALTER TABLE single_choice_answers ADD CONSTRAINT constraint_single_choice_answers_unique_session_and_question UNIQUE (session_id, question_id);

CREATE TABLE IF NOT EXISTS multiple_choice_answers(
    id   UUID PRIMARY KEY NOT NULL,
    session_id UUID REFERENCES quiz_sessions(id) NOT NULL,
    question_id UUID REFERENCES multiple_choice_questions(uuid) ON DELETE CASCADE NOT NULL,
    answers TEXT[4] -- list of letters e.g. ABD, maximum 4 answers are possible
);
CREATE INDEX idx_multiple_choice_answers_session_id ON multiple_choice_answers(session_id);
ALTER TABLE multiple_choice_answers ADD CONSTRAINT constraint_multiple_choice_answers_unique_session_and_question UNIQUE (session_id, question_id);

CREATE TABLE IF NOT EXISTS true_or_false_answers(
    id   UUID PRIMARY KEY NOT NULL,
    session_id UUID REFERENCES quiz_sessions(id) NOT NULL,
    question_id UUID REFERENCES true_or_false_questions(uuid) ON DELETE CASCADE NOT NULL,
    answer BOOLEAN
);
CREATE INDEX idx_true_or_false_answers_session_id ON true_or_false_answers(session_id);
ALTER TABLE true_or_false_answers ADD CONSTRAINT constraint_true_or_false_answers_unique_session_and_question UNIQUE (session_id, question_id);

CREATE TABLE IF NOT EXISTS quiz_results(
    id UUID PRIMARY KEY NOT NULL,
    session_id UUID REFERENCES quiz_sessions(id) ON DELETE CASCADE NOT NULL,
    max_score FLOAT,
    score FLOAT
);
CREATE INDEX idx_quiz_result_session_id ON quiz_results(session_id);

CREATE TABLE IF NOT EXISTS answer_scores(
    id UUID PRIMARY KEY NOT NULL,
    quiz_result_id UUID REFERENCES quiz_results(id) ON DELETE CASCADE NOT NULL,
    single_choice_answer_id UUID REFERENCES single_choice_answers(id),
    multiple_choice_answer_id UUID REFERENCES multiple_choice_answers(id),
    true_or_false_answer_id UUID REFERENCES true_or_false_answers(id),
    CHECK (
        (single_choice_answer_id IS NOT NULL)::int +
        (multiple_choice_answer_id IS NOT NULL)::int +
        (true_or_false_answer_id IS NOT NULL)::int = 1
    ),
    max_score FLOAT,
    score FLOAT
);
CREATE INDEX idx_answer_score_quiz_results_id ON answer_scores(quiz_result_id);

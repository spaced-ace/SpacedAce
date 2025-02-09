package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"math"
	"net/http"
	"slices"
	"spaced-ace/constants"
	"spaced-ace/context"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"spaced-ace/models/request"
	"spaced-ace/render"
	"spaced-ace/utils"
	"spaced-ace/views/components"
	"spaced-ace/views/forms"
	"spaced-ace/views/pages"
	"strconv"
)

var (
	difficultyOptions = []business.Option{
		{Name: "Easy", Value: "easy"},
		{Name: "Medium", Value: "medium"},
		{Name: "Hard", Value: "hard"},
	}

	statusOptions = []business.Option{
		{Name: "Due", Value: "due"},
		{Name: "Not Due", Value: "not-due"},
	}
)

func handleCreateQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	var requestForm request.CreateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	errors := map[string]string{}
	if requestForm.Title == "" {
		errors["title"] = "Title is required"
	}
	if requestForm.Description == "" {
		errors["description"] = "Description is required"
	}
	if len(errors) > 0 {
		return render.TemplRender(c, 200, forms.CreateQuizForm(requestForm, errors))
	}

	quizInfo, err := cc.ApiService.CreateQuiz(requestForm.Title, requestForm.Description)
	if err != nil {
		errors["other"] = "Error creating a quiz: " + err.Error()
		return render.TemplRender(c, 200, forms.CreateQuizForm(requestForm, errors))
	}

	c.Response().Header().Set("HX-Redirect", "/quizzes/"+quizInfo.Id+"/edit")
	return c.NoContent(http.StatusCreated)
}
func handleGenerateQuestionStart(c echo.Context) error {
	errors := map[string]string{}

	var requestForm request.GenerateQuestionForm
	if err := c.Bind(&requestForm); err != nil {
		errors["other"] = "Parsing error: " + err.Error()
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	questionType := requestForm.QuestionType
	if questionType != models.SingleChoiceQuestion && questionType != models.MultipleChoiceQuestion && questionType != models.TrueOrFalseQuestion {
		if questionType == "open-ended" {
			errors["other"] = fmt.Sprintf("Currently not suppored question type: '%s'.", questionType)
		} else {
			errors["other"] = fmt.Sprintf("Invalid question type: '%s'.", questionType)
		}
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	if requestForm.QuizId == "" {
		errors["other"] = "quizId is required"
	}
	if requestForm.Context == "" {
		errors["context"] = "Context is required"
	}
	if len(errors) > 0 {
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	return render.TemplRender(
		c,
		200,
		forms.GenerateQuestionForm(
			true,
			nil,
			requestForm,
			errors,
		),
	)
}
func handleGenerateQuestion(c echo.Context) error {
	errors := map[string]string{}
	cc := c.(*context.AppContext)

	var requestForm request.GenerateQuestionForm
	if err := c.Bind(&requestForm); err != nil {
		errors["other"] = "Parsing error: " + err.Error()
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	questionType := requestForm.QuestionType
	if questionType != models.SingleChoiceQuestion && questionType != models.MultipleChoiceQuestion && questionType != models.TrueOrFalseQuestion {
		if questionType == "open-ended" {
			errors["other"] = fmt.Sprintf("Currently not suppored question type: '%s'.", questionType)
		} else {
			errors["other"] = fmt.Sprintf("Invalid question type: '%s'.", questionType)
		}
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	if requestForm.QuizId == "" {
		errors["other"] = "quizId is required"
	}
	if requestForm.Context == "" {
		errors["context"] = "Context is required"
	}

	switch questionType {
	case models.SingleChoiceQuestion:
		{
			question, err := cc.ApiService.GenerateSingleChoiceQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				errors["other"] = "Error generating question: " + err.Error()
				return render.TemplRender(
					c,
					200,
					forms.GenerateQuestionForm(
						false,
						components.QuestionPlaceholderRemover(),
						requestForm,
						errors,
					),
				)
			}

			questionComponent := components.SingleChoiceQuestion(components.SingleChoiceQuestionProps{
				QuizSession:               nil,
				Question:                  question,
				Answer:                    nil,
				AllowDeleting:             true,
				ReplacePlaceholderWithOOB: true,
			})
			return render.TemplRender(
				c,
				200,
				forms.GenerateQuestionForm(
					false,
					questionComponent,
					requestForm,
					errors,
				),
			)
		}
	case models.MultipleChoiceQuestion:
		{
			question, err := cc.ApiService.GenerateMultipleChoiceQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				errors["other"] = "Error generating question: " + err.Error()
				return render.TemplRender(
					c,
					200,
					forms.GenerateQuestionForm(
						false,
						components.QuestionPlaceholderRemover(),
						requestForm,
						errors,
					),
				)
			}

			questionComponent := components.MultipleChoiceQuestion(components.MultipleChoiceQuestionProps{
				QuizSession:               nil,
				Question:                  question,
				Answer:                    nil,
				AllowDeleting:             true,
				ReplacePlaceholderWithOOB: true,
			})
			return render.TemplRender(
				c,
				200,
				forms.GenerateQuestionForm(
					false,
					questionComponent,
					requestForm,
					errors,
				),
			)
		}
	default:
		{
			question, err := cc.ApiService.GenerateTrueOrFalseQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				errors["other"] = "Error generating question: " + err.Error()
				return render.TemplRender(
					c,
					200,
					forms.GenerateQuestionForm(
						false,
						components.QuestionPlaceholderRemover(),
						requestForm,
						errors,
					),
				)
			}

			questionComponent := components.TrueOrFalseQuestion(components.TrueOrFalseQuestionProps{
				QuizSession:               nil,
				Question:                  question,
				Answer:                    nil,
				AllowDeleting:             true,
				ReplacePlaceholderWithOOB: true,
			})
			return render.TemplRender(
				c,
				200,
				forms.GenerateQuestionForm(
					false,
					questionComponent,
					requestForm,
					errors,
				),
			)
		}
	}
}

func handleAnswerQuestion(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizSessionId := c.Param("quizSessionId")
	if quizSessionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing quizSessionId")
	}

	var commonRequestForm request.CreateOrUpdateAnswerForm
	if err := c.Bind(&commonRequestForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid CreateOrUpdateRequestForm: ", err.Error())
	}

	switch commonRequestForm.AnswerType {
	case "single-choice":
		var requestForm request.CreateOrUpdateSingleChoiceAnswerForm
		if err := c.Bind(&requestForm); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid CreateOrUpdateSingleChoiceAnswerForm: ", err.Error())
		}

		_, err := cc.ApiService.CreateOrUpdateSingleChoiceAnswer(
			quizSessionId,
			requestForm.CreateOrUpdateAnswerForm.QuestionId,
			requestForm.Answer,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	case "multiple-choice":
		var requestForm request.CreateOrUpdateMultipleChoiceAnswerForm
		if err := c.Bind(&requestForm); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid CreateOrUpdateMultipleChoiceAnswerForm: ", err.Error())
		}

		_, err := cc.ApiService.CreateOrUpdateMultipleChoiceAnswer(
			quizSessionId,
			requestForm.CreateOrUpdateAnswerForm.QuestionId,
			requestForm.Answers,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	case "true-or-false":
		var requestForm request.CreateOrUpdateTrueOrFalseAnswerForm
		if err := c.Bind(&requestForm); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid CreateOrUpdateTrueOrFalseAnswerForm: ", err.Error())
		}

		_, err := cc.ApiService.CreateOrUpdateTrueOrFalseAnswer(
			quizSessionId,
			requestForm.CreateOrUpdateAnswerForm.QuestionId,
			requestForm.Answer,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	}

	return echo.NewHTTPError(http.StatusBadRequest, "invalid answerType: ", commonRequestForm.AnswerType)
}

func handleUpdateQuiz(c echo.Context) error {
	errors := map[string]string{}
	messages := map[string]string{}

	cc := c.(*context.AppContext)

	var requestForm request.UpdateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		errors["other"] = "Parsing error: " + err.Error()
		return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
	}

	if requestForm.QuizId == "" {
		errors["other"] = "quizId is required"
	}
	if requestForm.Title == "" {
		errors["title"] = "Title is required"
	}
	if requestForm.Description == "" {
		errors["description"] = "Description is required"
	}
	if len(errors) > 0 {
		return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
	}

	_, err := cc.ApiService.UpdateQuiz(requestForm.QuizId, requestForm.Title, requestForm.Description)
	if err != nil {
		errors["other"] = fmt.Sprintf("Error updating %s, error: %s", requestForm.QuizId, err.Error())
		return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
	}

	messages["successful"] = fmt.Sprintf("Succesfuly updated '%s'!", requestForm.Title)
	return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
}
func handleDeleteQuestion(c echo.Context) error {
	cc := c.(*context.AppContext)

	questionId := c.Param("questionId")
	if questionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Question ID is required")
	}

	questionType := c.QueryParam("type")
	if !utils.StringInArray(questionType, []string{models.SingleChoiceQuestion, models.MultipleChoiceQuestion, models.TrueOrFalseQuestion}) {
		return echo.NewHTTPError(400, "Invalid question type: "+questionType)
	}

	quizId := c.QueryParam("quizId")
	if quizId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	err := cc.ApiService.DeleteQuestion(questionType, quizId, questionId)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
func handleDeleteQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")
	if quizId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	if err := cc.ApiService.DeleteQuiz(quizId); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func handleSubmitQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizSessionId := c.Param("quizSessionId")
	if quizSessionId == "" {
		log.Default().Print("missing quizSessionId in url params\n")
		return echo.NewHTTPError(http.StatusBadRequest, "missing quizSessionId in url param")
	}

	quizSession, err := cc.ApiService.GetQuizSession(quizSessionId)
	if err != nil {
		log.Default().Printf("invalid quiz session ID `%s`\n", quizSessionId)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid quizSessionId: %s", quizSessionId))
	}

	var quizResult *business.QuizResult
	if quizSession.Finished {
		quizResult, err = cc.ApiService.GetQuizResult(quizSessionId)
	} else {
		quizResult, err = cc.ApiService.SubmitQuiz(quizSessionId)
	}
	if err != nil {
		return err
	}

	quiz, err := cc.ApiService.GetQuiz(quizSession.QuizId)
	if err != nil {
		return err
	}

	var answerLists *business.AnswerLists
	answers, err := cc.ApiService.GetAnswers(quizSession.Id)
	if err == nil {
		answerLists = answers
	}

	viewModel := pages.QuizResulPageViewModel{
		QuizSession: quizSession,
		Quiz:        quiz,
		AnswerLists: answerLists,
		QuizResult:  quizResult,
	}
	c.Response().Header().Set("HX-Replace-Url", fmt.Sprintf("/quiz-history/%s", quizSessionId))
	return render.TemplRender(c, 200, pages.QuizResultPage(viewModel))
}

func handleShowLearnListPopup(c echo.Context) error {
	cc := c.(*context.AppContext)

	learnList, err := cc.ApiService.GetLearnList()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	props := components.LearnListPopupProps{
		LearnList: *learnList,
	}
	return render.TemplRender(c, 200, components.LearnListPopup(props))
}
func handleAddQuizToLearnList(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizID := c.Param("quizID")
	if quizID == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "missing path param quizID")
	}

	learnList, err := cc.ApiService.AddQuizToLearnList(quizID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	props := components.LearnListPopupProps{
		LearnList: *learnList,
	}
	return render.TemplRender(c, 200, components.LearnListPopup(props))
}
func handleRemoveQuizFromLearnList(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizID := c.Param("quizID")
	if quizID == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "missing path param quizID")
	}

	learnList, err := cc.ApiService.RemoveQuizFromLearnList(quizID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	props := components.LearnListPopupProps{
		LearnList: *learnList,
	}
	return render.TemplRender(c, 200, components.LearnListPopup(props))
}

func handleGetReviewItemList(c echo.Context) error {
	cc := c.(*context.AppContext)

	quiz := c.QueryParam("quiz")

	difficulty := c.QueryParam("difficulty")
	if !slices.Contains([]string{"easy", "medium", "hard"}, difficulty) {
		difficulty = ""
	}

	selectedDifficulty := business.Option{}
	for _, option := range difficultyOptions {
		if option.Value == difficulty {
			selectedDifficulty = option
			break
		}
	}

	status := c.QueryParam("status")
	if !slices.Contains([]string{"due", "not-due"}, status) {
		status = ""
	}

	selectedStatus := business.Option{}
	for _, option := range statusOptions {
		if option.Value == status {
			selectedStatus = option
			break
		}
	}

	query := c.QueryParam("query")

	pageString := c.QueryParam("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}

	reviewItems, quizOptions, maxReviewItemCount, err := cc.ApiService.GetReviewItemListData(quiz, difficulty, status, query, page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("getting review item list data: %w", err))
	}

	selectedQuizOption := business.Option{}
	for _, option := range quizOptions {
		if option.Value == quiz {
			selectedQuizOption = option
			break
		}
	}

	pageCount := int(math.Ceil(float64(maxReviewItemCount) / float64(constants.REVIEW_ITEM_PAGE_SIZE)))
	pageOptions := make([]int, 0, 1)
	for i := 1; i <= pageCount; i++ {
		if page >= i && len(pageOptions) < 5 {
			pageOptions = append(pageOptions, i)
		}
	}

	previousPage := page - 1
	if previousPage < 1 {
		previousPage = -1
	}

	nextPage := page + 1
	if nextPage >= pageCount {
		nextPage = -1
	}

	props := components.ReviewItemListProps{
		SelectedQuizOption: selectedQuizOption,
		QuizOptions:        quizOptions,
		SelectedDifficulty: selectedDifficulty,
		DifficultyOptions:  difficultyOptions,
		SelectedStatus:     selectedStatus,
		StatusOptions:      statusOptions,
		Query:              query,
		ReviewItems:        reviewItems,
		PageOptions:        pageOptions,
		CurrentPage:        page,
		PreviousPage:       previousPage,
		NextPage:           nextPage,
	}
	return render.TemplRender(c, 200, components.ReviewItemList(props))
}
func handleSubmitReviewItemQuestion(c echo.Context) error {
	err := handleReviewItemQuestionSubmission(c)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/learn")
	return c.NoContent(http.StatusOK)
}
func handleSubmitReviewItemQuestionAndNext(c echo.Context) error {
	err := handleReviewItemQuestionSubmission(c)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/learn/review-all")
	return c.NoContent(http.StatusOK)
}

func handleReviewItemQuestionSubmission(c echo.Context) error {
	reviewItemID := c.Param("reviewItemID")

	answerForm := new(request.SubmitReviewItemQuestionForm)
	if err := c.Bind(answerForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("parsing review item question submission form: %w\n", err))
	}

	cc := c.(*context.AppContext)

	_, err := cc.ApiService.SubmitReviewItemQuestion(reviewItemID, *answerForm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("submitting review item question form: %w\n", err))
	}

	return nil
}
func handleQuizPreviewPopup(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")
	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "quizId is required")
	}

	// TODO
	_, err = cc.ApiService.HasQuizSession(cc.Session.User.Id, quizId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error checking open sessions: "+err.Error())
	}

	// TODO get all sessions and then find the open one

	quizSessions, err := cc.ApiService.GetQuizSessions(cc.Session.User.Id, quizId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error getting open sessions: "+err.Error())
	}

	var quizSession *business.QuizSession
	for _, session := range quizSessions {
		if !session.Finished {
			quizSession = session
			break
		}
	}

	props := components.PreviewQuizPopupProps{
		Quiz:        quiz,
		QuizSession: quizSession,
	}
	return render.TemplRender(c, 200, components.PreviewQuizPopup(props))
}
func handleClosePopup(c echo.Context) error {
	return c.NoContent(200)
}
func handleDrawerPopup(c echo.Context) error {
	return render.TemplRender(c, 200, components.NavbarDrawerPopup())
}
func handleQuizDrawerPopup(c echo.Context) error {
	cc := c.(*context.AppContext)

	props := components.QuizDrawerProps{
		Username: cc.Session.User.Name,
	}
	return render.TemplRender(c, 200, components.QuizDrawerPopup(props))
}

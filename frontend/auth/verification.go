package auth

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/render"
	"spaced-ace/views/components"
	"spaced-ace/views/pages"
)

func GetVerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return render.TemplRender(c, http.StatusBadRequest, pages.VerifyEmailPage("error", "Missing verification token"))
	}
	resp, err := http.Get(constants.BACKEND_URL + "/verify-email?token=" + token)
	if err != nil {
		log.Println("Error verifying email:", err)
		return render.TemplRender(c, http.StatusInternalServerError, pages.VerifyEmailPage("error", "Error connecting to verification service"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			errorResp.Message = "Verification failed"
		}
		return render.TemplRender(c, resp.StatusCode, pages.VerifyEmailPage("error", errorResp.Message))
	}

	return render.TemplRender(c, http.StatusOK, pages.VerifyEmailPage("success", ""))
}

// GetEmailVerificationNeeded renders the email verification needed page
func GetEmailVerificationNeeded(c echo.Context) error {
	email := c.QueryParam("email")
	viewModel := pages.EmailVerificationNeededViewModel{
		Email: email,
	}
	return render.TemplRender(c, http.StatusOK, pages.EmailVerificationNeededPage(viewModel))
}

func PostResendVerification(c echo.Context) error {
	var request struct {
		Email string `form:"email" json:"email"`
	}
	if err := c.Bind(&request); err != nil {
		return render.TemplRender(c, http.StatusBadRequest, components.VerificationFailed("Invalid request"))
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Println("Error marshaling resend request:", err)
		return render.TemplRender(c, http.StatusInternalServerError, components.VerificationFailed("Error processing request"))
	}

	resp, err := http.Post(constants.BACKEND_URL+"/resend-verification", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error resending verification email:", err)
		return render.TemplRender(c, http.StatusInternalServerError, components.VerificationFailed("Error connecting to verification service"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			errorResp.Message = "Failed to resend verification email"
		}
		return render.TemplRender(c, resp.StatusCode, components.VerificationFailed(errorResp.Message))
	}

	return render.TemplRender(c, http.StatusOK, components.VerificationEmailSent())
}

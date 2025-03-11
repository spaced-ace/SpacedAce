package auth

import (
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
)

var emailVerificationService *EmailVerificationService

type EmailVerificationService struct {
	client     *resend.Client
	fromEmail  string
	appBaseURL string
	mu         *sync.Mutex
}

// Initializes the email service
//
// unsafe to call concurrently
func InitEmailService() error {
	if emailVerificationService != nil {
		return nil
	}
	svc, err := newEmailVerificationService()
	if err != nil {
		return fmt.Errorf("Failed to initialize email client: %w", err)
	}
	emailVerificationService = svc
	return nil
}

// Returns the service or panics
func GetEmailVerificationService() *EmailVerificationService {
	if emailVerificationService == nil {
		panic("Email service was not initialized")
	}
	return emailVerificationService
}

func newEmailVerificationService() (*EmailVerificationService, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY environment variable not set")
	}

	client := resend.NewClient(apiKey)

	fromEmail := os.Getenv("EMAIL_FROM_ADDRESS")
	if fromEmail == "" {
		fromEmail = "verification@spacedace.hu"
	}

	appBaseURL := os.Getenv("APP_BASE_URL")
	if appBaseURL == "" {
		appBaseURL = "http://localhost" // Default for local development
	}

	return &EmailVerificationService{
		client:     client,
		fromEmail:  fromEmail,
		appBaseURL: appBaseURL,
		mu:         &sync.Mutex{},
	}, nil
}

func GenerateVerificationToken() string {
	return uuid.NewString()
}

func (s *EmailVerificationService) SendVerificationEmail(email, name, token string) error {
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", s.appBaseURL, token)

	s.mu.Lock()
	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{email},
		Subject: "Verify your SpacedAce account",
		Html: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2>Welcome to SpacedAce!</h2>
				<p>Hi %s,</p>
				<p>Thank you for signing up. Please verify your email address by clicking the button below:</p>
				<p style="text-align: center;">
					<a href="%s" style="display: inline-block; background-color: #4F46E5; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Verify Email</a>
				</p>
				<p>If you didn't create an account, you can safely ignore this email.</p>
				<p>Best regards,<br>The SpacedAce Team</p>
			</div>
		`, name, verificationLink),
	}

	_, err := s.client.Emails.Send(params)
	s.mu.Unlock()
	return err
}

package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/matcornic/hermes/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// EmailService defines the capabilities of our email provider
type EmailService interface {
	SendOTPEmail(to, name, otp string) error
}

type gmailService struct {
	srv *gmail.Service
}

// NewEmailService initializes the Gmail client. If token.json is missing,
// it will trigger the interactive web flow in the terminal.
func NewEmailService(ctx context.Context, clientID, clientSecret, tokenFile string) (EmailService, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{gmail.GmailSendScope, gmail.GmailReadonlyScope},
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob", // Use out-of-band for CLI apps
	}

	client := getClient(ctx, config, tokenFile)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %w", err)
	}

	return &gmailService{srv: srv}, nil
}

func (g *gmailService) SendOTPEmail(to, name, otp string) error {
	emailBody := fmt.Sprintf("From: ztas.global@gmail.com\r\n"+
		"To: %s\r\n"+
		"Subject: Your Z-TAS Registration OTP\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"utf-8\"\r\n"+
		"\r\n"+
		generateOTPUI(otp, name), to)

	var msg gmail.Message
	msg.Raw = base64.URLEncoding.EncodeToString([]byte(emailBody))

	_, err := g.srv.Users.Messages.Send("me", &msg).Do()
	if err != nil {
		return fmt.Errorf("failed to send OTP email to %s: %w", to, err)
	}

	return nil
}

func generateOTPUI(otpCode, name string) string {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Z-TAS",
			Link: "https://z-tas.com/",
			Logo: "https://z-tas.com/Assets/ZTAS-Text.webp",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				"Welcome to Z-TAS! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Z-TAS, here is your OTP:",
					InviteCode:   otpCode,
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		// Log this in production rather than panic, but returning empty string for safety
		fmt.Printf("Hermes HTML generation failed: %v\n", err)
		return ""
	}

	return emailBody
}

// --- Token Management Helpers ---

func getClient(ctx context.Context, config *oauth2.Config, tokFile string) *http.Client {
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(ctx, config)
		saveToken(tokFile, tok)
	}
	return config.Client(ctx, tok)
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("\n=== GOOGLE OAUTH REQUIRED ===\n")
	fmt.Printf("Go to the following link in your browser: \n\n%v\n\n", authURL)
	fmt.Printf("Type the authorization code here and press Enter: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		panic(fmt.Sprintf("Unable to read authorization code: %v", err))
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		panic(fmt.Sprintf("Unable to retrieve token from web: %v", err))
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		panic(fmt.Sprintf("Unable to cache oauth token: %v", err))
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

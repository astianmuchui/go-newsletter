package services

import (
	"fmt"
	"log"
	"os"

	"github.com/astianmuchui/go-newsletter/models"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var Store = session.New()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}


func SendEmails(subject string, msg string) error {
    var dbHandler models.UserFunctions = &models.Subscriber{}
    subscribers := dbHandler.GetSubscribers()

    dialerHost := os.Getenv("GMAIL_CLIENT")
    dialerPort := 587
    dialerUsername := os.Getenv("EMAIL_ADDRESS")
    dialerPassword := os.Getenv("GMAIL_APP_PASSWORD")

    dialer := gomail.NewDialer(dialerHost, dialerPort, dialerUsername, dialerPassword)
    s, err := dialer.Dial()
    if err != nil {
        return fmt.Errorf("failed to connect to SMTP server: %w", err)
    }
    defer s.Close()

    m := gomail.NewMessage()
    for _, r := range subscribers {
        m.SetHeader("From", dialerUsername)
        m.SetAddressHeader("To", r.Email, r.Name)
        m.SetHeader("Subject", subject)
        m.SetBody("text/html", fmt.Sprintf("Hello %s <br> %s", r.Name, msg))

        if err := gomail.Send(s, m); err != nil {
            log.Printf("Could not send email to %q: %v", r.Email, err)
            continue
        }

        log.Printf("Email sent to %q", r.Email)
        m.Reset()
    }
    return nil
}
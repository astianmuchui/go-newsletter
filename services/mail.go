package services

import (
	"fmt"
	"log"
	"os"

	"github.com/astianmuchui/go-newsletter/models"
	"github.com/gofiber/fiber/v2"
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
func SendEmails(subject string, msg string, c *fiber.Ctx) {
	var db_handler models.UserFunctions = &models.Subscriber{}
	var data = db_handler.GetSubscribers()

	dialerHost := os.Getenv("GMAIL_CLIENT")
	dialerPort := 587
	dialerUsername := os.Getenv("EMAIL_ADDRESS")
	dialerPassword := os.Getenv("GMAIL_APP_PASSWORD")

	log.Println("Dialer Host:", dialerHost)
	log.Println("Dialer Port:", dialerPort)
	log.Println("Dialer Username:", dialerUsername)
	log.Println("Dialer Password:", dialerPassword)

	dialer := gomail.NewDialer(dialerHost, dialerPort, dialerUsername, dialerPassword)

	s, err := dialer.Dial()
	sess, _ := Store.Get(c)
	if sess.Get("errors") == nil {
		sess.Set("errors", []error{})
	}


	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
	for _, r := range data {
		m.SetHeader("From", dialerUsername)
		m.SetAddressHeader("To", r.Email, r.Name)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", fmt.Sprintf("Hello %s <br> %s", r.Name, msg))

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", r.Email, err)
			sess.Set("errors", append(sess.Get("errors").([]error), fmt.Errorf("could not send email to %q: %v", r.Email, err)))
			sess.Save()
		}

		sess.Set("mail_success", "Process over !")
		sess.Save()
		m.Reset()
	}

}

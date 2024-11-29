package handlers

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/astianmuchui/go-newsletter/models"
	"github.com/astianmuchui/go-newsletter/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Message struct {
	Subject string
	Message string
}

var Store = session.New()

func HomeHandler(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)
	var err, success string

	if sess.Get("err") != nil {
		err = sess.Get("err").(string)
		sess.Delete("err")
		sess.Save()
	}

	if sess.Get("mail_success") != nil {
		success = sess.Get("mail_success").(string)
		sess.Delete("mail_success")
		sess.Save()
	}

	var db_handler models.UserFunctions = &models.Subscriber{}

	data := db_handler.GetSubscribers()

	return c.Render("index", fiber.Map{
		"error":       err,
		"success":     success,
		"subscribers": data,
	})
}

func SendHandler(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)
	var errors []error
	var mailSuccess string

	if sess.Get("errors") != nil {
		errors = sess.Get("errors").([]error)
		sess.Delete("errors")
		sess.Save()
	}

	if sess.Get("mail_success") != nil {
		mailSuccess = sess.Get("mail_success").(string)
		sess.Delete("mail_success")
		sess.Save()
	}


	return c.Render("send_message", fiber.Map{
		"errors":       errors,
		"success": mailSuccess,
	})
}

func SubscribeHandler(c *fiber.Ctx) error {
	payload := new(models.Subscriber)

	if err := c.BodyParser(payload); err != nil {

		return c.Status(fiber.StatusInternalServerError).SendString("A server error occured")

	} else {

		var db_handler models.UserFunctions = &models.Subscriber{}
		sess, _ := Store.Get(c)

		if db_handler.Exists(payload.Email) {
			sess.Set("err", "Subscriber already exists")
			sess.Save()
			return c.Redirect("/")
		}
		var created bool = db_handler.CreateSubscriber(*payload)

		if created {
			sess.Set("success", "Subscriber added")
			sess.Save()

			return c.Redirect("/")
		} else {
			sess.Set("err", "Subscriber not added")
			sess.Save()

			return c.Redirect("/")
		}
	}
}

func SendEmailHandler(c *fiber.Ctx) error {
    payload := new(Message)
    sess, err := Store.Get(c)
    if err != nil {
        log.Error("Failed to get session:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to initialize session")
    }

    if err := c.BodyParser(payload); err != nil {
        log.Error("Failed to parse payload:", err)
        return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
    }

    sess.Set("feedback", "Emails being sent")
    if err := sess.Save(); err != nil {
        log.Error("Failed to save session:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
    }

    go func(subject, message string) {
        if err := services.SendEmails(subject, message); err != nil {
            log.Error("Failed to send emails:", err)
        }
    }(payload.Subject, payload.Message)

    return c.Redirect("/send")
}
package email

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

func SendEmail(payerEmail string, payeeDisplayName string, amount float64, timeSubmitted time.Time) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("EMAIL"),
		os.Getenv("EMAIL_PASS"),
		"smtp.gmail.com",
	)

	to := []string{payerEmail}
	msg := []byte(
		fmt.Sprintf("To: %s\r\n"+
			"Subject: Results from Poker Night %s\r\n"+
			"\r\n"+
			"You need to transfer %s %.2f$.\r\n",
			payerEmail,
			timeSubmitted.Format("2006-01-02 15:04:05"),
			payeeDisplayName,
			amount,
		),
	)
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("EMAIL"),
		to,
		msg,
	)

	return err
}

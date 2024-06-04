package email

import (
	"backend/internal/config"
	"backend/internal/constants"
	"backend/internal/models"
	"log"
	"time"

	"fmt"

	"github.com/google/uuid"
)

type EmailWorkerQueue struct {
	queue chan models.PaymentDetails
}

var emailWorkerQueue *EmailWorkerQueue

func CreateEmailChannel() {
	emailWorkerQueue = &EmailWorkerQueue{
		queue: make(chan models.PaymentDetails),
	}
	go consumeEmailMessages()
}

func QueueEmail(details models.PaymentDetails) {
	emailWorkerQueue.queue <- details
}

func consumeEmailMessages() {
	for paymentDetail := range emailWorkerQueue.queue {
		var payer, payee models.User
		r1 := config.DB.Find(&payer, "id = ?", paymentDetail.PayerID)
		r2 := config.DB.Find(&payee, "id = ?", paymentDetail.PayeeID)
		if r1.Error != nil || r2.Error != nil {
			updatePaymentDetails(paymentDetail.ID, constants.PAYMENT_DETAILS_FAILED)
			log.Println("NULLS")
		}

		err := SendEmail(payer.Email, payee.DisplayName, paymentDetail.Amount, paymentDetail.TimeSubmitted)

		// update the status on payment detail
		if err != nil {
			updatePaymentDetails(paymentDetail.ID, constants.PAYMENT_DETAILS_FAILED)
		} else {
			updatePaymentDetails(paymentDetail.ID, constants.PAYMENT_DETAILS_EMAIL_SENT)
		}
		fmt.Println(err)
	}
}

func updatePaymentDetails(id uuid.UUID, status string) {
	var paymentDetails models.PaymentDetails
	result := config.DB.First(&paymentDetails, "id = ?", id)
	if result.Error != nil {
		log.Fatal("payment not found")
		return
	}

	result = config.DB.Model(&paymentDetails).Updates(models.PaymentDetails{
		Status:        status,
		TimeCompleted: time.Now(),
	})
	if result.Error != nil {
		log.Fatal("payment not found")
		return
	}
}

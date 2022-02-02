package external

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct{}
type SendDTO struct {
	CryptoId   string
	Email      string
	Percentage float64
	Value      float64
}

func (sg *SendGrid) Send(payload SendDTO) error {
	from := mail.NewEmail("Crypto Tracing", "contato@piperapp.co")
	subject := "CryptoTracking - " + fmt.Sprint(payload.CryptoId)
	to := mail.NewEmail(payload.Email, payload.Email)
	plainTextContent := payload.CryptoId
	htmlContent := "<p>increased in <b>" + fmt.Sprint(payload.Percentage) + "% </b> with value <b>" + fmt.Sprint(payload.Value) + "</b></p>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}

package usecases

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	"github.com/DKhorkov/medods/internal/config"
	gomail "gopkg.in/gomail.v2"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(
		rand.NewSource(
			time.Now().UnixNano(),
		),
	)

	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(bytes)
}

func sendEmail(
	subject string,
	body string,
	emailsTo []string,
	smtpConfig config.SMTPConfig,
	logger *slog.Logger,
) {
	message := gomail.NewMessage()
	message.SetHeader("From", smtpConfig.Login)
	message.SetHeader("To", emailsTo...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	smtpClient := gomail.NewDialer(
		smtpConfig.Host,
		smtpConfig.Port,
		smtpConfig.Login,
		smtpConfig.Password,
	)

	if err := smtpClient.DialAndSend(message); err != nil {
		logger.Error(
			"Failed to send email",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
	}

	// address := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)
	// message := []byte(subject + body)
	// auth := smtp.PlainAuth("", smtpConfig.Login, smtpConfig.Password, smtpConfig.Host)
	//err := smtp.SendMail(address, auth, smtpConfig.Login, emailsTo, message)
	//if err != nil {
	//	logger.Error(
	//		"Failed to send email",
	//		"Traceback",
	//		logging.GetLogTraceback(),
	//		"Error",
	//		err,
	//	)
	//}
}

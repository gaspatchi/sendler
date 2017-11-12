package sendqueue

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

type authorization struct {
	smtpServer   string
	smtpUser     string
	smtpPassword string
	smtpName     string
}

var auth *authorization = &authorization{}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Panic("Ошибка при чтении конфига")
		panic(err)
	}
	auth.smtpServer = viper.GetString("email.server")
	auth.smtpUser = viper.GetString("email.user")
	auth.smtpPassword = viper.GetString("email.password")
	auth.smtpName = viper.GetString("email.name")

}

func SendMail(subject string, address string, name string, body string) error {
	message := gomail.NewMessage()
	message.SetAddressHeader("From", auth.smtpUser, auth.smtpName)
	message.SetAddressHeader("To", address, name)
	message.SetHeader("Bcc", auth.smtpUser)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)
	dialer := gomail.NewDialer(auth.smtpServer, 465, auth.smtpUser, auth.smtpPassword)
	err := dialer.DialAndSend(message)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendMail"}).Error("Ошибка при отправке письма  ", err)
		return err
	}
	return nil
}

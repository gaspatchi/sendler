package sendqueue

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"encoding/json"

	"time"

	"../../lib/inforer"
	"../../lib/schedler"
	"github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
)

type VerifyRegistrationStruct struct {
	Email     string
	Firstname string
	Lastname  string
	Token     string
	Template  bytes.Buffer
}

type ResetPasswordStruct struct {
	Email     string
	Firstname string
	Lastname  string
	Password  string
	Token     string
	Template  bytes.Buffer
}

type SetNumberStruct struct {
	Email     string
	Firstname string
	Lastname  string
	Template  bytes.Buffer
}

type AddMoneyStruct struct {
	Email     string
	Firstname string
	Lastname  string
	Template  bytes.Buffer
}

type SendClaimStruct struct {
	Type      string
	Email     string
	Number    string
	Firstname string
	Lastname  string
	Text      string
	Template  bytes.Buffer
}

type SendMessageStruct struct {
	Email     string
	Firstname string
	Lastname  string
	Text      string
	Template  bytes.Buffer
}

type TeacherSchedule struct {
	Info     inforer.TeacherInfo      `json:"info"`
	Schedule schedler.TeacherSchedule `json:"schedule"`
}

type sendSchedule struct {
	Email     string
	Firstname string
	Lastname  string
	Date      string
	Text      string
	Groups    []struct {
		Info struct {
			Result struct {
				Group struct {
					Group string `json:"group"`
				} `json:"group"`
			} `json:"result"`
		} `json:"info"`
		Schedule struct {
			Schedule []struct {
				Teacher struct {
					Lastname   string `json:"lastname"`
					Firstname  string `json:"firstname"`
					Patronymic string `json:"patronymic"`
				} `json:"teacher"`
				Lesson struct {
					Lesson string `json:"lesson"`
				} `json:"lesson"`
				Index   int `json:"index"`
				Cabinet struct {
					Cabinet string `json:"cabinet"`
				} `json:"cabinet"`
			} `json:"schedule"`
		} `json:"schedule"`
	} `json:"groups"`
	Teachers []struct {
		Info struct {
			Result struct {
				Teacher struct {
					Lastname   string `json:"lastname"`
					Firstname  string `json:"firstname"`
					Patronymic string `json:"patronymic"`
				} `json:"teacher"`
			} `json:"result"`
		} `json:"info"`
		Schedule struct {
			Schedule []struct {
				Lesson struct {
					Lesson string `json:"lesson"`
				} `json:"lesson"`
				Group struct {
					Group string `json:"group"`
				} `json:"group"`
				Date    string `json:"date"`
				Cabinet struct {
					Cabinet string `json:"cabinet"`
				} `json:"cabinet"`
				Index int `json:"index"`
			} `json:"schedule"`
		} `json:"schedule"`
	} `json:"teachers"`
	Template bytes.Buffer
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func VerifyRegistration(connection *tarantool.Connection, taskdata SendTask) error {
	var registerScheme VerifyRegistrationStruct
	registerScheme.Token = taskdata.Data.Token
	email, _, err := taskdata.GetAddresses(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "VerifyRegistration"}).Error("Ошибка при получении адресо  ", err)
		return err
	}
	firstname, lastname, err := taskdata.GetInitials(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "VerifyRegistration"}).Error("Ошибка при получении инициало  ", err)
		return err
	}
	model, err := GetTemplate(connection, "verifyRegistration")
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "VerifyRegistration"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	registerScheme.Email = email
	registerScheme.Firstname = firstname
	registerScheme.Lastname = lastname
	emaiTemplate, _ := template.New("VerifyRegistration").Parse(model)
	err = emaiTemplate.Execute(&registerScheme.Template, registerScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "VerifyRegistration"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail("👤 Активация аккаунта", registerScheme.Email, fmt.Sprintf("%s %s", registerScheme.Firstname, registerScheme.Lastname), registerScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "VerifyRegistration"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

func ResetPassword(connection *tarantool.Connection, taskdata SendTask) error {
	var resetScheme ResetPasswordStruct
	resetScheme.Token = taskdata.Data.Token
	resetScheme.Password = taskdata.Data.Password
	email, _, err := taskdata.GetAddresses(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "ResetPassword"}).Error("Ошибка при получении адресо  ", err)
		return err
	}
	resetScheme.Email = email
	firstname, lastname, err := taskdata.GetInitials(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "ResetPassword"}).Error("Ошибка при получении инициало  ", err)
		return err
	}
	resetScheme.Firstname = firstname
	resetScheme.Lastname = lastname
	model, err := GetTemplate(connection, "resetPassword")
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "ResetPassword"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	emaiTemplate, _ := template.New("resetPassword").Parse(model)
	err = emaiTemplate.Execute(&resetScheme.Template, resetScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "ResetPassword"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail("🔑 Востановление пароля", resetScheme.Email, fmt.Sprintf("%s %s", resetScheme.Firstname, resetScheme.Lastname), resetScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "ResetPassword"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

func SetNumber(connection *tarantool.Connection, taskdata SendTask) error {
	var setNumberScheme SetNumberStruct
	email, _, err := taskdata.GetAddresses(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SetNumber"}).Error("Ошибка при получении адресо  ", err)
		return err
	}
	setNumberScheme.Email = email
	firstname, lastname, err := taskdata.GetInitials(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SetNumber"}).Error("Ошибка при получении инициало  ", err)
		return err
	}
	setNumberScheme.Firstname = firstname
	setNumberScheme.Lastname = lastname
	model, err := GetTemplate(connection, "setNumber")
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SetNumber"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	emaiTemplate, _ := template.New("setNumber").Parse(model)
	err = emaiTemplate.Execute(&setNumberScheme.Template, setNumberScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SetNumber"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail("📵 Укажите номер", setNumberScheme.Email, fmt.Sprintf("%s %s", setNumberScheme.Firstname, setNumberScheme.Lastname), setNumberScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SetNumber"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

func AddMoney(connection *tarantool.Connection, taskdata SendTask) error {
	var addMoneyScheme AddMoneyStruct
	email, _, err := taskdata.GetAddresses(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "AddMoney"}).Error("Ошибка при получении адресо  ", err)
		return err
	}
	addMoneyScheme.Email = email
	firstname, lastname, err := taskdata.GetInitials(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "AddMoney"}).Error("Ошибка при получении инициало  ", err)
		return err
	}
	addMoneyScheme.Firstname = firstname
	addMoneyScheme.Lastname = lastname
	model, err := GetTemplate(connection, "addMoney")
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "AddMoney"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	emaiTemplate, _ := template.New("addMoney").Parse(model)
	err = emaiTemplate.Execute(&addMoneyScheme.Template, addMoneyScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "AddMoney"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail("💰 Пополните счет", addMoneyScheme.Email, fmt.Sprintf("%s %s", addMoneyScheme.Firstname, addMoneyScheme.Lastname), addMoneyScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "AddMoney"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

func SendClaim(connection *tarantool.Connection, taskdata SendTask) error {
	var sendClaimScheme SendClaimStruct
	sendClaimScheme.Type = taskdata.Data.Type
	sendClaimScheme.Email = taskdata.Data.Email
	sendClaimScheme.Number = taskdata.Data.Number
	sendClaimScheme.Firstname = taskdata.Data.Firstname
	sendClaimScheme.Lastname = taskdata.Data.Lastname
	sendClaimScheme.Text = taskdata.Data.Text
	model, err := GetTemplate(connection, "sendClaim")
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendClaim"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	emaiTemplate, _ := template.New("sendClaim").Parse(model)
	err = emaiTemplate.Execute(&sendClaimScheme.Template, sendClaimScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendClaim"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail(sendClaimScheme.Type, taskdata.Data.Address, fmt.Sprintf("%s %s", sendClaimScheme.Firstname, sendClaimScheme.Lastname), sendClaimScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendClaim"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

func SendMessage(connection *tarantool.Connection, taskdata SendTask) error {
	var sendMessageScheme SendMessageStruct
	email, _, err := taskdata.GetAddresses(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendMessage"}).Error("Ошибка при получении адресо  ", err)
		return err
	}
	sendMessageScheme.Email = email
	firstname, lastname, err := taskdata.GetInitials(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendMessage"}).Error("Ошибка при получении инициало  ", err)
		return err
	}
	sendMessageScheme.Firstname = firstname
	sendMessageScheme.Lastname = lastname
	sendMessageScheme.Text = taskdata.Data.Text
	model, err := GetTemplate(connection, "sendMessage")
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendMessage"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	emaiTemplate, _ := template.New("sendMessage").Parse(model)
	err = emaiTemplate.Execute(&sendMessageScheme.Template, sendMessageScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendMessage"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail("⚠️ Информация", sendMessageScheme.Email, fmt.Sprintf("%s %s", sendMessageScheme.Firstname, sendMessageScheme.Lastname), sendMessageScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendMessage"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

func SendNewSchedule(connection *tarantool.Connection, taskdata SendTask) error {
	var sendScheduleScheme sendSchedule
	scheduleTime, _ := time.Parse("2006-01-02", taskdata.Date)
	sendScheduleScheme.Date = fmt.Sprintf("%d.%d.%d", scheduleTime.Day(), int(scheduleTime.Month()), scheduleTime.Year())
	schedule, err := GetSchedule(connection, taskdata.Data.UserID, taskdata.Direction, taskdata.Date, taskdata.Template)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendNewSchedule"}).Error("Ошибка при получении расписания")
		return err
	}
	json.Unmarshal([]byte(schedule), &sendScheduleScheme)
	email, _, err := taskdata.GetAddresses(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendNewSchedule"}).Error("Ошибка при получении адресо  ", err)
		return err
	}
	sendScheduleScheme.Email = email
	firstname, lastname, err := taskdata.GetInitials(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendNewSchedule"}).Error("Ошибка при получении инициало  ", err)
		return err
	}
	sendScheduleScheme.Email = email
	sendScheduleScheme.Firstname = firstname
	sendScheduleScheme.Lastname = lastname
	model, err := GetTemplate(connection, taskdata.Template)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendNewSchedule"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	emaiTemplate, _ := template.New(taskdata.Template).Parse(model)
	err = emaiTemplate.Execute(&sendScheduleScheme.Template, sendScheduleScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendNewSchedule"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail(fmt.Sprintf("📅 Расписание занятий на %s", sendScheduleScheme.Date), sendScheduleScheme.Email, fmt.Sprintf("%s %s", sendScheduleScheme.Firstname, sendScheduleScheme.Lastname), sendScheduleScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendNewSchedule"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

func SendUpdateSchedule(connection *tarantool.Connection, taskdata SendTask) error {
	var sendScheduleScheme sendSchedule
	scheduleTime, _ := time.Parse("2006-01-02", taskdata.Date)
	sendScheduleScheme.Date = fmt.Sprintf("%d.%d.%d", scheduleTime.Day(), int(scheduleTime.Month()), scheduleTime.Year())
	schedule, err := GetSchedule(connection, taskdata.Data.UserID, taskdata.Direction, taskdata.Date, taskdata.Template)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendUpdateSchedule"}).Error("Ошибка при получении расписания")
		return err
	}
	json.Unmarshal([]byte(schedule), &sendScheduleScheme)
	email, _, err := taskdata.GetAddresses(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendUpdateSchedule"}).Error("Ошибка при получении адресо  ", err)
		return err
	}
	sendScheduleScheme.Email = email
	firstname, lastname, err := taskdata.GetInitials(connection)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendUpdateSchedule"}).Error("Ошибка при получении инициало  ", err)
		return err
	}
	sendScheduleScheme.Email = email
	sendScheduleScheme.Firstname = firstname
	sendScheduleScheme.Lastname = lastname
	model, err := GetTemplate(connection, taskdata.Template)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendUpdateSchedule"}).Error("Ошибка при получении шаблон  ", err)
		return err
	}
	emaiTemplate, _ := template.New(taskdata.Template).Parse(model)
	err = emaiTemplate.Execute(&sendScheduleScheme.Template, sendScheduleScheme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendUpdateSchedule"}).Error("Ошибка при формировании письм  ", err)
		return err
	}
	err = SendMail(fmt.Sprintf("🔃 Обновленное расписание занятий на %s", sendScheduleScheme.Date), sendScheduleScheme.Email, fmt.Sprintf("%s %s", sendScheduleScheme.Firstname, sendScheduleScheme.Lastname), sendScheduleScheme.Template.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send", "section": "SendUpdateSchedule"}).Error("Ошибка при отправке письм  ", err)
		return err
	}
	return nil
}

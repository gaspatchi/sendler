package sendqueue

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	tarantool "github.com/tarantool/go-tarantool"
)

type SendTask struct {
	Direction string `json:"direction"`
	Template  string `json:"template"`
	Date      string `json:"date"`
	Data      struct {
		Try       int    `json:"try"`
		UserID    string `json:"user_id"`
		Password  string `json:"password"`
		Token     string `json:"token"`
		Address   string `json:"address"`
		Type      string `json:"type"`
		Email     string `json:"email"`
		Number    string `json:"number"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Text      string `json:"text"`
	} `json:"data"`
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func (task *SendTask) Unmarshal(data []byte, point interface{}) (err error) {
	return json.Unmarshal(data, point)
}

func (task *SendTask) Marshal() (string, error) {
	result, err := json.Marshal(task)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при формировании JSON  ", err)
		return "", err
	}
	return string(result), nil
}

func (task *SendTask) GetAddresses(connection *tarantool.Connection) (email string, number string, err error) {
	response, err := connection.Call("getAddresses", []interface{}{task.Data.UserID})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении адресов  ", err)
		return "", "", err
	}
	if response.Data[0].([]interface{})[0].(bool) == true {
		email = response.Data[1].([]interface{})[0].(string)
		number = response.Data[2].([]interface{})[0].(string)
		return email, number, nil
	}
	logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении адресов  ", err)
	return "", "", errors.New(response.Data[1].([]interface{})[0].(string))
}

func (task *SendTask) GetInitials(connection *tarantool.Connection) (firstname string, lastname string, err error) {
	response, err := connection.Call("getInitials", []interface{}{task.Data.UserID})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении инициалов  ", err)
		return "", "", err
	}
	if response.Data[0].([]interface{})[0].(bool) == true {
		firstname = response.Data[1].([]interface{})[0].(string)
		lastname = response.Data[2].([]interface{})[0].(string)
		return firstname, lastname, nil
	}
	logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении инициалов  ", err)
	return "", "", errors.New(response.Data[1].([]interface{})[0].(string))
}

func GetTemplate(connection *tarantool.Connection, id string) (model string, err error) {
	response, err := connection.Call("getTemplate", []interface{}{id})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении шаблона письма  ", err)
		return "", err
	}
	if response.Data[0].([]interface{})[0].(bool) == true {
		return response.Data[1].([]interface{})[0].(string), err
	}
	logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении шаблона письма  ", err)
	return "", errors.New(response.Data[1].([]interface{})[0].(string))
}

func GetSchedule(connection *tarantool.Connection, userId string, direction string, date string, template string) (schedule string, err error) {
	response, err := connection.Call("selectSchedule", []interface{}{userId, direction, date, template})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении расписания  ", err)
		return "", err
	}
	if response.Data[0].([]interface{})[0].(bool) == true {
		return response.Data[1].([]interface{})[0].(string), err
	}
	logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении расписания  ", err)
	return "", errors.New(response.Data[1].([]interface{})[0].(string))
}

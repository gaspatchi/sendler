package fetchqueue

import (
	"encoding/json"
	"os"

	"errors"

	"../../lib/inforer"
	"../../lib/schedler"
	"github.com/sirupsen/logrus"
	tarantool "github.com/tarantool/go-tarantool"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

type TeacherSchedule struct {
	Info     inforer.TeacherInfo      `json:"info"`
	Schedule schedler.TeacherSchedule `json:"schedule"`
}

type GroupSchedule struct {
	Info     inforer.GroupInfo      `json:"info"`
	Schedule schedler.GroupSchedule `json:"schedule"`
}

type FetchTask struct {
	Data struct {
		Try    int    `json:"try"`
		UserID string `json:"user_id"`
	} `json:"data"`
	ID          int    `json:"id"`
	Destination string `json:"destination"`
	Action      string `json:"action"`
	Type        string `json:"type"`
	Date        string `json:"date"`
}

func (task *FetchTask) Unmarshal(data []byte, point interface{}) (err error) {
	return json.Unmarshal(data, point)
}

func (task *FetchTask) Marshal() (string, error) {
	result, err := json.Marshal(task)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Ошибка при формировании JSON  ", err)
		return "", err
	}
	return string(result), nil
}

func (schedule *TeacherSchedule) Post(connection *tarantool.Connection, task FetchTask) error {
	jsonSchedule, err := json.Marshal(schedule)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch", "section": "Post"}).Error("Невозможно сформировать JSON  ", err)
		return err
	}

	response, err := connection.Call("formatSchedule", []interface{}{task.Data.UserID, task.Date, task.Action, task.Destination, task.Type, jsonSchedule})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch", "section": "Post"}).Error("Невозможно отправить расписание по преподавателю  ", err)
		return err
	}
	if response.Data[0].([]interface{})[0].(bool) == false {
		logrus.WithFields(logrus.Fields{"module": "Fetch", "section": "Post"}).Error("Невозможно отправить расписание по преподавателю  ", err)
		return errors.New("Невозможно добавить расписание")
	}
	return nil
}

func (schedule *GroupSchedule) Post(connection *tarantool.Connection, task FetchTask) error {
	jsonSchedule, err := json.Marshal(schedule)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch", "section": "Post"}).Error("Невозможно сформировать JSON  ", err)
		return err
	}
	response, err := connection.Call("formatSchedule", []interface{}{task.Data.UserID, task.Date, task.Action, task.Destination, task.Type, jsonSchedule})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch", "section": "Post"}).Error("Невозможно отправить расписание по группе  ", err)
		return err
	}
	if response.Data[0].([]interface{})[0].(bool) == false {
		logrus.WithFields(logrus.Fields{"module": "Fetch", "section": "Post"}).Error("Невозможно отправить расписание по группе  ", err)
		return errors.New("Невозможно добавить расписание")
	}
	return nil
}

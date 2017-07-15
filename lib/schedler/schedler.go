package schedler

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/resty.v0"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

var endpoint string

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
	client, err := api.NewClient(api.DefaultConfig())
	catalog := client.Catalog()
	response, _, err := catalog.Service("schedler", "", &api.QueryOptions{Datacenter: "dc1"})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Schedler"}).Panic("Consul недоступен")
		panic(err)
	}
	if len(response) == 0 {
		logrus.WithFields(logrus.Fields{"module": "Schedler"}).Panic("Schedler не зарегистрирован")
		panic(err)
	} else {
		endpoint = fmt.Sprintf("%s:%d", response[0].ServiceAddress, response[0].ServicePort)
	}
}

func GetTeacher(id int, date string) (schedule TeacherSchedule, err error) {
	var teacherSchedule TeacherSchedule
	response, err := resty.R().SetResult(&teacherSchedule).Get(fmt.Sprintf("http://%s/teacher/%d/%s", endpoint, id, date))
	if response.StatusCode() == 200 {
		return teacherSchedule, nil
	}
	logrus.WithFields(logrus.Fields{"module": "Schedler"}).Error("Невозможно получить расписание по преподавателю")
	return teacherSchedule, errors.New("При получении расписания по преподавателю произошла ошибка")
}

func GetGroup(id int, date string) (schedule GroupSchedule, err error) {
	var groupSchedule GroupSchedule
	response, err := resty.R().SetResult(&groupSchedule).Get(fmt.Sprintf("http://%s/group/%d/%s", endpoint, id, date))
	if response.StatusCode() == 200 {
		return groupSchedule, nil
	}
	logrus.WithFields(logrus.Fields{"module": "Schedler"}).Error("Невозможно получить расписание по группе")
	return groupSchedule, errors.New("При получении расписания по группе произошла ошибка")
}

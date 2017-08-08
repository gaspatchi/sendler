package inforer

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v0"
)

var endpoint string

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
	consulconfig := api.DefaultConfig()
	consulconfig.Address = fmt.Sprintf("%s:%d", viper.GetString("consul.address"), viper.GetInt("consul.port"))
	client, err := api.NewClient(consulconfig)
	catalog := client.Catalog()
	response, _, err := catalog.Service("inforer", "", &api.QueryOptions{Datacenter: "dc1"})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Inforer"}).Panic("Consul недоступен")
		panic(err)
	}
	if len(response) == 0 {
		logrus.WithFields(logrus.Fields{"module": "Inforer"}).Panic("Inforer не зарегистрирован")
		panic(err)
	} else {
		endpoint = fmt.Sprintf("%s:%d", response[0].ServiceAddress, response[0].ServicePort)
	}
}

func GetTeacher(id int) (info TeacherInfo, err error) {
	var teacherInfo TeacherInfo
	response, err := resty.R().SetResult(&teacherInfo).Get(fmt.Sprintf("http://%s/teacher/%d", endpoint, id))
	if response.StatusCode() == 200 {
		return teacherInfo, nil
	}
	logrus.WithFields(logrus.Fields{"module": "Inforer"}).Error("Невозможно получить информацию о преподавателе")
	return teacherInfo, errors.New("При получении информации о преподавателе произошла ошибка")
}

func GetGroup(id int) (info GroupInfo, err error) {
	var groupInfo GroupInfo
	response, err := resty.R().SetResult(&groupInfo).Get(fmt.Sprintf("http://%s/group/%d", endpoint, id))
	if response.StatusCode() == 200 {
		return groupInfo, nil
	}
	logrus.WithFields(logrus.Fields{"module": "Inforer"}).Error("Невозможно получить информацию о группе")
	return groupInfo, errors.New("При получении информации о группе произошла ошибка")
}

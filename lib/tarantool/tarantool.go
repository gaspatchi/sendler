package tarantool

import (
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/queue"
)

type instance struct {
	Connection *tarantool.Connection
	FetchQueue queue.Queue
	SendQueue  queue.Queue
	endpoint   string
}

var trntl *instance = &instance{}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Ошибка при чтении конфига")
		panic(err)
	}
	consulconfig := api.DefaultConfig()
	consulconfig.Address = fmt.Sprintf("%s:%d", viper.GetString("consul.address"), viper.GetInt("consul.port"))
	client, err := api.NewClient(consulconfig)
	catalog := client.Catalog()
	response, _, err := catalog.Service("tarantool", "", &api.QueryOptions{Datacenter: "dc1"})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Consul недоступен")
		panic(err)
	}
	if len(response) == 0 {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Tarantool не зарегистрирован")
		panic(err)
	} else {
		trntl.endpoint = fmt.Sprintf("%s:%d", response[0].ServiceAddress, response[0].ServicePort)
	}
	setup(viper.GetString("tarantool.user"), viper.GetString("tarantool.password"))
}

func setup(user string, password string) {
	opts := tarantool.Opts{User: user, Pass: password, MaxReconnects: 5, Reconnect: 1}
	conn, err := tarantool.Connect(trntl.endpoint, opts)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Ошибка при подключении к Tarantool")
		panic(err)
	}
	trntl.Connection = conn
	trntl.FetchQueue = queue.New(conn, "fetch_queue")
	trntl.SendQueue = queue.New(conn, "send_queue")
}

func GetTarantool() *tarantool.Connection {
	return trntl.Connection
}

func GetFetchQueue() queue.Queue {
	return trntl.FetchQueue
}

func GetSendQueue() queue.Queue {
	return trntl.SendQueue
}

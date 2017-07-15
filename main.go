package main

import (
	"fmt"
	"net/http"

	"./lib/prometheus"
	"./lib/tarantool"
	"./services/fetch"
	"./services/send"
	"./utils/startup"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

var address string

func init() {
	startup.InitService()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err)
	}
	address = fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
}

func main() {
	connection := tarantool.GetTarantool()
	fetchQueue := tarantool.GetFetchQueue()
	sendQueue := tarantool.GetSendQueue()

	go fetchqueue.Listen(fetchQueue, connection)
	go sendqueue.Listen(sendQueue, connection)

	prometheus.MustRegister(metrics.CountSends, metrics.CountRelease, metrics.CountDeletes, metrics.CountFetchs)

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(address, nil)
}

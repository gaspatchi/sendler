package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var CountSends = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "sendler_send_mails", Help: "Количество отправленных писем"}, []string{"template"})
var CountFetchs = prometheus.NewCounter(prometheus.CounterOpts{Name: "sender_count_fetch", Help: "Количество полученных расписаний"})
var CountRelease = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "sendler_tasks_release", Help: "Количество возвращённых задач"}, []string{"queue"})
var CountDeletes = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "sendler_tasks_deletes", Help: "Количество удалённых задач"}, []string{"queue"})

func init() {
	CountSends.WithLabelValues("").Add(0)
	CountFetchs.Add(0)
	CountRelease.WithLabelValues("send").Add(0)
	CountDeletes.WithLabelValues("send").Add(0)
}

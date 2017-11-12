package fetchqueue

import (
	"os"

	"time"

	"../../lib/prometheus"

	"github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/queue"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func Listen(fetchQueue queue.Queue, connection *tarantool.Connection) {
	for {
		var taskdata FetchTask
		task, err := fetchQueue.TakeTimeout(time.Second * 5)
		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("При получении задачи произошла ошибка  ", err)
			continue
		}
		if task == nil {
			continue
		}
		err = taskdata.Unmarshal([]byte(task.Data().(string)), &taskdata)
		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("При разборе JSON произошла ошибка  ", err)
			continue
		}
		if taskdata.Data.Try == 3 {
			err = task.Delete()
			if err != nil {
				logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Ошибка при удалении задачи  ", err)
				continue
			}
			metrics.CountDeletes.WithLabelValues("Fetch").Inc()
			continue
		}
		err = fetchSchedule(connection, taskdata)
		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Ошибка при получении расписания  ", err)
			err = releaseTask(fetchQueue, task, taskdata)
			if err != nil {
				logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Невозможно вернуть задачу в очередь  ", err)
				continue
			}
			continue
		}
		err = task.Ack()
		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
			continue
		}
		metrics.CountFetchs.Inc()
		continue
	}
}

func releaseTask(fetchQueue queue.Queue, task *queue.Task, taskdata FetchTask) error {
	taskdata.Data.Try++
	err := task.Delete()
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Ошибка при удалении задачи  ", err)
		return err
	}
	newtask, err := taskdata.Marshal()
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Ошибка при создании JSON  ", err)
		return err
	}
	_, err = fetchQueue.PutWithOpts(newtask, queue.Opts{Delay: 20 * time.Second})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Ошибка при добавлении задачи в очередь  ", err)
		return err
	}
	metrics.CountRelease.WithLabelValues("Fetch").Inc()
	return nil
}

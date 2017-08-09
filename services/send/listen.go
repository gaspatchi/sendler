package sendqueue

import (
	"os"

	"time"

	"../../lib/prometheus"
	"github.com/sirupsen/logrus"
	tarantool "github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/queue"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func Listen(sendQueue queue.Queue, connection *tarantool.Connection) {
	for {
		var taskdata SendTask
		task, err := sendQueue.Take()

		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при получении задачи произошла ошибка  ", err)
			continue
		}
		err = taskdata.Unmarshal([]byte(task.Data().(string)), &taskdata)
		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при разборе JSON произошла ошибка  ", err)
			continue
		}
		if taskdata.Data.Try == 3 {
			err = task.Delete()
			if err != nil {
				logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при удалении задачи  ", err)
				continue
			}
			metrics.CountDeletes.WithLabelValues("Send").Inc()
			continue
		}
		if taskdata.Direction == "email" {
			switch taskdata.Template {
			case "verifyRegistration":
				err = VerifyRegistration(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма для верификации аккаунта  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("verifyRegistration").Inc()
				continue
			case "resetPassword":
				err = ResetPassword(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма для сброса пароля  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("resetPassword").Inc()
				continue
			case "createSchedule":
				err = SendNewSchedule(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма для верификации аккаунта  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("createSchedule").Inc()
				continue
			case "updateSchedule":
				err = SendUpdateSchedule(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма с обновлённым расписанием  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("updateSchedule").Inc()
				continue
			case "setNumber":
				err = SetNumber(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма с просьбой указать номер  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("setNumber").Inc()
				continue
			case "addMoney":
				err = AddMoney(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма с просьбой пополнить счёт  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("addMoney").Inc()
				continue
			case "feedback", "driving", "hairdresser", "florist":
				err = SendClaim(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма с заявкой  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("Feedback").Inc()
				continue
			case "sendMessage":
				err = SendMessage(connection, taskdata)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при отправке письма с текстом  ", err)
					err = releaseTask(sendQueue, task, taskdata)
					if err != nil {
						logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Невозможно вернуть задачу в очередь  ", err)
						continue
					}
					continue
				}
				err = task.Ack()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("При подтверждении выполнении задачи произошла ошибка  ", err)
					continue
				}
				metrics.CountSends.WithLabelValues("sendMessage").Inc()
				continue
			default:
				err = task.Release()
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при возвращении задачи в очередь  ", err)
					continue
				}
				continue
			}
		} else {
			err = task.Release()
			if err != nil {
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при возвращении задачи в очередь  ", err)
					continue
				}
			}
			continue
		}
	}
}

func releaseTask(sendQueue queue.Queue, task *queue.Task, taskdata SendTask) error {
	taskdata.Data.Try++
	err := task.Delete()
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при удалении задачи  ", err)
		return err
	}
	newtask, err := taskdata.Marshal()
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при создании JSON  ", err)
		return err
	}
	_, err = sendQueue.PutWithOpts(newtask, queue.Opts{Delay: 20 * time.Second})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Send"}).Error("Ошибка при добавлении задачи в очередь  ", err)
		return err
	}
	metrics.CountRelease.WithLabelValues("Send").Inc()
	return nil
}

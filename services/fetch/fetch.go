package fetchqueue

import (
	"os"

	"../../lib/inforer"
	"../../lib/schedler"
	"github.com/sirupsen/logrus"
	tarantool "github.com/tarantool/go-tarantool"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func fetchSchedule(connection *tarantool.Connection, task FetchTask) error {
	if task.Type == "teacher" {
		err := fetchTeacher(connection, task)
		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Невозможно получить расписание преподавателя  ",err)
			return err
		}
	} else if task.Type == "group" {
		err := fetchGroup(connection, task)
		if err != nil {
			logrus.WithFields(logrus.Fields{"module": "Fetch"}).Error("Невозможно получить расписание группы  ",err)
			return err
		}
	}
	return nil
}

func fetchTeacher(connection *tarantool.Connection, task FetchTask) error {
	var schedule TeacherSchedule
	var teacherInfo inforer.TeacherInfo
	var teacherSchedule schedler.TeacherSchedule

	teacherInfo, err := inforer.GetTeacher(task.ID)
	if err != nil {
		return err
	}
	teacherSchedule, err = schedler.GetTeacher(task.ID, task.Date)
	if err != nil {
		return err
	}

	schedule.Info = teacherInfo
	schedule.Schedule = teacherSchedule

	err = schedule.Post(connection, task)
	if err != nil {
		return err
	}
	return nil
}

func fetchGroup(connection *tarantool.Connection, task FetchTask) error {
	var schedule GroupSchedule
	var groupInfo inforer.GroupInfo
	var groupSchedule schedler.GroupSchedule

	groupInfo, err := inforer.GetGroup(task.ID)
	if err != nil {
		return err
	}
	groupSchedule, err = schedler.GetGroup(task.ID, task.Date)
	if err != nil {
		return err
	}

	schedule.Info = groupInfo
	schedule.Schedule = groupSchedule

	err = schedule.Post(connection, task)
	if err != nil {
		return err

	}
	return nil
}

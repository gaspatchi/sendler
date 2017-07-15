package schedler

type TeacherSchedule struct {
	Schedule []struct {
		Date  string `json:"date"`
		Index int    `json:"index"`
		Group struct {
			GroupID int         `json:"group_id"`
			Group   string      `json:"group"`
			Course  interface{} `json:"course"`
		} `json:"group"`
		Lesson struct {
			LessonID int    `json:"lesson_id"`
			Lesson   string `json:"lesson"`
		} `json:"lesson"`
		Cabinet struct {
			CabinetID int    `json:"cabinet_id"`
			Cabinet   string `json:"cabinet"`
		} `json:"cabinet"`
	} `json:"schedule"`
}

type GroupSchedule struct {
	Schedule []struct {
		Date  string `json:"date"`
		Index int    `json:"index"`
		Group struct {
			GroupID int         `json:"group_id"`
			Group   string      `json:"group"`
			Course  interface{} `json:"course"`
		} `json:"group"`
		Lesson struct {
			LessonID int    `json:"lesson_id"`
			Lesson   string `json:"lesson"`
		} `json:"lesson"`
		Teacher struct {
			TeacherID  int    `json:"teacher_id"`
			Firstname  string `json:"firstname"`
			Lastname   string `json:"lastname"`
			Patronymic string `json:"patronymic"`
		} `json:"teacher"`
		Cabinet struct {
			CabinetID int    `json:"cabinet_id"`
			Cabinet   string `json:"cabinet"`
		} `json:"cabinet"`
	} `json:"schedule"`
}

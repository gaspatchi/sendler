package inforer

type TeacherInfo struct {
	Result struct {
		Teacher struct {
			TeacherID  int    `json:"teacher_id"`
			Firstname  string `json:"firstname"`
			Lastname   string `json:"lastname"`
			Patronymic string `json:"patronymic"`
		} `json:"teacher"`
	} `json:"result"`
}

type GroupInfo struct {
	Result struct {
		Group struct {
			GroupID int    `json:"group_id"`
			Group   string `json:"group"`
			Course  int    `json:"course"`
		} `json:"group"`
	} `json:"result"`
}

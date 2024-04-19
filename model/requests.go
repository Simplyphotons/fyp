package model

type AuthorizationRequest struct { //400
	Code         string `json:"code"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type GetQuestionsResponse struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	ID            string `json:"id"`
	StudentID     string `json:"studentID,omitempty"`
	SupervisorID  string `json:"supervisorID,omitempty"`
	QuestionShort string `json:"questionShort,omitempty"`
	Question      string `json:"question,omitempty"`
	Answer        string `json:"answer"`
	IsAnswered    bool   `json:"isAnswered"`

	//Question string `json:"question,omitempty"`
}

type GetGanttResponse struct {
	Gantt []Gantt `json:"gantt"`
}

type Gantt struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectID"`
	GanttName   string `json:"ganttName"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Description string `json:"description"`
	Links       string `json:"links"`
	Feedback    string `json:"feedback"`
	newFeedback string `json:"newFeedback"`
	Colour      string `json:"colour"`
}

type GanttChartRow struct {
	Content []Gantt `json:"content"`
}

type AccountSupervisor struct { //for checking account type
	ID           string `json:"id"`
	IsSupervisor bool   `json:"isSupervisor"`
}

type QuestionData struct {
	ID            string `json:"id"`
	StudentID     string `json:"studentID"`
	SupervisorID  string `json:"supervisorID"`
	QuestionShort string `json:"questionShort"`
	QuestionLong  string `json:"questionLong"`
}

type ApplicationData struct {
	ID           string `json:"id,omitempty"`
	StudentID    string `json:"student_id"`
	StudentName  string `json:"student_name"`
	SupervisorID string `json:"supervisor_id"`
	Heading      string `json:"heading"`
	Description  string `json:"description"`
	Accepted     bool   `json:"accepted"`
	Declined     bool   `json:"declined"`
}

type ProjectData struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	StudentID    string `json:"studentID"`
	SupervisorID string `json:"supervisorID"`
}

type UserData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	isSupervisor bool   `json:"isSupervisor"`
}

type UserCreateRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Verify struct {
	UserId string `json:"userId"`
	Found  bool   `json:"found"`
}

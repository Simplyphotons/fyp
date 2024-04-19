package db

type User struct {
	Id   string //id
	Name string
}

type Project struct {
	ID           string
	Name         string
	StudentID    string
	SupervisorID string
}

type Application struct {
	ID           string
	StudentID    string
	StudentName  string
	SupervisorID string
	Heading      string
	Description  string
	Accepted     bool
	Declined     bool
}

type Gantt struct {
	Id          string
	ProjectID   string
	GanttName   string
	StartDate   string
	EndDate     string
	Description string
	Links       string
	Feedback    string
	newFeedBack string
}

type Question struct {
	Id            string
	studentID     string
	supervisorID  string
	questionShort string
	questionLong  string
	answer        string
	is_answered   bool
}

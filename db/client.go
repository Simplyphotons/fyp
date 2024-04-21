package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Simplyphotons/fyp.git/model"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func (db Client) GetQuestions(ctx context.Context) ([]model.Question, error) {
	rows, err := db.conn.QueryContext(ctx, "SELECT ticket_id, questionshort from tickets")
	if err != nil {
		log.Printf("cannot execute query to get questions: %v", err)
		return nil, err
	}

	result := []model.Question{}

	var (
		id       string
		question string
	)
	for rows.Next() {
		err = rows.Scan(&id, &question)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
			return nil, err
		}

		result = append(result, model.Question{
			ID:       id,
			Question: question,
		})
	}

	// Nullable fields
	//var (
	//	id       string
	//	question sql.NullString
	//)
	//for rows.Next() {
	//	err = rows.Scan(&id, &question)
	//	if err != nil {
	//		log.Printf("cannot read data while getting questions: %v", err)
	//	}
	//
	//	item := model.Question{
	//		ID: id,
	//	}
	//
	//	if question.Valid {
	//		item.Question = &question.String
	//	}
	//	temp = append(temp)
	//}

	// Not nullable fields
	//for rows.Next() {
	//	item := model.Question{}
	//
	//	err := rows.Scan(&item.ID, &item.Question)
	//	if err != nil {
	//		log.Printf("cannot read data while getting questions: %v", err)
	//	}
	//	temp = append(temp, item)
	//}
	return result, nil

}

func (db Client) GetGantt(ctx context.Context, projectIdentifier string) ([]model.GanttChartRow, error) { //gets all milestones within a project
	rows, err := db.conn.QueryContext(ctx, "SELECT item_id, project_id, gantt_name, start_date, end_date, description, links, feedback, colour from gantt_items where project_id = $1 order by start_date", projectIdentifier)
	if err != nil {
		log.Printf("cannot execute query to get questions: %v", err)
		return nil, err
	}

	result := []model.GanttChartRow{}
	var (
		id          string
		projectID   string
		startDate   string
		ganttName   string
		endDate     string
		description string
		links       string
		feedback    string
		colour      string
	)
	for rows.Next() {
		err = rows.Scan(&id, &projectID, &ganttName, &startDate, &endDate, &description, &links, &feedback, &colour)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
			return nil, err
		}

		row := []model.Gantt{
			{
				ID:          id,
				ProjectID:   projectID,
				GanttName:   ganttName,
				StartDate:   startDate,
				EndDate:     endDate,
				Description: description,
				Links:       links,
				Feedback:    feedback,
				Colour:      colour,
			},
		}

		result = append(result, model.GanttChartRow{
			Content: row,
		})
	}
	println(id, projectID, ganttName, startDate, endDate)
	return result, nil
}

func (db Client) GetGanttItem(ctx context.Context, milestoneIdentifier string) ([]model.Gantt, error) { //gets one milestone
	rows, err := db.conn.QueryContext(ctx, "SELECT item_id, project_id, gantt_name, start_date, end_date, description, links, feedback from gantt_items where item_id = $1", milestoneIdentifier)
	if err != nil {
		log.Printf("cannot execute query to get questions: %v", err)
		return nil, err
	}

	result := []model.Gantt{}
	var (
		id          string
		projectID   string
		ganttName   string
		startDate   string
		endDate     string
		description string
		links       string
		feedback    string
	)
	for rows.Next() {
		err = rows.Scan(&id, &projectID, &ganttName, &startDate, &endDate, &description, &links, &feedback)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
			return nil, err
		}

		result = append(result, model.Gantt{
			ID:          id,
			ProjectID:   projectID,
			GanttName:   ganttName,
			StartDate:   startDate,
			EndDate:     endDate,
			Description: description,
			Links:       links,
			Feedback:    feedback,
		})

	}
	return result, nil
}

func (db Client) GetSupervisors(ctx context.Context) ([]model.UserData, error) { //for use in displaying all available supervisors when a student is creating a new project application.
	rows, err := db.conn.QueryContext(ctx, "SELECT id, name from users where is_supervisor = true")
	if err != nil {
		log.Printf("cannot execute query to get users: %v", err)
		return nil, err
	}

	result := []model.UserData{}

	var (
		id   string
		name string
	)
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Printf("cannot read data while getting users: %v", err)
		}

		result = append(result, model.UserData{
			ID:   id,
			Name: name,
		})
	}
	return result, nil
}

func (db Client) GetHasProjectStatus(ctx context.Context, userID string) (bool, error) {
	rows, err := db.conn.QueryContext(ctx, "SELECT has_project from users where id = $1", userID)
	if err != nil {
		log.Printf("cannot execute query to get users: %v", err)
		return false, err
	}

	var status bool

	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			log.Printf("cannot read data while getting users: %v", err)
		}
	}
	if status {
		return true, nil
	}
	return false, nil
}

func (db Client) GetApplications(ctx context.Context, supervisor_ID string) ([]model.ApplicationData, error) {
	query := `SELECT a.id, a.student_id, u.name, a.supervisor_id, a.heading, a.description, a.accepted, a.declined
FROM applications a INNER JOIN users u
    ON a.student_id = u.id
WHERE supervisor_id = $1`
	rows, err := db.conn.QueryContext(ctx, query, supervisor_ID)
	if err != nil {
		log.Printf("cannot execute query to get applications: %v", err)
		return nil, err
	}

	result := []model.ApplicationData{}
	var (
		id           string
		studentID    string
		studentName  string
		supervisorID string
		heading      string
		description  string
		accepted     bool
		declined     bool
	)
	for rows.Next() {
		err = rows.Scan(&id, &studentID, &studentName, &supervisorID, &heading, &description, &accepted, &declined)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
		}
		if !accepted {
			result = append(result, model.ApplicationData{
				ID:           id,
				StudentID:    studentID,
				StudentName:  studentName,
				SupervisorID: supervisorID,
				Heading:      heading,
				Description:  description,
				Accepted:     accepted,
				Declined:     declined,
			})
		}
	}
	return result, nil

}

func (db Client) GetAllAcceptedRequests(ctx context.Context) ([]model.ApplicationData, error) {
	query := `SELECT a.id, a.student_id, u.name, a.supervisor_id, a.heading, a.description, a.accepted, a.declined
FROM applications a INNER JOIN users u
    ON a.student_id = u.id
WHERE a.accepted = true`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		log.Printf("cannot execute query to get applications: %v", err)
		return nil, err
	}

	result := []model.ApplicationData{}
	var (
		id           string
		studentID    string
		studentName  string
		supervisorID string
		heading      string
		description  string
		accepted     bool
		declined     bool
	)
	for rows.Next() {
		err = rows.Scan(&id, &studentID, &studentName, &supervisorID, &heading, &description, &accepted, &declined)
		if err != nil {
			log.Printf("cannot read data while getting reading prompts: %v", err)
		}
		if accepted {
			result = append(result, model.ApplicationData{
				ID:           id,
				StudentID:    studentID,
				StudentName:  studentName,
				SupervisorID: supervisorID,
				Heading:      heading,
				Description:  description,
				Accepted:     accepted,
				Declined:     declined,
			})
		}
	}
	return result, nil

}

func (db Client) GetApplicationsForStudent(ctx context.Context, student_ID string) ([]model.ApplicationData, error) {
	query := `SELECT a.id, a.student_id, u.name, a.supervisor_id, a.heading, a.description, a.accepted, a.declined
FROM applications a INNER JOIN users u
    ON a.supervisor_id = u.id
WHERE student_id = $1`
	rows, err := db.conn.QueryContext(ctx, query, student_ID)
	if err != nil {
		log.Printf("cannot execute query to get applications: %v", err)
		return nil, err
	}

	result := []model.ApplicationData{}
	var (
		id           string
		studentID    string
		studentName  string
		supervisorID string
		heading      string
		description  string
		accepted     bool
		declined     bool
	)
	for rows.Next() {
		err = rows.Scan(&id, &studentID, &studentName, &supervisorID, &heading, &description, &accepted, &declined)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
		}

		result = append(result, model.ApplicationData{
			ID:           id,
			StudentID:    studentID,
			StudentName:  studentName,
			SupervisorID: supervisorID,
			Heading:      heading,
			Description:  description,
			Accepted:     accepted,
			Declined:     declined,
		})
	}
	return result, nil

}
func (db Client) GetProjects(ctx context.Context, supervisor_id string) ([]model.ProjectData, error) {
	rows, err := db.conn.QueryContext(ctx, "SELECT project_id, project_name, student_id, supervisor_id from projects where supervisor_id = $1", supervisor_id)
	if err != nil {
		log.Printf("cannot execute query to get projects: %v", err)
		return nil, err
	}
	result := []model.ProjectData{}
	var (
		projectID    string
		projectName  string
		studentID    string
		supervisorID string
	)
	for rows.Next() {
		err = rows.Scan(&projectID, &projectName, &studentID, &supervisorID)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
		}

		result = append(result, model.ProjectData{
			ID:           projectID,
			Name:         projectName,
			StudentID:    studentID,
			SupervisorID: supervisorID,
		})
	}
	return result, nil

}

func (db Client) GetSecondProjects(ctx context.Context, supervisor_id string) ([]model.ProjectData, error) {
	rows, err := db.conn.QueryContext(ctx, "SELECT project_id, project_name, student_id, supervisor_id from projects where second_reader_id = $1", supervisor_id)
	if err != nil {
		log.Printf("cannot execute query to get projects: %v", err)
		return nil, err
	}
	result := []model.ProjectData{}
	var (
		projectID    string
		projectName  string
		studentID    string
		supervisorID string
	)
	for rows.Next() {
		err = rows.Scan(&projectID, &projectName, &studentID, &supervisorID)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
		}

		result = append(result, model.ProjectData{
			ID:           projectID,
			Name:         projectName,
			StudentID:    studentID,
			SupervisorID: supervisorID,
		})
	}
	return result, nil

}

func (db Client) GetProjectID(ctx context.Context, userID string) (*model.ProjectData, error) {
	rows, err := db.conn.QueryContext(ctx, "select project_id from projects where student_id = $1", userID)
	if err != nil {
		log.Printf("cannot execute query to get project id: %v", err)
		return nil, err
	}

	result := &model.ProjectData{}
	for rows.Next() {
		err = rows.Scan(&result.ID)
		if err != nil {
			log.Printf("cannot read project ID: %v", err)
			return nil, err
		}
	}
	println(result.ID)
	return result, err
}
func (db Client) GetProjectName(ctx context.Context, ProjectID string) (*model.ProjectData, error) {
	rows, err := db.conn.QueryContext(ctx, "select project_name from projects where project_id = $1", ProjectID)
	if err != nil {
		log.Printf("cannot execute query to get project name: %v", err)
		return nil, err
	}

	result := &model.ProjectData{}
	for rows.Next() {
		err = rows.Scan(&result.Name)
		if err != nil {
			log.Printf("cannot read project Name: %v", err)
			return nil, err
		}
	}
	println(result.ID)
	return result, err
}

func (db Client) GetSecondReaderStatus(ctx context.Context, ProjectID string, userID string) (bool, error) {
	rows, err := db.conn.QueryContext(ctx, "select second_reader_id from projects where project_id = $1", ProjectID)
	if err != nil {
		log.Printf("cannot execute query to get project name: %v", err)
		return false, err
	}

	var result string
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			log.Printf("cannot read project Name: %v", err)
			return false, err
		}
	}
	println(result)

	if result == userID {
		return true, nil
	}

	return false, nil
}

func (db Client) GetSpecificApplications(ctx context.Context, appID string) ([]model.ApplicationData, error) {
	rows, err := db.conn.QueryContext(ctx, "SELECT id, student_id, supervisor_id, heading, description, accepted, declined from applications where id = $1", appID)
	if err != nil {
		log.Printf("cannot execute query to get applications: %v", err)
		return nil, err
	}
	result := []model.ApplicationData{}
	var (
		id           string
		studentID    string
		supervisorID string
		heading      string
		description  string
		accepted     bool
		declined     bool
	)
	for rows.Next() {
		err = rows.Scan(&id, &studentID, &supervisorID, &heading, &description, &accepted, &declined)
		if err != nil {
			log.Printf("cannot read data while getting questions: %v", err)
		}

		result = append(result, model.ApplicationData{
			ID:           id,
			StudentID:    studentID,
			SupervisorID: supervisorID,
			Heading:      heading,
			Description:  description,
			Accepted:     accepted,
			Declined:     declined,
		})

	}
	return result, nil

}

func (db Client) NewQuestion(ctx context.Context, question Question) error { //adds new question to db
	//create question item for db
	//variables that are created and not taken: ID, Isanswered will be false

	id := GenerateUUID()
	studentID := question.studentID
	supervisorID := question.supervisorID
	questionShort := question.questionShort
	questionLong := question.questionLong
	answer := ""
	isAnswered := false

	updateQuery := "INSERT INTO tickets (ticket_id, student_id, supervispr_id, questionshort, questionlong, answer, is_answered) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	result, err := db.conn.Exec(updateQuery, id, studentID, supervisorID, questionShort, questionLong, answer, isAnswered)

	if err != nil {
		log.Printf("failed to add answer to corresponding ticket in db")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)
	return nil
}

func (db Client) NewAnswer(ctx context.Context) error { //adds new question to db
	//todo
	var questionData model.Question

	updateQuery := "UPDATE tickets SET answer = ?, is_answered = ? WHERE ticket_id = ?"

	id := questionData.ID
	answer := questionData.Answer
	isAnswered := true

	result, err := db.conn.Exec(updateQuery, answer, isAnswered, id)
	if err != nil {
		log.Printf("failed to add answer to corresponding ticket in db")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Updated %d rows.\n", rowsAffected)
	return nil
}

func (db Client) CreateProject(ctx context.Context, application Application, supervisor_id string) error {
	id := GenerateUUID()
	name := application.StudentName
	studentID := application.StudentID
	supervisorID := supervisor_id //takes the auth userid since only the supervisor can access this

	updateQuery := "INSERT INTO projects (project_id, project_name, student_id, supervisor_id) VALUES ($1, $2, $3, $4)"

	result, err := db.conn.Exec(updateQuery, id, name, studentID, supervisorID)
	if err != nil {
		log.Printf("failed to add answer to corresponding ticket in db")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)
	db.updateUser(studentID)
	db.acceptApplication(application.ID)
	db.deleteApplication(studentID)
	return nil

}

func (db Client) CreateSupervisorUser(ctx context.Context, user User) error {
	id := user.Id
	name := user.Name

	updateQuery := "INSERT INTO users (id, name, is_supervisor, has_project) VALUES ($1, $2, $3, $4)"

	result, err := db.conn.Exec(updateQuery, id, name, true, false)
	if err != nil {
		log.Printf("failed to add user to corresponding ticket in db")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)
	return nil

}

func (db Client) CreateStudentUser(ctx context.Context, user User) error {
	id := user.Id
	name := user.Name

	updateQuery := "INSERT INTO users (id, name, is_supervisor, has_project) VALUES ($1, $2, $3, $4)"

	result, err := db.conn.Exec(updateQuery, id, name, false, false)
	if err != nil {
		log.Printf("failed to add user to corresponding ticket in db")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)
	return nil
}

func (db Client) CreateApplication(ctx context.Context, application Application, student_id string) error {
	applicationID := GenerateUUID()
	studentID := student_id
	supervisorID := application.SupervisorID
	heading := application.Heading
	description := application.Description
	accepted := false
	declined := false

	updateQuery := "INSERT INTO applications (id, student_id, supervisor_id, heading, description, accepted, declined) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	//println(applicationID + "\n\n")
	//println(applicationID + "\n" + student_id + "\n" + supervisorID + "\n" + heading + "\n" + description)

	result, err := db.conn.Exec(updateQuery, applicationID, studentID, supervisorID, heading, description, accepted, declined)
	if err != nil {
		log.Printf("failed to add new appliction")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)
	return nil
}

func (db Client) acceptApplication(id string) error {
	println(id)
	println("call succ")
	updateQuery := "UPDATE applications SET accepted = $1 WHERE id = $2"

	result, err := db.conn.Exec(updateQuery, true, id)
	if err != nil {
		log.Printf("failed to accept application")
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("updated %d row.\n", rowsAffected)

	return nil
}

func (db Client) CompleteGanttItem(ctx context.Context, gantt Gantt) error {

	updateQuery := "UPDATE gantt_items SET colour = $1 WHERE item_id = $2"

	newColour := "#2C59C7"
	result, err := db.conn.Exec(updateQuery, newColour, gantt.Id)
	if err != nil {
		log.Printf("failed to accept application")
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)
	return nil
}

func (db Client) DeclineApplication(ctx context.Context, application Application) error {
	db.deleteSpecificApplication(application.ID)
	return nil
}

func (db Client) deleteApplication(condition string) error {
	query := "DELETE FROM applications WHERE student_id = $1 AND accepted = $2"

	result, err := db.conn.Exec(query, condition, false)
	if err != nil {
		log.Printf("failed to delete table")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Row %d deleted.\n", condition)
	log.Printf("%d Rows affected.\n", rowsAffected)
	return nil
}

func (db Client) deleteSpecificApplication(condition string) error {
	query := "DELETE FROM applications WHERE id = $1 AND accepted = $2"

	result, err := db.conn.Exec(query, condition, false)
	if err != nil {
		log.Printf("failed to delete table")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Row %d deleted.\n", condition)
	log.Printf("%d Rows affected.\n", rowsAffected)
	return nil
}

func (db Client) DeleteGanttItem(id string) error {
	query := "DELETE FROM gantt_items WHERE item_id = $1"

	result, err := db.conn.Exec(query, id)
	if err != nil {
		log.Printf("failed to delete gantt item")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Row %d deleted.\n", id)
	log.Printf("%d Rows affected.\n", rowsAffected)
	return nil
}

func (db Client) CreateGanttItem(ctx context.Context, gantt Gantt) error {

	id := GenerateUUID()
	projectID := gantt.ProjectID
	description := gantt.Description
	startDate := gantt.StartDate
	endDate := gantt.EndDate
	feedback := ""
	links := gantt.Links
	ganttName := gantt.GanttName
	colour := "#2A9D39"
	tracker := 0

	updateQuery := "INSERT INTO gantt_Items (item_id, project_id, description, start_date, end_date, feedback, links, gantt_name, colour, feedback_update_tracker) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	result, err := db.conn.Exec(updateQuery, id, projectID, description, startDate, endDate, feedback, links, ganttName, colour, tracker)
	if err != nil {
		log.Printf("failed to create new gantt item/milestone to the project")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)
	return nil
}

func (db Client) UpdateFeedback(ctx context.Context, gantt Gantt, userID string) error {
	newText := ""
	updateQuery := "UPDATE gantt_items SET feedback = $1 WHERE item_id = $2"
	isSupervisor, err := db.getAccountStatus(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve account status")
		return err
	}
	if isSupervisor {
		newText = gantt.Feedback + "Supervisor: " + gantt.NewFeedBack + "\n\n"
		db.enableAlert(1, gantt.Id)
	} else {
		newText = gantt.Feedback + "Student: " + gantt.NewFeedBack + "\n\n"
		db.enableAlert(2, gantt.Id)
	}

	result, err := db.conn.Exec(updateQuery, newText, gantt.Id)
	if err != nil {
		log.Printf("failed to update feedback")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\nFeedback updated", rowsAffected)
	return nil

}

func (db Client) enableAlert(number int, itemID string) error {
	alert, err := db.conn.Exec("UPDATE gantt_items SET feedback_update_tracker = $1, colour = '#e6e600' WHERE item_id = $2", number, itemID)
	if err != nil {
		log.Printf("failed to update feedback status and colour")
		return err
	}
	rowsAffected, _ := alert.RowsAffected()
	log.Printf("created %d row.\nFeedback status and alert colour updated", rowsAffected)
	return nil
}

func (db Client) DisableAlert(ctx context.Context, userID string, ganttID string) error {
	isSupervisor, err := db.getAccountStatus(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve account status")
		return err
	}
	if isSupervisor {
		row := db.conn.QueryRowContext(ctx, "SELECT feedback_update_tracker FROM gantt_items WHERE item_id = $1", ganttID)
		var val int
		err := row.Scan(&val)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user with ID %s not found", ganttID)
			}
			// Handle other errors
			log.Printf("cannot read data while getting alert: %v", err)
			return err
		}
		if val == 2 {
			alert, err := db.conn.Exec("UPDATE gantt_items SET feedback_update_tracker = 0, colour = '#2A9D39' WHERE item_id = $1", ganttID)
			if err != nil {
				log.Printf("failed to update feedback status and colour")
				return err
			}
			rowsAffected, _ := alert.RowsAffected()
			log.Printf("created %d row.\nFeedback status and alert colour updated", rowsAffected)
			return nil
		} else {
			log.Printf("failed to remove alert as account is type supervisor when should be type student")
			return nil
		}
	} else {

		row := db.conn.QueryRowContext(ctx, "SELECT feedback_update_tracker FROM gantt_items WHERE item_id = $1", ganttID)
		var val int
		err := row.Scan(&val)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user with ID %s not found", ganttID)
			}
			// Handle other errors
			log.Printf("cannot read data while getting alert: %v", err)
			return err
		}

		if val == 1 {
			alert, err := db.conn.Exec("UPDATE gantt_items SET feedback_update_tracker = 0, colour = '#2A9D39' WHERE item_id = $1", ganttID)
			if err != nil {
				log.Printf("failed to update feedback status and colour")
				return err
			}
			rowsAffected, _ := alert.RowsAffected()
			log.Printf("created %d row.\nFeedback status and alert colour updated", rowsAffected)
			return nil
		} else {
			log.Printf("failed to remove alert as account is type student when should be type supervisor")
			return nil
		}
	}

}

func (db Client) getAccountStatus(ctx context.Context, id string) (bool, error) {

	row := db.conn.QueryRowContext(ctx, "SELECT is_supervisor FROM users WHERE id = $1", id)

	var is_supervisor bool
	err := row.Scan(&is_supervisor)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("user with ID %s not found", id)
		}
		log.Printf("cannot read data while getting account status: %v", err)
		return false, err
	}
	return is_supervisor, nil
}

func (db Client) GetUsername(ctx context.Context, userId string) (string, error) {
	row := db.conn.QueryRowContext(ctx, "SELECT name FROM users WHERE id = $1", userId)
	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user with ID %s not found", userId)
		}
		// Handle other errors
		log.Printf("cannot read data while getting username: %v", err)
		return "", err
	}
	return name, nil
}

func (db Client) GetFeedback(ctx context.Context, ganttID string) (string, error) {
	row := db.conn.QueryRowContext(ctx, "SELECT feedback FROM gantt_items WHERE item_id = $1", ganttID)
	var feedback string
	err := row.Scan(&feedback)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("feedback for gantt ID %s not found", ganttID)
		}
		// Handle other errors
		log.Printf("cannot read data while getting feedback: %v", err)
		return "", err
	}
	return feedback, nil
}

func (db Client) Verify(ctx context.Context, userID string) (*model.Verify, error) {
	result := &model.Verify{}
	rows, err := db.conn.QueryContext(ctx, "SELECT id, name from users where id = $1", userID)
	if err != nil {
		log.Printf("cannot get user.  %v", err)
		return nil, err
	}

	result.UserId = userID

	if !rows.Next() {
		result.Found = false
		return result, nil
	}

	result.Found = true
	return result, nil

}

func (db Client) updateUser(student_id string) error {

	updateQuery := "UPDATE users SET has_project = $1 WHERE id = $2"

	result, err := db.conn.Exec(updateQuery, true, student_id)
	if err != nil {
		log.Printf("failed to accept application")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)

	return nil

}

func (db Client) AddSecondReader(ctx context.Context, readerID string, appID string) error {

	query := `SELECT p.project_id
FROM applications a INNER JOIN projects p
    ON a.student_id = p.student_id
WHERE id = $1`

	var projectID string
	rows, err := db.conn.QueryContext(ctx, query, appID)
	if err != nil {
		log.Printf("failed to read application")
		return err
	}
	for rows.Next() {
		err = rows.Scan(&projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("feedback for gantt ID %s not found", projectID)
			}
			// Handle other errors
			log.Printf("cannot read data while getting feedback: %v", err)
			return err
		}
		println("got project id: ", projectID)
	}
	updateQuery := "UPDATE projects SET second_reader_id = $1 WHERE project_id = $2"

	result, err := db.conn.Exec(updateQuery, readerID, projectID)
	if err != nil {
		log.Printf("failed to accept application")
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("created %d row.\n", rowsAffected)

	return nil

}

func GenerateUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}

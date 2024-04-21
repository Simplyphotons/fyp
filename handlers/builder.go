package handlers

import (
	"context"
	"github.com/Simplyphotons/fyp.git/auth0"
	"github.com/Simplyphotons/fyp.git/db"
	"github.com/Simplyphotons/fyp.git/model"
)

type DBClient interface {
	GetGanttItem(ctx context.Context, milestoneIdentifier string) ([]model.Gantt, error)
	CreateProject(ctx context.Context, project db.Application, supervisor_id string) error
	CreateApplication(ctx context.Context, application db.Application, student_id string) error
	DeclineApplication(ctx context.Context, application db.Application) error
	GetQuestions(ctx context.Context) ([]model.Question, error)
	GetSupervisors(ctx context.Context) ([]model.UserData, error)
	GetHasProjectStatus(ctx context.Context, userID string) (bool, error)
	GetApplications(ctx context.Context, supervisor_id string) ([]model.ApplicationData, error)
	GetApplicationsForStudent(ctx context.Context, student_id string) ([]model.ApplicationData, error)
	GetSpecificApplications(ctx context.Context, appID string) ([]model.ApplicationData, error)
	GetProjects(ctx context.Context, supervisor_id string) ([]model.ProjectData, error)
	GetProjectID(ctx context.Context, userID string) (*model.ProjectData, error)
	GetProjectName(ctx context.Context, projectID string) (*model.ProjectData, error)
	GetFeedback(ctx context.Context, ganttID string) (string, error)
	GetUsername(ctx context.Context, userId string) (string, error)
	NewQuestion(ctx context.Context, question db.Question) error
	NewAnswer(ctx context.Context) error
	GetGantt(ctx context.Context, projectIdentifier string) ([]model.GanttChartRow, error)
	CreateGanttItem(ctx context.Context, gantt db.Gantt) error
	UpdateFeedback(ctx context.Context, gantt db.Gantt, userID string) error
	DisableAlert(ctx context.Context, userID string, ganttID string) error
	DeleteGanttItem(id string) error
	AddSecondReader(ctx context.Context, readerID string, appID string) error
	GetAllAcceptedRequests(ctx context.Context) ([]model.ApplicationData, error)
	CreateSupervisorUser(ctx context.Context, user db.User) error
	CreateStudentUser(ctx context.Context, user db.User) error
	GetSecondProjects(ctx context.Context, supervisor_id string) ([]model.ProjectData, error)
	GetSecondReaderStatus(ctx context.Context, ProjectID string, userID string) (bool, error)
	CompleteGanttItem(ctx context.Context, gantt db.Gantt) error
	Verify(ctx context.Context, userID string) (*model.Verify, error)
}

type Auth0Client interface {
	AddRole(ctx context.Context, userId string, roleId string) error
	DoesUserExist(ctx context.Context, email string) (bool, error)
	AddUser(ctx context.Context, r auth0.UserCreateRequest) (string, error)
}

type Controller struct {
	dbClient         DBClient
	auth0Client      Auth0Client
	supervisorRoleID string
}

func New(client DBClient, auth0Client Auth0Client, supervisorRoleID string) *Controller {
	return &Controller{
		dbClient:         client,
		auth0Client:      auth0Client,
		supervisorRoleID: supervisorRoleID,
	}
}

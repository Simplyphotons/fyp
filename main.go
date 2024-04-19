package main

import (
	"github.com/Simplyphotons/fyp.git/auth0"
	"github.com/Simplyphotons/fyp.git/db"
	"github.com/Simplyphotons/fyp.git/handlers"
	"github.com/Simplyphotons/fyp.git/oauth2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	dbClient := db.MustCreate(os.Getenv("DB_URL"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")) // create db client

	debug := false
	debugStr := os.Getenv("DEBUG")
	if debugStr != "" {
		if strings.ToUpper(debugStr) == "TRUE" {
			debug = true
		}
	}

	auth0Url := os.Getenv("AUTH0_BASE_URL")
	if auth0Url == "" {
		slog.Error("AUTH0_BASE_URL must be specified")
		os.Exit(1)
	}

	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	if auth0Audience == "" {
		slog.Error("AUTH0_AUDIENCE must be specified")
		os.Exit(1)
	}

	auth0ClientID := os.Getenv("AUTH0_CLIENT_ID")
	if auth0ClientID == "" {
		slog.Error("AUTH0_CLIENT_ID must be specified")
		os.Exit(1)
	}

	auth0ClientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	if auth0ClientSecret == "" {
		slog.Error("AUTH0_CLIENT_SECRET must be specified")
		os.Exit(1)
	}

	supervisorRoleID := os.Getenv("SUPERVISOR_ROLE_ID")
	if supervisorRoleID == "" {
		slog.Error("SUPERVISOR_ROLE_ID must be specified")
		os.Exit(1)
	}

	auth0Client, err := auth0.Build(
		auth0.Debug(debug),
		auth0.BaseUrl(auth0Url),
		auth0.Audience(auth0Audience),
		auth0.ClientId(auth0ClientID),
		auth0.ClientSecret(auth0ClientSecret),
		auth0.HTTPClient(&http.Client{}),
	)
	if err != nil {
		log.Printf("cannot create Auth0 client: %v", err)
		os.Exit(2)
	}

	controller := handlers.New(dbClient, auth0Client, supervisorRoleID) //dependency injection

	// Initialize oauth2 middleware
	oauth2Config, err := oauth2.Build(
		oauth2.Debug(debug),
		oauth2.URL(os.Getenv("JWKS_URL")),
		oauth2.Unmatched(true),
		oauth2.Audience(os.Getenv("AUDIENCE")),
		oauth2.Issuer(os.Getenv("ISSUER")),
		oauth2.HTTPClient(&http.Client{}),
		oauth2.Request("GET", "/questions", []string{"read:questions"}),
		oauth2.Request("POST", "/createApplication", []string{"read:student"}),
		oauth2.Request("POST", "/newQuestion", []string{"read:student"}),
		oauth2.Request("POST", "/newAnswer", []string{"read:supervisor"}),
		oauth2.Request("GET", "/getQuestions", []string{"read:supervisor", "read:student"}),
		oauth2.Request("GET", "/getApplications", []string{"read:supervisor", "read:student"}),
		oauth2.Request("GET", "/GetAllAcceptedRequests", []string{"read:supervisor", "read:student"}),
		oauth2.Request("GET", "/getApplicationsForStudent", []string{"read:student"}),
		oauth2.Request("GET", "/getSpecificApplications", []string{"read:supervisor", "read:student"}),
		oauth2.Request("GET", "/getGanttItem/:id", []string{"read:supervisor", "read:student"}),
		oauth2.Request("GET", "/getGantt/:id", []string{"read:supervisor", "read:student"}),
		oauth2.Request("GET", "/getSupervisors", []string{"read:student"}),
		oauth2.Request("GET", "/getProjects", []string{"read:student", "read:supervisor"}),
		oauth2.Request("GET", "/getUsername/:id", []string{"read:student", "read:supervisor"}),
		oauth2.Request("GET", "/getProjectID", []string{"read:student"}),
		oauth2.Request("GET", "/getProjectStatus", []string{"read:student"}),
		oauth2.Request("GET", "/getProjectName/:id", []string{"read:student", "read:supervisor"}),
		oauth2.Request("GET", "/getGantt/:id", []string{"read:student", "read:supervisor"}),
		oauth2.Request("GET", "/verify", []string{"read:student"}),
		oauth2.Request("POST", "/createProject", []string{"read:supervisor", "read:student"}),
		oauth2.Request("POST", "/createApplication", []string{"read:student"}),
		oauth2.Request("PATCH", "/acceptApplication", []string{"read:supervisor"}),
		oauth2.Request("PATCH", "/declineApplication", []string{"read:supervisor"}),
		oauth2.Request("POST", "/createGanttItem", []string{"read:supervisor", "read:student"}),
		oauth2.Request("PATCH", "/updateFeedback/:feedback", []string{"read:supervisor"}),
		oauth2.Request("DELETE", "/deleteGanttItem/:id", []string{"read:supervisor", "read:student"}),
		oauth2.Request("POST", "/createSupervisorUser", []string{"read:admin"}),
		oauth2.Request("POST", "/createStudentUser", []string{"read:student"}),
		oauth2.Request("POST", "/completeGanttItem", []string{"read:student"}),
	)
	if err != nil {
		log.Printf("cannot create OAuth2 middleware: %v", err)
		os.Exit(2)
	}

	app.Use(oauth2.New(oauth2Config))

	app.Post("/authorize", controller.AuthorizeHandler)
	app.Post("/newQuestion", controller.NewQuestion) //creates new question
	app.Post("/newAnswer", controller.NewAnswer)     //creates new answer for particular question and adds to db
	app.Get("/getQuestions", controller.GetQuestionsHandler)
	//app.Get("/isSupervisor", controller.GetSupervisorHandler)
	app.Get("/getApplications", controller.GetApplicationsHandler)                     //retrieves all applications from db
	app.Get("/getApplicationsForStudent", controller.GetApplicationsForStudentHandler) //retrieves all applications from db
	app.Get("/GetAllAcceptedRequests", controller.GetAllAcceptedRequestsHandler)       //retrieves all applications from db
	app.Get("/getSpecificApplications/:id", controller.GetSpecificApplicationsHandler) //retrieves one specific applications
	app.Get("/getGanttItem/:id", controller.GetGanttItem)
	app.Get("/getGantt/:id", controller.GetGantt)
	app.Get("/getSupervisors", controller.GetSupervisorHandler)
	app.Get("/getProjectStatus", controller.GetHasProjectStatusHandler)
	app.Get("/getProjects", controller.GetProjectsHandler)
	app.Get("/getProjectID", controller.GetProjectIDHandler)
	app.Get("/getFeedback/:id", controller.GetFeedback)
	app.Get("/getProjectName/:id", controller.GetProjectNameHandler)
	app.Get("/getUsername/:id", controller.GetUsernameHandler)
	app.Get("/verify", controller.VerifyHandler)
	app.Post("/createProject", controller.CreateProjectHandler)         //post createproject
	app.Post("/createApplication", controller.CreateApplicationHandler) //post createapplication
	app.Post("/createSupervisorUser", controller.CreateSupervisorHandler)
	//patch acceptapplication
	app.Patch("/declineApplication", controller.DeclineApplicationHandler) //patch declineapplication
	app.Patch("/completeGanttItem", controller.CompleteGanttItemHandler)
	app.Post("/createGanttItem", controller.CreateGanttItemHandler) //creates Gantt item in db
	app.Patch("/updateFeedback/:id ", controller.AddFeedbackHandler)
	app.Post("/createStudentUser", controller.CreateStudentHandler)
	app.Delete("/deleteGanttItem/:id", controller.DeleteGanttItemHandler)
	app.Listen(":3000")
}

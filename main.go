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
	programLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))

	allowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = ""
	}

	allowMethods := os.Getenv("CORS_ALLOW_METHODS")
	if allowMethods == "" {
		allowMethods = ""
	}

	allowHeaders := os.Getenv("CORS_ALLOW_HEADERS")
	if allowHeaders == "" {
		allowHeaders = ""
	}

	dbClient := db.MustCreate(os.Getenv("DB_URL"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")) // create db client

	debug := false
	debugStr := os.Getenv("DEBUG")
	if debugStr != "" {
		if strings.ToUpper(debugStr) == "TRUE" {
			debug = true
		}
	}

	if debug {
		programLevel.Set(slog.LevelDebug)
	}
	slog.SetDefault(logger)
	slog.Debug("debugging mode is on")

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
	//oauth2Config, err := oauth2.Build(
	//	oauth2.Debug(debug),
	//	oauth2.URL(os.Getenv("JWKS_URL")),
	//	oauth2.Unmatched(true),
	//	oauth2.Audience(os.Getenv("AUDIENCE")),
	//	oauth2.Issuer(os.Getenv("ISSUER")),
	//	oauth2.HTTPClient(&http.Client{}),
	//	oauth2.Request("GET", "/questions", []string{"read:questions"}),
	//	oauth2.Request("POST", "/createApplication", []string{"read:student"}),
	//	oauth2.Request("POST", "/newQuestion", []string{"read:student"}),
	//	oauth2.Request("POST", "/newAnswer", []string{"read:supervisor"}),
	//	oauth2.Request("GET", "/getQuestions", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("GET", "/getApplications", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("GET", "/GetAllAcceptedRequests", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("GET", "/getApplicationsForStudent", []string{"read:student"}),
	//	oauth2.Request("GET", "/getSpecificApplications", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("GET", "/getGanttItem/:id", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("GET", "/getGantt/:id", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("GET", "/getSupervisors", []string{"read:student"}),
	//	oauth2.Request("GET", "/getProjects", []string{"read:student", "read:supervisor"}),
	//	oauth2.Request("GET", "/getUsername/:id", []string{"read:student", "read:supervisor"}),
	//	oauth2.Request("GET", "/getProjectID", []string{"read:student"}),
	//	oauth2.Request("GET", "/getProjectStatus", []string{"read:student"}),
	//	oauth2.Request("GET", "/getProjectName/:id", []string{"read:student", "read:supervisor"}),
	//	oauth2.Request("GET", "/getGantt/:id", []string{"read:student", "read:supervisor"}),
	//	oauth2.Request("GET", "/verify", []string{"read:student"}),
	//	oauth2.Request("POST", "/createProject", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("POST", "/createApplication", []string{"read:student"}),
	//	oauth2.Request("PATCH", "/acceptApplication", []string{"read:supervisor"}),
	//	oauth2.Request("PATCH", "/declineApplication", []string{"read:supervisor"}),
	//	oauth2.Request("PATCH", "/disableAlert/:id", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("POST", "/createGanttItem", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("PATCH", "/updateFeedback", []string{"read:student", "read:supervisor"}),
	//	oauth2.Request("DELETE", "/deleteGanttItem/:id", []string{"read:supervisor", "read:student"}),
	//	oauth2.Request("POST", "/createSupervisorUser", []string{"read:admin"}),
	//	oauth2.Request("POST", "/createStudentUser", []string{"read:student"}),
	//	oauth2.Request("POST", "/completeGanttItem", []string{"read:student"}),
	//)
	if err != nil {
		log.Printf("cannot create OAuth2 middleware: %v", err)
		os.Exit(2)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: allowMethods,
		AllowHeaders: allowHeaders,
	}))

	//app.Use(oauth2.New(oauth2Config))

	app.Post("/authorize", controller.AuthorizeHandler)
	app.Patch("/disableAlert/:id", MustGetOAuth2Handler("PATCH", "/disableAlert/:id", []string{"read:supervisor", "read:student"}, debug), controller.DisableAlertHandler)

	app.Post("/newQuestion", MustGetOAuth2Handler("POST", "/newQuestion", []string{"read:student"}, debug), controller.NewQuestion)              //creates new question
	app.Post("/newAnswer", MustGetOAuth2Handler("POST", "/newAnswer", []string{"read:supervisor", "read:student"}, debug), controller.NewAnswer) //creates new answer for particular question and adds to db
	app.Get("/getQuestions", MustGetOAuth2Handler("GET", "/getQuestions", []string{"read:supervisor", "read:student"}, debug), controller.GetQuestionsHandler)
	//app.Get("/isSupervisor", controller.GetSupervisorHandler)
	app.Get("/getApplications", MustGetOAuth2Handler("GET", "/getApplications", []string{"read:supervisor", "read:student"}, debug), controller.GetApplicationsHandler)                                 //retrieves all applications from db
	app.Get("/getApplicationsForStudent", MustGetOAuth2Handler("POST", "/getApplicationsForStudent", []string{"read:student"}, debug), controller.GetApplicationsForStudentHandler)                     //retrieves all applications from db
	app.Get("/getAllAcceptedRequests", MustGetOAuth2Handler("GET", "/getAllAcceptedRequests", []string{"read:supervisor", "read:student"}, debug), controller.GetAllAcceptedRequestsHandler)            //retrieves all applications from db
	app.Get("/getSpecificApplications/:id", MustGetOAuth2Handler("GET", "/getSpecificApplications/:id", []string{"read:supervisor", "read:student"}, debug), controller.GetSpecificApplicationsHandler) //retrieves one specific applications
	app.Get("/getGanttItem/:id", MustGetOAuth2Handler("GET", "/getGanttItem/:id", []string{"read:supervisor", "read:student"}, debug), controller.GetGanttItem)
	app.Get("/getGantt/:id", MustGetOAuth2Handler("GET", "/getGantt/:id", []string{"read:supervisor", "read:student"}, debug), controller.GetGantt)
	app.Get("/getSupervisors", MustGetOAuth2Handler("GET", "/getSupervisors", []string{"read:supervisor", "read:student"}, debug), controller.GetSupervisorHandler)
	app.Get("/getProjectStatus", MustGetOAuth2Handler("GET", "/getProjectStatus", []string{"read:supervisor", "read:student"}, debug), controller.GetHasProjectStatusHandler)
	app.Get("/getProjects", MustGetOAuth2Handler("GET", "/getProjects", []string{"read:supervisor", "read:student"}, debug), controller.GetProjectsHandler)
	app.Get("/getProjectID", MustGetOAuth2Handler("GET", "/getProjectID", []string{"read:supervisor", "read:student"}, debug), controller.GetProjectIDHandler)
	app.Get("/getFeedback/:id", MustGetOAuth2Handler("GET", "/getFeedback/:id", []string{"read:supervisor", "read:student"}, debug), controller.GetFeedback)
	app.Get("/getProjectName/:id", MustGetOAuth2Handler("GET", "/getProjectName/:id", []string{"read:supervisor", "read:student"}, debug), controller.GetProjectNameHandler)
	app.Get("/getUsername/:id", MustGetOAuth2Handler("GET", "/getUsername/:id", []string{"read:supervisor", "read:student"}, debug), controller.GetUsernameHandler)
	app.Get("/verify", MustGetOAuth2Handler("GET", "/verify", []string{"read:supervisor", "read:student"}, debug), controller.VerifyHandler)
	app.Post("/createProject", MustGetOAuth2Handler("POST", "/createProject", []string{"read:supervisor", "read:student"}, debug), controller.CreateProjectHandler)             //post createproject
	app.Post("/createApplication", MustGetOAuth2Handler("POST", "/createApplication", []string{"read:supervisor", "read:student"}, debug), controller.CreateApplicationHandler) //post createapplication
	app.Post("/createSupervisorUser", MustGetOAuth2Handler("POST", "/createSupervisorUser", []string{"read:supervisor", "read:student"}, debug), controller.CreateSupervisorHandler)
	//patch acceptapplication
	app.Patch("/declineApplication", MustGetOAuth2Handler("PATCH", "/declineApplication", []string{"read:supervisor", "read:student"}, debug), controller.DeclineApplicationHandler) //patch declineapplication
	app.Patch("/addSecondReader/:id", MustGetOAuth2Handler("PATCH", "/addSecondReader/:id", []string{"read:supervisor"}, debug), controller.DeclineApplicationHandler)               //patch declineapplication
	app.Patch("/completeGanttItem", MustGetOAuth2Handler("PATCH", "/completeGanttItem", []string{"read:supervisor", "read:student"}, debug), controller.CompleteGanttItemHandler)
	app.Post("/createGanttItem", MustGetOAuth2Handler("POST", "/createGanttItem", []string{"read:supervisor", "read:student"}, debug), controller.CreateGanttItemHandler) //creates Gantt item in db
	app.Patch("/updateFeedback", MustGetOAuth2Handler("PATCH", "/updateFeedback", []string{"read:supervisor", "read:student"}, debug), controller.AddFeedbackHandler)
	app.Post("/createStudentUser", MustGetOAuth2Handler("POST", "/createStudentUser", []string{"read:supervisor", "read:student"}, debug), controller.CreateStudentHandler)
	app.Delete("/deleteGanttItem/:id", MustGetOAuth2Handler("DELETE", "/deleteGanttItem/:id", []string{"read:supervisor", "read:student"}, debug), controller.DeleteGanttItemHandler)

	app.Listen(":3000")
}

func MustGetOAuth2Handler(method, path string, permissions []string, debug bool) fiber.Handler {
	oauth2Config, err := oauth2.Build(
		oauth2.Debug(debug),
		oauth2.URL(os.Getenv("JWKS_URL")),
		oauth2.Unmatched(true),
		oauth2.Audience(os.Getenv("AUDIENCE")),
		oauth2.Issuer(os.Getenv("ISSUER")),
		oauth2.HTTPClient(&http.Client{}),
		oauth2.Request(method, path, permissions))

	if err != nil {
		log.Fatalf("cannot build OAuth2 middleware: %v", err)
	}
	return oauth2.New(oauth2Config)
}

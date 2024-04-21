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

	//Initialize oauth2 middleware
	oauth2Config, err := oauth2.Build(
		oauth2.Debug(debug),
		oauth2.URL(os.Getenv("JWKS_URL")),
		oauth2.Unmatched(true),
		oauth2.Audience(os.Getenv("AUDIENCE")),
		oauth2.Issuer(os.Getenv("ISSUER")),
		oauth2.HTTPClient(&http.Client{}))
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

	app.Post("/authorize", controller.AuthorizeHandler)
	app.Patch("/disableAlert/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.DisableAlertHandler)

	app.Post("/newQuestion", oauth2Config.Authorize([]string{"read:student"}), controller.NewQuestion)                //creates new question
	app.Post("/newAnswer", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.NewAnswer) //creates new answer for particular question and adds to db
	app.Get("/getQuestions", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetQuestionsHandler)
	//app.Get("/isSupervisor", controller.GetSupervisorHandler)
	app.Get("/getApplications", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetApplicationsHandler)                     //retrieves all applications from db
	app.Get("/getApplicationsForStudent", oauth2Config.Authorize([]string{"read:student"}), controller.GetApplicationsForStudentHandler)                    //retrieves all applications from db
	app.Get("/getAllAcceptedRequests", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetAllAcceptedRequestsHandler)       //retrieves all applications from db
	app.Get("/getSpecificApplications/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetSpecificApplicationsHandler) //retrieves one specific applications
	app.Get("/getGanttItem/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetGanttItem)
	app.Get("/getGantt/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetGantt)
	app.Get("/getSupervisors", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetSupervisorHandler)
	app.Get("/getProjectStatus", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetHasProjectStatusHandler)
	app.Get("/getProjects", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetProjectsHandler)
	app.Get("/getProjectID", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetProjectIDHandler)
	app.Get("/getFeedback/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetFeedback)
	app.Get("/getProjectName/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetProjectNameHandler)
	app.Get("/getUsername/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.GetUsernameHandler)
	app.Get("/getSecondProjects", oauth2Config.Authorize([]string{"read:supervisor"}), controller.GetSecondProjectsHandler)
	app.Get("/getSecondReaderStatus/:id", oauth2Config.Authorize([]string{"read:student", "read:supervisor"}), controller.GetSecondReaderStatusHandler)
	app.Get("/verify", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.VerifyHandler)
	app.Post("/createProject", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.CreateProjectHandler)         //post createproject
	app.Post("/createApplication", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.CreateApplicationHandler) //post createapplication
	app.Post("/createSupervisorUser", oauth2Config.Authorize([]string{"read:admin"}), controller.CreateSupervisorHandler)
	//patch acceptapplication
	app.Patch("/declineApplication", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.DeclineApplicationHandler) //patch declineapplication
	app.Patch("/addSecondReader/:id", oauth2Config.Authorize([]string{"read:supervisor"}), controller.AddSecondReaderHandler)                   //patch declineapplication
	app.Patch("/completeGanttItem", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.CompleteGanttItemHandler)
	app.Post("/createGanttItem", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.CreateGanttItemHandler) //creates Gantt item in db
	app.Patch("/updateFeedback", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.AddFeedbackHandler)
	app.Post("/createStudentUser", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.CreateStudentHandler)
	app.Delete("/deleteGanttItem/:id", oauth2Config.Authorize([]string{"read:supervisor", "read:student"}), controller.DeleteGanttItemHandler)

	app.Listen(":3000")
}

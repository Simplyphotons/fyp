package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Simplyphotons/fyp.git/auth0"
	"github.com/Simplyphotons/fyp.git/db"
	"github.com/Simplyphotons/fyp.git/model"
	"github.com/Simplyphotons/fyp.git/security"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
)

func (c Controller) GetApplicationsHandler(ctx *fiber.Ctx) error { //get all gantt item for a particular project

	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	fmt.Printf("%s\n", authority.UserID)

	response, err := c.dbClient.GetApplications(ctx.Context(), authority.UserID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetAllAcceptedRequestsHandler(ctx *fiber.Ctx) error { //get all gantt item for a particular project

	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	fmt.Printf("%s\n", authority.UserID)

	response, err := c.dbClient.GetAllAcceptedRequests(ctx.Context())
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetApplicationsForStudentHandler(ctx *fiber.Ctx) error { //get all gantt item for a particular project

	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	fmt.Printf("%s\n", authority.UserID)

	response, err := c.dbClient.GetApplicationsForStudent(ctx.Context(), authority.UserID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetSpecificApplicationsHandler(ctx *fiber.Ctx) error { //get all gantt item for a particular project
	response, err := c.dbClient.GetSpecificApplications(ctx.Context(), ctx.Params("id"))
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetUsernameHandler(ctx *fiber.Ctx) error {
	response, err := c.dbClient.GetUsername(ctx.Context(), ctx.Params("id"))
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetProjectsHandler(ctx *fiber.Ctx) error { //get all projects

	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	fmt.Printf("%s\n", authority.UserID)

	response, err := c.dbClient.GetProjects(ctx.Context(), authority.UserID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetProjectIDHandler(ctx *fiber.Ctx) error { //get all projects

	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	fmt.Printf("%s\n", authority.UserID)

	response, err := c.dbClient.GetProjectID(ctx.Context(), authority.UserID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetProjectNameHandler(ctx *fiber.Ctx) error { //get all projects

	response, err := c.dbClient.GetProjectName(ctx.Context(), ctx.Params("id"))
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) CreateApplicationHandler(ctx *fiber.Ctx) error {
	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	fmt.Printf("%s", authority.UserID)

	// Read the request body
	var application model.ApplicationData

	err := json.Unmarshal(ctx.Body(), &application)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	// Translate it to the db request
	applicationRequest := db.Application{
		ID:           application.ID,
		StudentID:    application.StudentID,
		SupervisorID: application.SupervisorID,
		Heading:      application.Heading,
		Description:  application.Description,
		Accepted:     false,
		Declined:     false,
	}
	println(ctx)
	// Execute db request
	err = c.dbClient.CreateApplication(ctx.Context(), applicationRequest, authority.UserID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) CreateSupervisorHandler(ctx *fiber.Ctx) error {
	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	println(authority.UserID)

	var fromFrontEndRequest model.UserCreateRequest

	err := json.Unmarshal(ctx.Body(), &fromFrontEndRequest)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	doesExist, err := c.auth0Client.DoesUserExist(ctx.Context(), fromFrontEndRequest.Email)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(http.StatusInternalServerError).JSON(message)
	}
	if doesExist {
		errorMessage := model.ErrorMessage{
			Message: fmt.Sprintf("User '%s' already registerd in auth0", fromFrontEndRequest.Email),
		}
		return ctx.Status(http.StatusBadRequest).JSON(errorMessage)
	}

	authRequest := auth0.UserCreateRequest{
		Email:         fromFrontEndRequest.Email,
		Password:      fromFrontEndRequest.Password,
		FirstName:     &fromFrontEndRequest.FirstName,
		LastName:      &fromFrontEndRequest.LastName,
		Name:          fromFrontEndRequest.FirstName + " " + fromFrontEndRequest.LastName,
		VerifyEmail:   false,
		Connection:    "Username-Password-Authentication",
		EmailVerified: true,
	}

	newSupervisorID, err := c.auth0Client.AddUser(ctx.Context(), authRequest) //returns new user id
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(http.StatusInternalServerError).JSON(message)
	}
	err = c.auth0Client.AddRole(ctx.Context(), newSupervisorID, c.supervisorRoleID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(http.StatusInternalServerError).JSON(message)
	}

	// Translate it to the db request
	userRequest := db.User{
		Id:   newSupervisorID,
		Name: authRequest.Name,
	}
	// Execute db request
	err = c.dbClient.CreateSupervisorUser(ctx.Context(), userRequest)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) CreateStudentHandler(ctx *fiber.Ctx) error {
	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	var user model.UserData
	err := json.Unmarshal(ctx.Body(), &user)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	userRequest := db.User{
		Id:   authority.UserID,
		Name: user.Name,
	}

	err = c.dbClient.CreateStudentUser(ctx.Context(), userRequest)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

//func (c Controller) AcceptApplicationHandler(ctx *fiber.Ctx) error {
//	var application model.ApplicationData
//
//	err := json.Unmarshal(ctx.Body(), &application)
//	if err != nil {
//		message := model.ErrorMessage{
//			Message: err.Error(),
//		}
//		return ctx.Status(400).JSON(message)
//	}
//
//	// Translate it to the db request
//	applicationRequest := db.Application{
//		ID:           application.ID,
//		StudentID:    application.StudentID,
//		SupervisorID: application.SupervisorID,
//		Heading:      application.Heading,
//		Description:  application.Description,
//		Accepted:     false,
//		Declined:     false,
//	}
//
//	// Execute db request
//	err = c.dbClient.AcceptApplication(ctx.Context(), applicationRequest)
//	if err != nil {
//		message := model.ErrorMessage{
//			Message: err.Error(),
//		}
//		return ctx.Status(500).JSON(message)
//	}
//
//	return ctx.SendStatus(204)
//}

func (c Controller) DeclineApplicationHandler(ctx *fiber.Ctx) error {
	var application model.ApplicationData

	err := json.Unmarshal(ctx.Body(), &application)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	// Translate it to the db request
	applicationRequest := db.Application{
		ID:           application.ID,
		StudentID:    application.StudentID,
		SupervisorID: application.SupervisorID,
		Heading:      application.Heading,
		Description:  application.Description,
		Accepted:     false,
		Declined:     false,
	}

	// Execute db request
	err = c.dbClient.DeclineApplication(ctx.Context(), applicationRequest)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) VerifyHandler(ctx *fiber.Ctx) error {
	var (
		authority security.Authority
		ok        bool
	)
	if authority, ok = ctx.UserContext().Value(security.AuthorityKey{}).(security.Authority); !ok {
		message := model.ErrorMessage{
			Message: "cannot extract user id",
		}

		return ctx.Status(401).JSON(message)
	}

	response, err := c.dbClient.Verify(ctx.Context(), authority.UserID)
	if err != nil {
		slog.Error("cannot verify student account", "error", err)
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

package handlers

import (
	"encoding/json"
	"github.com/Simplyphotons/fyp.git/db"
	"github.com/Simplyphotons/fyp.git/model"
	"github.com/Simplyphotons/fyp.git/security"
	"github.com/gofiber/fiber/v2"
)

func (c Controller) GetGanttItem(ctx *fiber.Ctx) error {
	response, err := c.dbClient.GetGanttItem(ctx.Context(), ctx.Params("id"))
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetGantt(ctx *fiber.Ctx) error { //get all gantt item for a particular project
	response, err := c.dbClient.GetGantt(ctx.Context(), ctx.Params("id"))
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) GetFeedback(ctx *fiber.Ctx) error { //get all gantt item for a particular project
	response, err := c.dbClient.GetFeedback(ctx.Context(), ctx.Params("id"))
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}
	return ctx.Status(200).JSON(response)
}

func (c Controller) CreateProjectHandler(ctx *fiber.Ctx) error {

	// Read the request body
	var application model.ApplicationData

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

	err := json.Unmarshal(ctx.Body(), &application)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	// Translate it to the db request
	projectRequest := db.Application{
		ID:           application.ID,
		StudentName:  application.StudentID, //these two swap for some reason to reswapping fixes the issue
		StudentID:    application.StudentName,
		SupervisorID: authority.UserID,
	}

	// Execute db request
	err = c.dbClient.CreateProject(ctx.Context(), projectRequest, authority.UserID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) CreateGanttItemHandler(ctx *fiber.Ctx) error {

	// Read the request body
	var gantt model.Gantt

	err := json.Unmarshal(ctx.Body(), &gantt)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	// Translate it to the db request
	ganttRequest := db.Gantt{
		Id:          gantt.ID,
		ProjectID:   gantt.ProjectID,
		StartDate:   gantt.StartDate,
		EndDate:     gantt.EndDate,
		Description: gantt.Description,
		Links:       gantt.Links,
		Feedback:    gantt.Feedback,
		GanttName:   gantt.GanttName,
	}

	// Execute db request
	err = c.dbClient.CreateGanttItem(ctx.Context(), ganttRequest)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) AddFeedbackHandler(ctx *fiber.Ctx) error {
	// Read the request body
	var gantt model.Gantt

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

	err := json.Unmarshal(ctx.Body(), &gantt)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	// Translate it to the db request

	ganttRequest := db.Gantt{
		Id:          gantt.ID,
		Feedback:    gantt.Feedback,
		NewFeedBack: gantt.NewFeedback,
	}

	err = c.dbClient.UpdateFeedback(ctx.Context(), ganttRequest, authority.UserID)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) DisableAlertHandler(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

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

	err := c.dbClient.DisableAlert(ctx.Context(), authority.UserID, id)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) DeleteGanttItemHandler(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	err := c.dbClient.DeleteGanttItem(id)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

func (c Controller) CompleteGanttItemHandler(ctx *fiber.Ctx) error {

	var gantt model.Gantt

	err := json.Unmarshal(ctx.Body(), &gantt)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	// Translate it to the db request
	ganttRequest := db.Gantt{
		Id: gantt.ID,
	}

	err = c.dbClient.CompleteGanttItem(ctx.Context(), ganttRequest)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

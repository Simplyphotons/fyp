package handlers

import (
	"FYP/db"
	"FYP/model"
	"encoding/json"
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

func (c Controller) CreateProjectHandler(ctx *fiber.Ctx) error {

	// Read the request body
	var project model.ProjectData

	err := json.Unmarshal(ctx.Body(), &project)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(400).JSON(message)
	}

	// Translate it to the db request
	projectRequest := db.Project{
		ID:           project.ID,
		Name:         project.Name,
		StudentID:    project.StudentID,
		SupervisorID: project.SupervisorID,
		CreatedAt:    project.CreatedAt,
	}

	// Execute db request
	err = c.dbClient.CreateProject(ctx.Context(), projectRequest)
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

	feedback := ctx.Params("feedback")

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
	}

	// Execute db request
	err = c.dbClient.UpdateFeedback(ctx.Context(), ganttRequest, gantt.ID, feedback)
	if err != nil {
		message := model.ErrorMessage{
			Message: err.Error(),
		}
		return ctx.Status(500).JSON(message)
	}

	return ctx.SendStatus(204)
}

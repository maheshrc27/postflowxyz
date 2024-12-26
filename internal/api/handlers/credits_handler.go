package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maheshrc27/postflow/internal/service"
)

type CreditsHandler struct {
	c service.CreditsService
}

func NewCreditsHandler(service service.CreditsService) *CreditsHandler {
	return &CreditsHandler{c: service}
}

func (h *CreditsHandler) GetCredits(c *fiber.Ctx) error {
	userId := GetUserID(c)

	credits, err := h.c.GetCredits(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Couldn't find user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"credits": credits,
	})
}

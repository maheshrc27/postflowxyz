package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maheshrc27/postflow/internal/service"
	"github.com/maheshrc27/postflow/internal/transfer"
)

type VideoHandler struct {
	v service.VideoService
}

func NewVideoHandler(vs service.VideoService) *VideoHandler {
	return &VideoHandler{v: vs}
}

func (h *VideoHandler) GetVideos(c *fiber.Ctx) error {
	userId := GetUserID(c)

	videos, err := h.v.GetVideos(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to get videos",
		})
	}

	return c.Status(fiber.StatusOK).JSON(videos)
}

func (h *VideoHandler) CreateVideo(c *fiber.Ctx) error {
	userId := GetUserID(c)

	var video transfer.VideoTransfer
	err := c.BodyParser(&video)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse request",
		})
	}

	videoURL, err := h.v.RequestVideo(c.Context(), userId, string(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to generate video",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"video_url": videoURL,
	})

}

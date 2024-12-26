package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/maheshrc27/postflow/internal/service"
)

type PaymentHandler struct {
	c service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{c: service}
}

func (h *PaymentHandler) PaymentWebhook(c *fiber.Ctx) error {
	customerEmail := c.FormValue("email")
	productId := c.FormValue("short_product_id")
	productPrice := c.FormValue("price")

	if customerEmail == "" || productId == "" {
		log.Println("Error: Missing email or product_id in webhook payload")
		return c.Status(fiber.StatusBadRequest).SendString("Email or product_id is empty")
	}

	err := h.c.HandlePayment(c.Context(), customerEmail, productId, productPrice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong while saving account")
	}

	return c.SendStatus(fiber.StatusOK)
}

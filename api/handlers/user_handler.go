package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karsteneugene/top-up-system/api/models"
)

// GetUsers godoc
// @Summary Get users
// @Description Get users
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Router /users [get]
func GetUsers(c *fiber.Ctx) error {
	user := models.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
	}
	return c.JSON(user)
}

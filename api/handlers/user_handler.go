package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karsteneugene/top-up-system/api/models"
)

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Produce json
// @Success 200 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users [get]
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	// Get all users and check if there are any problems with the database
	if err := db.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error retrieving users"})
	}
	// Check if users is empty
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No users found"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"users": users}})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Router /users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	var user models.User
	id := c.Params("id")

	// Query user by ID
	if err := db.First(&user, id).Error; err != nil { // Check if user exists
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"user": user}})
}

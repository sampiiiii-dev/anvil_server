package resources

import (
	"github.com/labstack/echo/v4"
	"github.com/sampiiiii-dev/anvil_server/anvil/common" // Import the common package
	"github.com/sampiiiii-dev/anvil_server/anvil/models"
	"gorm.io/gorm"
	"net/http"
)

// RegisterUserRoutes registers RESTful CRUD routes for the User model.
func RegisterUserRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/users", func(c echo.Context) error {
		// List all users
		var users []models.User
		db.Find(&users)
		common.CustomResponse(c.Response().Writer, http.StatusOK, users, "List of users")
		return nil
	})

	e.GET("/users/:id", func(c echo.Context) error {
		// Get a single user by ID
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			common.CustomResponse(c.Response().Writer, http.StatusNotFound, "User not found", nil)
			return err
		}
		common.CustomResponse(c.Response().Writer, http.StatusOK, user, "User details")
		return nil
	})

	e.POST("/users", func(c echo.Context) error {
		// Create a new user
		user := new(models.User)
		if err := c.Bind(user); err != nil {
			common.CustomResponse(c.Response().Writer, http.StatusBadRequest, "Invalid input", nil)
			return err
		}
		db.Create(&user)
		common.CustomResponse(c.Response().Writer, http.StatusCreated, user, "User created")
		return nil
	})

	e.PUT("/users/:id", func(c echo.Context) error {
		// Update an existing user
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			common.CustomResponse(c.Response().Writer, http.StatusNotFound, "User not found", nil)
			return err
		}
		if err := c.Bind(&user); err != nil {
			common.CustomResponse(c.Response().Writer, http.StatusBadRequest, "Invalid input", nil)
			return err
		}
		db.Save(&user)
		common.CustomResponse(c.Response().Writer, http.StatusOK, user, "User updated")
		return nil
	})

	e.DELETE("/users/:id", func(c echo.Context) error {
		// Delete a user
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			common.CustomResponse(c.Response().Writer, http.StatusNotFound, "User not found", nil)
			return err
		}
		db.Delete(&user)
		common.CustomResponse(c.Response().Writer, http.StatusNoContent, nil, "User deleted")
		return nil
	})
}

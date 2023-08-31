package controller

import (
	"fmt"
	"net/http"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
)

// User

// GetUsers             godoc
// @Summary      Get users array
// @Description  Responds with the list of all users as JSON.
// @Tags         users
// @Produce      json
// @Success      200  {array}  models.User
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	users := []models.User{}
	if db.DB.Find(&users).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server side error"})
		return
	}
	c.JSON(http.StatusOK, &users)
}

// CreateUser             godoc
// @Summary      Create a new user
// @Description  Takes no arguments.
// @Tags         users
// @Produce      json
// @Success      200 {object} models.User
// @Router       /user [post]
func CreateUser(c *gin.Context) {
	var user models.User
	db.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %d was createdd", user.ID)})
}

// DeleteUser             godoc
// @Summary      Delete a user
// @Description  Takes user id that you want to delete.
// @Tags         users
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} models.User
// @Router       /user/:id [delete]
func DeleteUser(c *gin.Context) {
	var user models.User
	db.DB.Where("id = ?", c.Param("id")).Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %d was deleted", user.ID)})
}

// UpdateUser             godoc
// @Summary      Change user id
// @Description  Takes user id that you want to change as a parameter and in json user id as "id" that you want to change to.
// @Tags         users
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} models.User
// @Router       /user/:id [put]
func UpdateUser(c *gin.Context) {
	var user models.User
	db.DB.Where("id = ?", c.Param("id")).First(&user)
	c.BindJSON(&user)
	db.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s's id was updated to %d", c.Param("id"), user.ID)})
}

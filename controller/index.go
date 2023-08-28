package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index             godoc
// @Summary      Start page
// @Description  Just a start page with a link to API docs
// @Success      200
// @Router       / [get]
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "go to /docs/index.html for API docs"})
}

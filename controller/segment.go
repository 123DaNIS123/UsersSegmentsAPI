package controller

import (
	"fmt"
	"net/http"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
)

// GetSegments             godoc
// @Summary      Get segments array
// @Description  Responds with the list of all segments as JSON.
// @Tags         segments
// @Produce      json
// @Success      200  {array}  models.Segment
// @Router       /segments [get]
func GetSegments(c *gin.Context) {
	segments := []models.Segment{}
	db.DB.Find(&segments)
	c.JSON(http.StatusOK, &segments)
}

// GetSegment             godoc
// @Summary      Get a segment by id
// @Description  Takes segment id as a parameter.
// @Tags         segments
// @Produce      json
// @Param        id path int true "Segment ID"
// @Success      200 {object} models.Segment
// @Router       /segment/:id [get]
func GetSegment(c *gin.Context) {
	var segment models.Segment
	db.DB.Where("id = ?", c.Param("id")).First(&segment)
	c.JSON(http.StatusOK, &segment)
}

// CreateSegment             godoc
// @Summary      Create a segment
// @Description  Takes in JSON "name" of segment that you want to create.
// @Tags         segments
// @Produce      json
// @Param        name formData string true "Segment name"
// @Param        percentage formData int false "Percentage of users that will be added to this segment"
// @Success      200 {object} models.Segment
// @Router       /segment [post]
func CreateSegment(c *gin.Context) {
	var segment models.Segment
	// c.BindJSON(&segment)
	if err := c.BindJSON(&segment); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var tempSegment models.Segment
	if r := db.DB.Model(&models.Segment{}).Where("name = ?", segment.Name).Limit(1).Find(&tempSegment); r.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Segment with such name already exists"})
		return
	}
	db.DB.Create(&segment)
	if segment.Percentage > 0 {

		tmpUserSegments := make([]models.UserSegment, 0)
		db.DB.Model(&models.UserSegment{}).Find(&tmpUserSegments)

		userSequence(c, segment.Percentage, segment.Name)
	}
}

// DeleteSegment             godoc
// @Summary      Delete a segment
// @Description  Takes in JSON "name" of segment that you want to delete.
// @Tags         segments
// @Produce      json
// @Param        name body string true "Segment name"
// @Success      200 {object} models.Segment
// @Router       /segment [delete]
func DeleteSegment(c *gin.Context) {
	var segment models.Segment
	if c.BindJSON(&segment) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}
	if db.DB.Where("name = ?", segment.Name).Delete(&segment).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Already deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s was deleted", segment.Name)})
}

// UpdateSegment             godoc
// @Summary      Update a segment
// @Description  Takes in JSON "name" you want to change the segment to.
// @Tags         segments
// @Produce      json
// @Param        id path int true "Segment ID"
// @Param        name formData string true "Segment name"
// @Success      200 {object} models.Segment
// @Router       /segment/:id [put]
func UpdateSegment(c *gin.Context) {
	var segment models.Segment

	db.DB.Where("id = ?", c.Param("id")).First(&segment)
	if c.BindJSON(&segment) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}
	db.DB.Save(&segment)
	c.JSON(http.StatusOK, &segment)
}

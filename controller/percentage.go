package controller

import (
	"net/http"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
)

func userSequence(c *gin.Context, percentage uint, segmentName string) {
	var allUsersAmount int64
	var users []models.User
	db.DB.Table("users").Count(&allUsersAmount)
	limitUserAmount := int((float64(percentage) / 100.0) * float64(allUsersAmount))
	if limitUserAmount == 0 {
		c.JSON(http.StatusBadRequest, "low percentage amount")
		return
	}
	if err := db.DB.Table("users").Order("random()").Limit(limitUserAmount).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, &users)
		return
	}
	userSegmentsSlice := make([]models.UserSegment, 0)
	for _, v := range users {
		bindMessage := BindMessage{SegmentsAdd: []string{segmentName}, UserID: v.ID}
		// userSegmentsSlice := bindMessage.Add()
		userSegmentsSlice = append(userSegmentsSlice, bindMessage.AddUserSegments()[0])
	}
	if len(userSegmentsSlice) != 0 {
		db.DB.Save(&userSegmentsSlice)
	}
	c.JSON(http.StatusOK, &userSegmentsSlice)
}

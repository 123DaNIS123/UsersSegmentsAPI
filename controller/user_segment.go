package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BindMessage struct {
	SegmentsAdd    []string `json:"segments_add"`
	SegmentsRemove []string `json:"segments_remove"`
	UserID         uint     `json:"user_id"`
}

func (bm *BindMessage) AddUserSegments() []models.UserSegment {
	userSegmentSlice := make([]models.UserSegment, 0)
	for _, v := range bm.SegmentsAdd {
		var userSegment models.UserSegment
		var segment models.Segment
		if err := db.DB.Where("name = ?", v).Find(&segment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("incorrect name \"%s\" in segments_add\n", v)
			}
		} else {
			var user models.User
			if err := db.DB.Where("id = ?", bm.UserID).Find(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Printf("no user with given user_id \"%d\"\n", bm.UserID)
				}
			} else if db.DB.Where("user_id = ? AND segment_id = ?", bm.UserID, segment.ID).Find(&userSegment).Error == nil {
				fmt.Printf(" binding between User \"%d\" and Segment \"%d\" \"%s\"already exist\n", bm.UserID, segment.ID, segment.Name)
			} else {
				userSegment.UserID = bm.UserID
				userSegment.SegmentID = segment.ID
				userSegmentSlice = append(userSegmentSlice, userSegment)
			}
		}
	}
	return userSegmentSlice
}

func (bm *BindMessage) RemoveUserSegments(c *gin.Context) {
	for _, v := range bm.SegmentsRemove {
		var userSegment models.UserSegment
		var segment models.Segment
		if err := db.DB.Where("name = ?", v).First(&segment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("incorrect name \"%s\" in segments_add\n", v)
			}
		} else {
			var user models.User
			if err := db.DB.Where("id = ?", bm.UserID).Find(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Printf("unregistered user id \"%d\" in user_id\n", bm.UserID)
				}
			} else if err := db.DB.Where("user_id = ? AND segment_id = ?", bm.UserID, segment.ID).Delete(&userSegment).Error; err != nil {
				c.JSON(555, &userSegment)
			}
		}
	}
}

// Bind             godoc
// @Summary      Bind and unbind a user with segments
// @Description  Takes in JSON "segments_add" list, "segments_remove" list and "user_id". Return created User-Segment binds
// @Tags         UserSegment
// @Produce      json
// @Param        segments_add formData array false "array of segments that you want to add"
// @Param        segments_remove formData array false "array of segments that you want to remove"
// @Param        user_id formData int true "user id that you want to add"
// @Success      200
// @Router       /bind [post]
func Bind(c *gin.Context) {
	var bindMessage BindMessage
	c.BindJSON(&bindMessage)
	userSegmentsSlice := bindMessage.AddUserSegments()
	if len(userSegmentsSlice) != 0 {
		db.DB.Save(&userSegmentsSlice)
	}
	if len(bindMessage.SegmentsRemove) != 0 {
		bindMessage.RemoveUserSegments(c)
	}
	c.JSON(http.StatusOK, &userSegmentsSlice)
}

// GetBinds             godoc
// @Summary      Get all binds
// @Description  Get all User-Segment binds
// @Tags         UserSegment
// @Produce      json
// @Success      200
// @Router       /binds [get]
func GetBinds(c *gin.Context) {
	userSegments := []models.UserSegment{}
	db.DB.Find(&userSegments)
	c.JSON(http.StatusOK, &userSegments)
}

type UserRequest struct {
	ID uint `json:"user_id"`
}

// UserBinds             godoc
// @Summary      Show certain user's segment binds
// @Description  Takes in JSON "user_id" to show the user's segment binds.
// @Tags         UserSegment
// @Produce      json
// @Param        user_id formData int true "show segment binds of user with given id"
// @Success      200
// @Router       /userbinds [post]
func GetUserBinds(c *gin.Context) {
	var user UserRequest
	var segments []models.Segment
	c.BindJSON(&user)
	if err := db.DB.Table("user_segments").
		Where("user_id = ? AND deleted_at IS NULL", user.ID).
		Order("segment_id asc").Joins("join segments on segments.id = user_segments.segment_id").
		Select("segments.id", "segments.name", "segments.percentage").
		Find(&segments).Error; err != nil {
		c.JSON(http.StatusBadRequest, &segments)
		return
	}
	c.JSON(http.StatusOK, &segments)
}

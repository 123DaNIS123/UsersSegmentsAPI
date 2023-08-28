package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/123DaNIS123/UsersSegments/config"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	if config.DB.Find(&users).Error != nil {
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
	config.DB.Create(&user)
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
	config.DB.Where("id = ?", c.Param("id")).Delete(&user)
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
	config.DB.Where("id = ?", c.Param("id")).First(&user)
	c.BindJSON(&user)
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s's id was updated to %d", c.Param("id"), user.ID)})
}

// Segment

// GetSegments             godoc
// @Summary      Get segments array
// @Description  Responds with the list of all segments as JSON.
// @Tags         segments
// @Produce      json
// @Success      200  {array}  models.Segment
// @Router       /segments [get]
func GetSegments(c *gin.Context) {
	segments := []models.Segment{}
	config.DB.Find(&segments)
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
	config.DB.Where("id = ?", c.Param("id")).First(&segment)
	c.JSON(http.StatusOK, &segment)
}

// CreateSegment             godoc
// @Summary      Create a segment
// @Description  Takes in JSON "name" of segment that you want to create.
// @Tags         segments
// @Produce      json
// @Param        name body string true "Segment name"
// @Success      200 {object} models.Segment
// @Router       /segment [post]
func CreateSegment(c *gin.Context) {
	var segment models.Segment
	if c.BindJSON(&segment) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}
	// config.DB.Create(&segment)
	if config.DB.Create(&segment).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}
	c.JSON(http.StatusOK, &segment)
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
	if config.DB.Where("name = ?", segment.Name).Delete(&segment).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Already deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s was deleted", segment.Name)})
}

// DeleteSegment             godoc
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
	config.DB.Where("id = ?", c.Param("id")).First(&segment)
	if c.BindJSON(&segment) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}
	config.DB.Save(&segment)
	c.JSON(http.StatusOK, &segment)
}

// Bind User to Segment

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
		if err := config.DB.Where("name = ?", v).First(&segment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("incorrect name \"%s\" in segments_add\n", v)
			}
		} else {
			var user models.User
			if err := config.DB.Where("id = ?", bm.UserID).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Printf("no user with given user_id \"%d\"\n", bm.UserID)
				}
			} else if config.DB.Where("user_id = ? AND segment_id = ?", bm.UserID, segment.ID).First(&userSegment).Error == nil {
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
		if err := config.DB.Where("name = ?", v).First(&segment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("incorrect name \"%s\" in segments_add\n", v)
			}
		} else {
			var user models.User
			if err := config.DB.Where("id = ?", bm.UserID).First(&user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Printf("unregistered user id \"%d\" in user_id\n", bm.UserID)
				}
			} else if err := config.DB.Where("user_id = ? AND segment_id = ?", bm.UserID, segment.ID).Delete(&userSegment).Error; err != nil {
				c.JSON(555, &userSegment)
			}
		}
	}
}

// Bind             godoc
// @Summary      Bind and unbind a user with segments
// @Description  Takes in JSON "segments_add" list, "segments_remove" list and "user_id". Return created User-Segment binds
// @Tags         userSegment
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
		config.DB.Save(&userSegmentsSlice)
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
	config.DB.Find(&userSegments)
	c.JSON(http.StatusOK, &userSegments)
}

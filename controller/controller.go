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

func GetUsers(c *gin.Context) {
	users := []models.User{}
	if config.DB.Find(&users).Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server side error"})
		return
	}
	c.JSON(http.StatusOK, &users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if c.BindJSON(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}
	config.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %d was created", user.ID)})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %d was deleted", user.ID)})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).First(&user)
	c.BindJSON(&user)
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s's id was updated to %d", c.Param("id"), user.ID)})
}


// Segment

func GetSegments(c *gin.Context) {
	segments := []models.Segment{}
	config.DB.Find(&segments)
	c.JSON(http.StatusOK, &segments)
}

func GetSegment(c *gin.Context) {
	var segment models.Segment
	config.DB.Where("id = ?", c.Param("id")).First(&segment)
	c.JSON(http.StatusOK, &segment)
}

func CreateSegment(c *gin.Context) {
	var segment models.Segment
	if c.BindJSON(&segment) != nil{
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

func DeleteSegment(c *gin.Context) {
	var segment models.Segment
	if c.BindJSON(&segment) != nil{
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect data"})
		return
	}
	if config.DB.Where("name = ?", segment.Name).Delete(&segment).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Already deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s was deleted", segment.Name)})
}

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
			} else if config.DB.Where("user_id = ? AND segment_id = ?", bm.UserID, segment.ID).First(&userSegment).Error == nil{
				fmt.Printf(" binding between User \"%d\" and Segment \"%d\" \"%s\"already exist\n", bm.UserID, segment.ID, segment.Name)
			} else{
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

func Bind(c *gin.Context) {
	var bindMessage BindMessage
	c.BindJSON(&bindMessage)
	userSegmentsSlice := bindMessage.AddUserSegments()
	if len(userSegmentsSlice) != 0 {
		config.DB.Save(&userSegmentsSlice)
	}
	bindMessage.RemoveUserSegments(c)
	c.JSON(http.StatusOK, &userSegmentsSlice)
}

func GetBinds(c *gin.Context) {
	userSegments := []models.UserSegment{}
	config.DB.Find(&userSegments)
	c.JSON(http.StatusOK, &userSegments)
}

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
	"github.com/gin-gonic/gin"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type BindMessage struct {
	SegmentsAdd    []string `json:"segments_add"`
	SegmentsRemove []string `json:"segments_remove"`
	UserID         uint     `json:"user_id"`
	TTL            Duration `json:"ttl"`
}

func (bm *BindMessage) AddUserSegments() []models.UserSegment {
	userSegmentSlice := make([]models.UserSegment, 0)
	for _, v := range bm.SegmentsAdd {
		var userSegment models.UserSegment
		var segment models.Segment

		if r := db.DB.Model(&models.Segment{}).Where("name = ?", v).Limit(1).Find(&segment); r.RowsAffected == 0 {
			fmt.Printf("incorrect name \"%s\" in segments_add\n", v)
		} else {
			var user models.User
			if r := db.DB.Model(&models.User{}).Where("id = ?", bm.UserID).Limit(1).Find(&user); r.RowsAffected == 0 {
				fmt.Printf("unregistered user id \"%d\" in user_id\n", bm.UserID)
			}
			if r := db.DB.Model(&models.UserSegment{}).Limit(1).Find(&userSegment, "user_id = ? AND segment_id = ?", bm.UserID, segment.ID); r.RowsAffected > 0 {
				fmt.Printf("binding for User \"%d\" and Segment \"%d\" \"%s\"already exist\n", bm.UserID, segment.ID, segment.Name)
			} else {
				userSegment.UserID = bm.UserID
				userSegment.SegmentID = segment.ID
				userSegment.TTL = bm.TTL.Duration
				userSegment.CreatedAt = time.Now()
				userSegmentSlice = append(userSegmentSlice, userSegment)
				var history models.History = models.History{UserID: bm.UserID,
					SegmentID:   segment.ID,
					SegmentName: segment.Name,
					Operation:   "Add",
					Timestamp:   time.Now(),
				}
				db.DB.Model(&models.History{}).Save(&history)
			}
		}
	}
	return userSegmentSlice
}

func (bm BindMessage) RemoveUserSegments() []models.UserSegment {
	userSegmentsSlice := make([]models.UserSegment, 0)
	for _, v := range bm.SegmentsRemove {
		var userSegments models.UserSegment
		var segment models.Segment
		if r := db.DB.Model(&models.Segment{}).Where("name = ?", v).Limit(1).Find(&segment); r.RowsAffected == 0 {
			fmt.Printf("incorrect name \"%s\" in SegmentsRemove\n", v)
		} else {
			var user models.User
			if r := db.DB.Model(&models.User{}).Where("id = ?", bm.UserID).Limit(1).Find(&user); r.RowsAffected == 0 {
				fmt.Printf("unregistered user id \"%d\" in user_id\n", bm.UserID)
			} else if r := db.DB.Model(&models.UserSegment{}).Where("user_id = ? AND segment_id = ?", user.ID, segment.ID).Limit(1).Find(&userSegments); r.RowsAffected == 0 {
				fmt.Printf("error deleting segment %d from user %d\n", segment.ID, bm.UserID)
			} else {
				userSegmentsSlice = append(userSegmentsSlice, userSegments)
				var history models.History = models.History{UserID: bm.UserID,
					SegmentID:   segment.ID,
					SegmentName: segment.Name,
					Operation:   "Remove",
					Timestamp:   time.Now().UTC(),
				}
				db.DB.Model(&models.History{}).Save(&history)
				db.DB.Model(&models.UserSegment{}).Where("user_id = ? AND segment_id = ?", user.ID, segment.ID).Delete(&userSegments)
			}
		}
	}
	return userSegmentsSlice
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
	tmpUserSegments := make([]models.UserSegment, 0)
	db.DB.Model(&models.UserSegment{}).Find(&tmpUserSegments)
	userSegmentsSlice := bindMessage.AddUserSegments()
	if len(userSegmentsSlice) != 0 {
		db.DB.Save(&userSegmentsSlice)
	}
	userSegmentsSliceRemove := make([]models.UserSegment, 0)
	if len(bindMessage.SegmentsRemove) != 0 {
		userSegmentsSliceRemove = bindMessage.RemoveUserSegments()
	}
	c.JSON(http.StatusOK, gin.H{"Added": &userSegmentsSlice, "Removed": &userSegmentsSliceRemove})
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

	tmpUserSegments := make([]models.UserSegment, 0)
	db.DB.Model(&models.UserSegment{}).Find(&tmpUserSegments)

	c.BindJSON(&user)
	fmt.Printf("w.Body: %v\n", &user)
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

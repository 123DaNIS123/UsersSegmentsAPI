package controller

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/123DaNIS123/UsersSegments/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Operation string

const (
	Create Operation = "Create"
	Delete Operation = "Delete"
)

type History struct {
	UserID    uint
	SegmentID uint
	Segment   string
	Operation Operation
	Timestamp time.Time
}

type Data struct {
	UserID      uint           `json:"user_id" gorm:"column:user_id"`
	SegmentID   uint           `json:"segment_id" gorm:"column:segment_id"`
	SegmentName string         `json:"segmentName" gorm:"column:name"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

type YearMonth struct {
	Year  uint `json:"year"`
	Month uint `json:"month"`
	// Day   uint `json:"day"`
}

// History             godoc
// @Summary      Get all user-segment binds since date
// @Description  Takes in JSON "year" and "month" since which to get user-segment binds
// @Tags         UserSegment
// @Produce      json
// @Param        year formData int true "year"
// @Param        month formData int true "month"
// @Success      200
// @Router       /timedata [post]
func GetTimeData(c *gin.Context) {
	var yearMonth YearMonth
	c.BindJSON(&yearMonth)
	if yearMonth.Year == 0 {
		c.JSON(http.StatusBadRequest, "key \"year\" not correct")
		return
	}
	if yearMonth.Month == 0 || yearMonth.Month > 12 {
		c.JSON(http.StatusBadRequest, "key \"month\" not correct")
		return
	}
	userSegments := make([]Data, 0)
	config.DB.Table("user_segments").Find(&userSegments)
	if err := config.DB.Table("user_segments").Unscoped().
		Select("user_segments.user_id", "user_segments.segment_id", "user_segments.created_at", "user_segments.deleted_at", "segments.name").
		Joins("join segments on segments.id = user_segments.segment_id").
		Where("(EXTRACT('Year' FROM created_at) = ? AND EXTRACT('Month' FROM created_at) = ?) OR (EXTRACT('Year' FROM deleted_at) = ? AND EXTRACT('Month' FROM deleted_at) = ?)",
			yearMonth.Year,
			yearMonth.Month,
			yearMonth.Year,
			yearMonth.Month).
		Find(&userSegments).Error; err != nil {
		c.JSON(http.StatusBadRequest, &userSegments)
	} else {
		histories := make([]History, 0)
		for _, v := range userSegments {
			history := History{v.UserID, v.SegmentID, v.SegmentName, Operation(Create), v.CreatedAt.UTC()}
			histories = append(histories, history)
			if v.DeletedAt.Valid {
				history.Operation = Operation(Delete)
				history.Timestamp = v.DeletedAt.Time.UTC()
				histories = append(histories, history)
			}
		}
		sort.Slice(histories, func(i, j int) bool {
			return histories[i].Timestamp.Before(histories[i].Timestamp)
		})
		fileName := fmt.Sprintf("%d-%d.csv", yearMonth.Year, yearMonth.Month)
		filePath := fmt.Sprintf("./%s", fileName)

		file, err := os.Create(filePath)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()
		headers := []string{"N", "UserID", "SegmentId", "SegmentName", "Operation", "Timestamp"}
		writer.Write(headers)
		for i, history := range histories {
			row := []string{
				strconv.Itoa(i),
				strconv.FormatUint(uint64(history.UserID), 10),
				strconv.FormatUint(uint64(history.SegmentID), 10),
				history.Segment,
				string(history.Operation),
				history.Timestamp.Format("02.01.2006 15:04:05"),
			}
			if err := writer.Write(row); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"report_url": filePath})
	}
}

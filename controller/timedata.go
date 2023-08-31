package controller

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/123DaNIS123/UsersSegments/db"
	"github.com/123DaNIS123/UsersSegments/models"
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
	histories := make([]models.History, 0)

	tmpUserSegments := make([]models.UserSegment, 0)
	db.DB.Model(&models.UserSegment{}).Find(&tmpUserSegments)

	userSegments := make([]Data, 0)
	db.DB.Table("user_segments").Find(&userSegments)
	if r := db.DB.Model(&models.History{}).Where("(EXTRACT('Year' FROM timestamp) = ? AND EXTRACT('Month' FROM timestamp) = ?)",
		yearMonth.Year,
		yearMonth.Month).
		Find(&histories); r.RowsAffected == 0 {
		c.String(http.StatusOK, "no data in this period")
	} else {
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
		for _, history := range histories {
			row := []string{
				strconv.Itoa(int(history.ID)),
				strconv.FormatUint(uint64(history.UserID), 10),
				strconv.FormatUint(uint64(history.SegmentID), 10),
				history.SegmentName,
				history.Operation,
				history.Timestamp.Format(time.RFC3339),
			}
			if err := writer.Write(row); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"report_url": filePath})
	}
}

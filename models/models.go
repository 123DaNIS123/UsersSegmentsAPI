package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Segments []Segment `json:"-" gorm:"many2many:user_segments;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Segment struct {
	ID         uint   `gorm:"primaryKey;"`
	Name       string `json:"name" gorm:"uniqueIndex"`
	Percentage uint   `json:"percentage"`
	// UserSegment UserSegment `gorm:"Foreignkey:SegmentID;"`
}

type UserSegment struct {
	UserID    uint `gorm:"primaryKey;column:user_id"`
	SegmentID uint `gorm:"primaryKey;column:segment_id"`
	TTL       time.Duration
	CreatedAt time.Time `json:"created_at"`
}

type History struct {
	ID          uint      `json:"id" gorm:"primarykey;column:id"`
	UserID      uint      `json:"user_id" gorm:"column:user_id"`
	SegmentID   uint      `json:"segment_id" gorm:"column:segment_id"`
	SegmentName string    `json:"segment_name" gorm:"column:segment_name"`
	Operation   string    `json:"operation" gorm:"column:operation"`
	Timestamp   time.Time `json:"timestamp" gorm:"primarykey;column:timestamp"`
}

func (v UserSegment) AfterFind(tx gorm.DB) (err error) {
	currentTime := time.Now()
	fmt.Printf("CreatedAt: %v; ExpireTime: %v; CurrentTime: %v", v.CreatedAt, v.CreatedAt.Add(v.TTL), time.Now())
	if currentTime.After(v.CreatedAt.Add(v.TTL)) && v.TTL != 0 {
		err = tx.Delete(v).Error
	}
	// fmt.Println(v)
	return err
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint        `gorm:"primaryKey"`
	Segments    []Segment   `json:"-" gorm:"many2many:user_segments;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Segment struct {
	ID          uint        `gorm:"primaryKey;"`
	Name        string      `json:"name" gorm:"uniqueIndex"`
	// UserSegment UserSegment `gorm:"Foreignkey:SegmentID;"`
}

type UserSegment struct {
	UserID    uint `gorm:"primaryKey"`
	SegmentID uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// func (usersegment *UserSegment) BeforeCreate(tx *gorm.DB) (err error) {
// 	// if usersegment.UserID.exist
// 	return
// }

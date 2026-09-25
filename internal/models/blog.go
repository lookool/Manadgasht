package models

import "time"

// BlogPost یک مقاله‌ی وبلاگ سایت است.
type BlogPost struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"size:200;not null"`
	Slug       string `gorm:"size:220;uniqueIndex;not null"`
	Summary    string `gorm:"size:300"`
	Content    string `gorm:"type:text"`
	CoverImage string `gorm:"size:255"`
	IsPublished bool  `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

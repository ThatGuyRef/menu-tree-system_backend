package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Menu represents a hierarchical menu item
type Menu struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index" json:"parentId"`
	Depth     int        `gorm:"not null;default:0" json:"depth"`
	SortOrder int        `gorm:"not null;default:0" json:"sortOrder"`
	Children  []Menu     `gorm:"-" json:"children"`
	Parent    *Menu      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// BeforeCreate generates a UUID before inserting
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// --- Request DTOs ---

// CreateMenuInput is the payload for creating a menu item
type CreateMenuInput struct {
	Name     string     `json:"name" binding:"required"`
	ParentID *uuid.UUID `json:"parentId"`
}

// UpdateMenuInput is the payload for updating a menu item
type UpdateMenuInput struct {
	Name string `json:"name" binding:"required"`
}

// MoveMenuInput is the payload for moving a menu to a different parent
type MoveMenuInput struct {
	ParentID *uuid.UUID `json:"parentId"`
}

// ReorderMenuInput is the payload for reordering a menu within same level
type ReorderMenuInput struct {
	NewPosition int `json:"newPosition" binding:"required"`
}
package entities

import "github.com/google/uuid"

type BaseItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text;not null;default:''"`
	Category    string    `gorm:"type:text;not null;default:''"`
}

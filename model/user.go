package model

import (
	"github.com/jinzhu/gorm"
)

type UserRole int

const (
	RoleUser UserRole = iota
	RoleAdmin
)

type User struct {
	gorm.Model
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email" gorm:"index"`
	Password string   `json:"password,omitempty"`
	Role     UserRole `json:"role" gorm:"default:0"`
}

package model

import "gorm.io/gorm"

// Card model
type Card struct {
	gorm.Model
	Year   uint
	Set    string
	Maker  string
	Player string
}

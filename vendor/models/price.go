package models

import (
  "github.com/jinzhu/gorm"
)

type PriceModel struct {
	gorm.Model
	Price 		float64
}

func (p *PriceModel) SavePrice(db *gorm.DB) {
	db.Create(&p)
}
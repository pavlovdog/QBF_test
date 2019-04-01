package models

import (
    "github.com/go-pg/pg"
	// "time"
)

type PriceModel struct {
	Price 		float64
	Timestamp	int64
}

func (p *PriceModel) SavePrice(db *pg.DB) error {
	// p.Timestamp = int64(time.Now().Unix())
	return db.Insert(&p)
}
package models

import (
	"errors"
	//"html"
	//"strings"
	//"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type OfferingItem struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Price    uint64    `gorm:"not null;unique" json:"price"`
	StoreID  int  		`gorm:"not null;unique" json:"store_id"`
	Name     string    `gorm:"size:255;not null;unique" json:"name"`
	StoreName     string    `gorm:"size:255;not null;unique" json:"store_name"`
	Image     string    `gorm:"size:255;not null;unique" json:"image"`
	Description  string    `gorm:"size:255;not null;unique" json:"description"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *OfferingItem) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()

}

func (p *OfferingItem) Validate() error {

	if p.StoreID == 0 {
		return errors.New("Required UserID")
	}
	if p.Price == 0 {
		return errors.New("Required Số Lượng")
	}
	if p.Name =="" {
		return errors.New("Required Author")
	}
	return nil
}

func (p *OfferingItem) SaveOfferingItem(db *gorm.DB) (*OfferingItem, error) {

	
	var err error
	err = db.Debug().Model(&OfferingItem{}).Create(&p).Error
	if err != nil {
		return &OfferingItem{}, err
	}
	
	return p, nil
}



func (p *OfferingItem) FindOfferingItemByID(db *gorm.DB, pid uint64) ( *OfferingItem, error) {
	var err error
	err = db.Debug().Model(&OfferingItem{}).Where("id = ? and status = 1", pid).Take(&p).Error
	if err != nil {
		return &OfferingItem{}, err
	}
	
	return p, nil
	
}


package models

import (
	"errors"
	//	"html"
	//	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Tb Order
type OrderDetail struct {
	ID              uint64        `gorm:"primary_key;auto_increment" json:"id"`
	OrderID          uint64        `sql:"type:int64 REFERENCES users(id)" json:"order_id"`
	Amount          uint64        `sql:"type:int64" json:"amount"`
	ItemID       int8          `gorm:"size:11;not null;" json:"item_id"`
	Quantity     int8          `gorm:"size:11;not null;" json:"quantity"`
	CreatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status          int8          `gorm:"size:11;not null;" json:"status"`
}


// Person is a representation of a person
func (p *OrderDetail) Prepare() {
	p.ID = 0

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *OrderDetail) Validate() error {

	if p.OrderID == 0 {
		return errors.New("Required OrderID")
	}
	if p.Amount == 0 {
		return errors.New("Required Amount")
	}
	
	return nil
}

//Save
func (p *OrderDetail) SaveOrder(db *gorm.DB) (*OrderDetail, error) {
	var err error
	err = db.Debug().Model(&OrderDetail{}).Create(&p).Error
	if err != nil {
		return &OrderDetail{}, err
	}
	return p, nil
}

//FindALL
func (p *OrderDetail) FindAllOrders(db *gorm.DB) (*[]Order, error) {
	var err error
	rows := []Order{}
	err = db.Debug().Model(&Order{}).Limit(100).Find(&rows).Error
	if err != nil {
		return &[]Order{}, err
	}
	
	return &rows, nil
}



func (p *OrderDetail) UpdateAOrder(db *gorm.DB) (*OrderDetail, error) {

	var err error
	err = db.Debug().Model(&OrderDetail{}).Where("item_id = ? and order_id = ?", p.ItemID, p.OrderID).Updates(OrderDetail{Status: p.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &OrderDetail{}, err
	}
	return p, nil
}

func (p *OrderDetail) DeleteAOrder(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Order{}).Where("id = ?", id).Take(&Order{}).Delete(&Order{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

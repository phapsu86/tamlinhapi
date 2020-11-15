package models

import (
	"errors"
	//	"html"
	//	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Tb Order
type Order struct {
	ID              uint64        `gorm:"primary_key;auto_increment" json:"id"`
	UserID          uint64        `sql:"type:int64 REFERENCES users(id)" json:"user_id"`
	Amount          uint64        `sql:"type:int64" json:"amount"`
	OrderType       int8          `gorm:"size:11;not null;" json:"type"`
	Notes     		string          `gorm:"size:11;not null;" json:"notes"`
	CreatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status          int8          `gorm:"size:11;not null;" json:"status"`
}


// Person is a representation of a person
func (p *Order) Prepare() {
	p.ID = 0

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Order) Validate() error {

	if p.UserID == 0 {
		return errors.New("Required UserID")
	}
	if p.Amount == 0 {
		return errors.New("Required Amount")
	}
	
	return nil
}

//Save
func (p *Order) SaveOrder(db *gorm.DB) (*Order, error) {
	var err error
	err = db.Debug().Model(&Order{}).Create(&p).Error
	if err != nil {
		return &Order{}, err
	}
	return p, nil
}

//FindALL
func (p *Order) FindAllOrders(db *gorm.DB) (*[]Order, error) {
	var err error
	rows := []Order{}
	err = db.Debug().Model(&Order{}).Limit(100).Find(&rows).Error
	if err != nil {
		return &[]Order{}, err
	}
	
	return &rows, nil
}



func (p *Order) UpdateAOrder(db *gorm.DB) (*Order, error) {

	var err error
	err = db.Debug().Model(&Order{}).Where("id = ? and user_id = ?", p.ID, p.UserID).Updates(Order{Status: p.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Order{}, err
	}
	return p, nil
}

func (p *Order) DeleteAOrder(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Order{}).Where("id = ?", id).Take(&Order{}).Delete(&Order{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

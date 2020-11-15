package models

import (
	"errors"
	//	"html"
	//	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Tb TransactionHistory
type TransactionHistory struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	TranCode  string    `sql:"type:int64 REFERENCES users(id)" json:"tran_code"`
	OrderID   uint64    `sql:"type:int64" json:"order_id"`
	Amount    uint64    `sql:"type:int64" json:"amount"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status    int8      `gorm:"size:11;not null;" json:"status"`
}

// Person is a representation of a person
func (p *TransactionHistory) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *TransactionHistory) Validate() error {

	if p.ID == 0 {
		return errors.New("Required UserID")
	}
	if p.Amount == 0 {
		return errors.New("Required Amount")
	}

	return nil
}

//Save
func (p *TransactionHistory) SaveTransactionHistory(db *gorm.DB) (*TransactionHistory, error) {
	var err error
	err = db.Debug().Model(&TransactionHistory{}).Create(&p).Error
	if err != nil {
		return &TransactionHistory{}, err
	}
	return p, nil
}

//FindALL
func (p *TransactionHistory) FindAllTransactionHistorys(db *gorm.DB) (*[]TransactionHistory, error) {
	var err error
	rows := []TransactionHistory{}
	err = db.Debug().Model(&TransactionHistory{}).Limit(100).Find(&rows).Error
	if err != nil {
		return &[]TransactionHistory{}, err
	}

	return &rows, nil
}

func (p *TransactionHistory) UpdateATransactionHistory(db *gorm.DB) (*TransactionHistory, error) {

	var err error
	err = db.Debug().Model(&TransactionHistory{}).Where("id = ? and tran_code = ?", p.ID, p.TranCode).Updates(TransactionHistory{Status: p.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &TransactionHistory{}, err
	}
	return p, nil
}

func (p *TransactionHistory) DeleteATransactionHistory(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&TransactionHistory{}).Where("id = ?", id).Take(&TransactionHistory{}).Delete(&TransactionHistory{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

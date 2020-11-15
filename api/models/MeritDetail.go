package models

import (
	"errors"
	//	"html"
	//	"strings"
	//"time"
	"github.com/jinzhu/gorm"
)

//Tb MeritDetail
type MeritDetail struct {
	ID             uint64       `gorm:"primary_key;auto_increment" json:"id"`
	Amount         uint64       `gorm:"not null;" json:"amount"`
	Price          uint64       `gorm:"size:11;not null;" json:"price"`
	OfferingID     uint64       `gorm:"size:11;not null;" json:"offering_id"`
	OfferingDetail OfferingItem `json:"item_detail"`
	MeritID        uint64       `gorm:"size:11;not null;" json:"merit_id"`
	StoreID        int          `gorm:"size:11;not null;" json:"store_id"`
	Status         int8         `gorm:"size:11;not null;" json:"status"`
}

// Person is a representation of a person
func (p *MeritDetail) Prepare() {
	p.ID = 0

}

func (p *MeritDetail) Validate() error {

	if p.Amount == 0 {
		return errors.New("Required Amount")
	}
	if p.MeritID == 0 {
		return errors.New("Required MeritID")
	}

	return nil
}

//Save
func (p *MeritDetail) SaveMeritDetail(db *gorm.DB) (*MeritDetail, error) {
	var err error
	err = db.Debug().Model(&MeritDetail{}).Create(&p).Error
	if err != nil {
		return &MeritDetail{}, err
	}
	return p, nil
}

//FindALL
func (p *MeritDetail) FindAllMeritDetails(db *gorm.DB) (*[]MeritDetail, error) {
	var err error
	rows := []MeritDetail{}
	err = db.Debug().Model(&MeritDetail{}).Limit(100).Find(&rows).Error
	if err != nil {
		return &[]MeritDetail{}, err
	}
	// if len(rows) > 0 {
	// 	for i, _ := range rows {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", MeritDetail[i].AuthorID).Take(&MeritDetail[i].Author).Error
	// 		if err != nil {
	// 			return &[]MeritDetail{}, err
	// 		}
	// 	}
	// }

	return &rows, nil
}

func (p *MeritDetail) FindMeritDetailByID(db *gorm.DB, mid uint64) (*[]MeritDetail, error) {
	var err error
	rows := []MeritDetail{}
	err = db.Debug().Model(&MeritDetail{}).Where("merit_id = ?", mid).Take(&rows).Error
	if err != nil {
		return &[]MeritDetail{}, err
	}
	return &rows, nil
}

func (p *MeritDetail) UpdateAMeritDetail(db *gorm.DB) (*MeritDetail, error) {

	var err error

	err = db.Debug().Model(&MeritDetail{}).Where("id = ?", p.ID).Updates(MeritDetail{Status: p.Status}).Error
	if err != nil {
		return &MeritDetail{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &MeritDetail{}, err
	// 	}
	// }
	return p, nil
}

func (p *MeritDetail) DeleteAMeritDetail(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&MeritDetail{}).Where("id = ?", id).Take(&MeritDetail{}).Delete(&MeritDetail{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

package models

import (
	"errors"
	//	"html"
	//	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Tb MeritAddressList
type MeritAddressList struct {
	ID              uint64    `gorm:"primary_key;auto_increment" json:"id"`

	Meritter        User      `json:"Meritter"`
	UserID          uint64    `sql:"type:int64 REFERENCES users(id)" json:"user_id"`
	Amount          uint64    `sql:"type:int64" json:"amount"`
	Notes           string    `gorm:"size:255;not null;" json:"notes"`
	TransactionCode string    `gorm:"size:255;not null;" json:"transaction_code"`
	ProvinceID      string       `gorm:"size:11;not null;" json:"province_id"`
	DistrictID      string      `gorm:"size:11;not null;" json:"district_id"`
	WardsID         string      `gorm:"size:11;not null;" json:"wards_id"`
	Address         string    `gorm:"size:255;not null;" json:"address"`
	IsAnonymous     int8      `gorm:"size:11;not null;" json:"is_anonymous"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status          int8      `gorm:"size:11;not null;" json:"status"`
}

// Person is a representation of a person
func (p *MeritAddressList) Prepare() {
	p.ID = 0
	p.Meritter = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *MeritAddressList) Validate() error {

	if p.UserID == 0 {
		return errors.New("Required UserID")
	}
	if p.Amount == 0 {
		return errors.New("Required Amount")
	}
	if p.Address == "" {
		return errors.New("Required Author")
	}
	return nil
}

//Save
func (p *MeritAddressList) SaveMeritAddressList(db *gorm.DB) (*MeritAddressList, error) {
	var err error
	err = db.Debug().Model(&MeritAddressList{}).Create(&p).Error
	if err != nil {
		return &MeritAddressList{}, err
	}
	return p, nil
}

//FindALL
func (p *MeritAddressList) FindAllMeritAddressLists(db *gorm.DB) (*[]MeritAddressList, error) {
	var err error
	rows := []MeritAddressList{}
	err = db.Debug().Model(&MeritAddressList{}).Limit(100).Find(&rows).Error
	if err != nil {
		return &[]MeritAddressList{}, err
	}

	return &rows, nil
}

func (p *MeritAddressList) FindMeritAddressListByObjectID(db *gorm.DB, obj_id uint64, obj_type int, page int) (*[]MeritAddressList, error) {
	var err error
	rows := []MeritAddressList{}
	err = db.Debug().Model(&MeritAddressList{}).Where("object_id = ? and object_type = ?", obj_id, obj_type).Limit(10).Offset(10 * page).Find(&rows).Error
	if err != nil {
		return &[]MeritAddressList{}, err
	}

	if len(rows) > 0 {
		for i, _ := range rows {
			err := db.Debug().Model(&User{}).Where("id = ?", rows[i].UserID).Take(&rows[i].Meritter).Error
			if err != nil {
				return &[]MeritAddressList{}, err
			}
		}
	}

	return &rows, nil
}

// =================Lấy danh sách lễ vật của user======================
func (p *MeritAddressList) GetMeritInfoByStatus(db *gorm.DB, user_id uint64, status int, merit_type int, page int) (*[]MeritAddressList, error) {
	var err error
	rows := []MeritAddressList{}
	err = db.Debug().Model(&MeritAddressList{}).Where("user_id = ? and status = ? and merit_type = ?", user_id, status, merit_type).Limit(10).Offset(10 * page).Find(&rows).Error
	if err != nil {
		return &[]MeritAddressList{}, err
	}

	// if merit_type == 0 {
	// 	return &rows, nil
	// }

	// if len(rows) > 0 {
	// 	for i, _ := range rows {
	// 		details := []MeritDetail{}
	// 		err := db.Debug().Model(&MeritDetail{}).Where("merit_id = ?", rows[i].ID).Find(&details).Error
	// 		if err != nil {
	// 			return &[]MeritAddressList{}, err
	// 		}

	// 		if len(details) > 0 {
	// 			for j, _ := range details {
	// 				err := db.Debug().Model(&OfferingItem{}).Where("id = ?", details[i].OfferingID).Take(&details[j].OfferingDetail).Error
	// 				if err != nil {
	// 					return &[]MeritAddressList{}, err
	// 				}

	// 			}

	// 		}

	// 		rows[i].MeritDetail = details

	// 	}
	// }

	return &rows, nil
}

//========================================

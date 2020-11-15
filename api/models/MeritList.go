package models

import (
	"errors"
	//	"html"
	//	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Tb MeritList
type MeritList struct {
	ID              uint64        `gorm:"primary_key;auto_increment" json:"id"`
	StoreID         int           `gorm:"size:11;not null;" json:"store_id"`
	Meritter        Meritter      `json:"meritter"`
	UserID          uint64        `sql:"type:int64 REFERENCES users(id)" json:"user_id"`
	Amount          uint64        `sql:"type:int64" json:"amount"`
	ObjectName      string        `gorm:"size:255;not null;" json:"object_name"`
	Notes           string        `gorm:"size:255;not null;" json:"notes"`
	TextToStore     string        `gorm:"size:255;not null;" json:"text_to_store"`
	MeritDetail     []MeritDetail `json:"meritdetails"`
	TransactionCode string        `gorm:"size:255;not null;" json:"transaction_code"`
	ObjectID        int           `gorm:"size:11;not null;" json:"object_id"`
	ObjectType      int8          `gorm:"size:11;not null;" json:"object_type"`
	MeritType       int8          `gorm:"size:11;not null;" json:"type"`
	IsAnonymous     int8          `gorm:"size:11;not null;" json:"is_anonymous"`
	CreatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status          int8          `gorm:"size:11;not null;" json:"status"`
}

type Meritter struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

type Response struct {
	User  Meritter `json:"user_info"`
	Total uint64   `json:"total"`
}

// Person is a representation of a person
func (p *MeritList) Prepare() {
	p.ID = 0
	//p.Meritter = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *MeritList) Validate() error {

	if p.UserID == 0 {
		return errors.New("Required UserID")
	}
	if p.Amount == 0 {
		return errors.New("Required Amount")
	}
	if p.ObjectID == 0 {
		return errors.New("Required Author")
	}
	return nil
}

//Save
func (p *MeritList) SaveMeritList(db *gorm.DB) (*MeritList, error) {
	var err error
	err = db.Debug().Model(&MeritList{}).Create(&p).Error
	if err != nil {
		return &MeritList{}, err
	}
	return p, nil
}

//FindALL
func (p *MeritList) FindAllMeritLists(db *gorm.DB) (*[]MeritList, error) {
	var err error
	rows := []MeritList{}
	err = db.Debug().Model(&MeritList{}).Limit(100).Find(&rows).Error
	if err != nil {
		return &[]MeritList{}, err
	}
	// if len(rows) > 0 {
	// 	for i, _ := range rows {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", MeritList[i].AuthorID).Take(&MeritList[i].Author).Error
	// 		if err != nil {
	// 			return &[]MeritList{}, err
	// 		}
	// 	}
	// }

	return &rows, nil
}

func (p *MeritList) FindMeritListByObjectID(db *gorm.DB, obj_id uint64, obj_type int, page int) (*[]Response, error) {
	var err error
	rows := []MeritList{}
	err = db.Debug().Model(&MeritList{}).Select("user_id, sum(amount) as amount").Where("object_id = ? and object_type = ?", obj_id, obj_type).Group("user_id").Limit(10).Offset(10 * page).Find(&rows).Error
	if err != nil {
		return &[]Response{}, err
	}
	rs := []Response{}
	if len(rows) > 0 {
		for i, _ := range rows {
			u := User{}
			err := db.Debug().Model(&User{}).Where("id = ?", rows[i].UserID).Take(&u).Error
			if err != nil {
				return &[]Response{}, err
			}
			rows[i].Meritter = Meritter{Name: u.Nickname, Mobile: u.Mobile}
			item := Response{User: rows[i].Meritter, Total: rows[i].Amount}
			rs = append(rs, item)
		}
	}

	return &rs, nil
}

// =================Lấy danh sách lễ vật của user======================
func (p *MeritList) GetMeritInfoByStatus(db *gorm.DB, user_id uint64, status int, merit_type int, page int) ([]MeritList, error) {
	var err error
	rows := []MeritList{}
	err = db.Debug().Model(&MeritList{}).Where("user_id = ? and status = ? and merit_type = ?", user_id, status, merit_type).Limit(10).Offset(10 * page).Find(&rows).Error
	if err != nil {
		return []MeritList{}, err
	}

	if merit_type == 0 {
		return rows, nil
	}

	if len(rows) > 0 {
		for i, _ := range rows {
			details := []MeritDetail{}
			err := db.Debug().Model(&MeritDetail{}).Where("merit_id = ?", rows[i].ID).Find(&details).Error
			if err != nil {
				return []MeritList{}, err
			}

			if len(details) > 0 {
				for j, _ := range details {
					err := db.Debug().Model(&OfferingItem{}).Where("id = ?", details[j].OfferingID).Take(&details[j].OfferingDetail).Error
					if err != nil {
						return []MeritList{}, err
					}

				}

			}

			rows[i].MeritDetail = details

		}
	}

	return rows, nil
}

//========================================

func (p *MeritList) UpdateAMeritList(db *gorm.DB) (*MeritList, error) {

	var err error

	err = db.Debug().Model(&MeritList{}).Where("object_id = ? and user_id = ?", p.ObjectID, p.UserID).Updates(MeritList{Status: p.Status, Notes: p.Notes, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &MeritList{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &MeritList{}, err
	// 	}
	// }
	return p, nil
}

func (p *MeritList) DeleteAMeritList(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&MeritList{}).Where("id = ?", id).Take(&MeritList{}).Delete(&MeritList{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

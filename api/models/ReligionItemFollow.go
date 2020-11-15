package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//Tb UserRemind
type ReligionItemFollow struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"size:255;not null;unique" json:"user_id"`
	ItemID    uint64    `gorm:"size:255;not null;" json:"item_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status    int       `gorm:"size:255;not null;" json:"status"`
}

// Person is a representation of a person
func (p *ReligionItemFollow) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *ReligionItemFollow) Validate() error {

	return nil
}

//Save
func (p *ReligionItemFollow) SaveUserItemFollow(db *gorm.DB, uid uint64, itemID uint64, actionType int) error {
	var err error
	if actionType == 1 {
		err = db.Debug().Model(&ReligionItemFollow{}).Create(&p).Error
		if err != nil {
			return err
		}

	} else if actionType == 0 {
		db = db.Debug().Model(&ReligionItemFollow{}).Where("item_id = ? and user_id = ?", itemID, uid).Take(&ReligionItemFollow{}).Delete(&ReligionItemFollow{})
		if db.Error != nil {
			if gorm.IsRecordNotFoundError(db.Error) {
				return errors.New("NOT_FOUND")
			}
			return db.Error
		}

	}

	return nil
}

//FindALL
func (p *ReligionItemFollow) FindAllReligionItemFollow(db *gorm.DB, uid uint64) (*[]ReligionItemFollow, error) {
	var err error
	rows := []ReligionItemFollow{}
	err = db.Debug().Model(&ReligionItemFollow{}).Where("user_id = ?", uid).Limit(100).Find(&rows).Error
	if err != nil {
		return &[]ReligionItemFollow{}, err
	}
	// if len(rows) > 0 {
	// 	for i, _ := range rows {
	// 		err := db.Debug().Model(&ReligionEvent{}).Where("id = ?", rows[i].ID).Take(&rows[i].EventDetail).Error
	// 		if err != nil {
	// 			return &[]ReligionItemFollow{}, err
	// 		}
	// 	}
	// }

	return &rows, nil
}



func (p *ReligionItemFollow) GetReligionItemFollowByID(db *gorm.DB, uid uint64,item_id uint64) (*ReligionItemFollow, error) {
	var err error
	item := ReligionItemFollow{}
	err = db.Debug().Model(&ReligionItemFollow{}).Where("user_id = ? and item_id = ?", uid,item_id).Take(&item).Error
	if err != nil {
		return &ReligionItemFollow{}, err
	}
	
	return &item, nil
}



func (p *ReligionItemFollow) FindUserReligionItemFollowID(db *gorm.DB, uid uint64, itemid uint64) (*ReligionItemFollow, error) {
	var err error
	rows := ReligionItemFollow{}
	err = db.Debug().Model(&ReligionItemFollow{}).Where("item_id = ? and user_id = ?", itemid, uid).Take(&rows).Error
	if err != nil {
		return &ReligionItemFollow{}, err
	}

	return &rows, nil

}

func (p *ReligionItemFollow) DeleteAReligionItemFollow(db *gorm.DB, id uint64, uid uint64) (int64, error) {

	db = db.Debug().Model(&UserRemind{}).Where("object_id = ? and user_id = ?", id, uid).Take(&UserRemind{}).Delete(&UserRemind{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

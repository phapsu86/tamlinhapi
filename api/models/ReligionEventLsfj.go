package models

import (
	"errors"
	//	"html"
	//	"strings"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// type result struct {
// 	TotalJoin int
// 	TotalShare int
// 	TotalFollow int
//   }

//Tb ReligionEventLsfj
type ReligionEventLsfj struct {
	ID           uint64        `gorm:"primary_key;auto_increment" json:"id"`
	UserID       uint64        `gorm:"size:11;not null;unique" json:"user_id"`
	EventID      uint64           `gorm:"size:11;not null;" json:"event_id"`
	EventDetails ReligionEvent `json:"event_details"`
	IsJoin       *int           `gorm:"size:1;not null;" json:"is_join"`
	IsShare      *int           `gorm:"size:1;not null;" json:"is_share"`
	IsFollow     *int           `gorm:"size:1;not null;" json:"is_follow"`
	CreatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Person is a representation of a person
func (p *ReligionEventLsfj) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *ReligionEventLsfj) Validate() error {

	if p.EventID == 0 {
		return errors.New("Required Event_ID")
	}

	return nil
}

//Save
func (p *ReligionEventLsfj) SaveReligionEventLsfj(db *gorm.DB,uid uint64,eventID uint64, value int, actionType int) (error) {
	var err error

	item := ReligionEventLsfj{}

	err = db.Debug().Model(&ReligionEventLsfj{}).Where("event_id = ? and user_id = ?", eventID, uid).First(&item).Error
	if item.UserID != 0 {
	
		if actionType == 0 { //Join
			err = db.Debug().Model(&ReligionEventLsfj{}).Where("event_id = ? and user_id = ?", eventID, uid).Updates(ReligionEventLsfj{IsJoin: &value, UpdatedAt: time.Now()}).Error
		} else if actionType == 1 { //share
			err = db.Debug().Model(&ReligionEventLsfj{}).Where("event_id = ? and user_id = ?", eventID, uid).Updates(ReligionEventLsfj{IsShare: &value, UpdatedAt: time.Now()}).Error
		} else if actionType == 2 { //Follow
			err = db.Debug().Model(&ReligionEventLsfj{}).Where("event_id = ? and user_id = ?", eventID, uid).Updates(ReligionEventLsfj{IsFollow: &value, UpdatedAt: time.Now()}).Error
		}
		if err != nil {
			return err
		}

		} else {
			item.Prepare()
			item.UserID = uid
			item.EventID = eventID
			
			
			if actionType == 0 { //Join
				item.IsJoin = &value
				err = db.Debug().Model(&ReligionEventLsfj{}).Create(&item).Error
				
			} else if actionType == 1 { //share
				item.IsShare = &value
				err = db.Debug().Model(&ReligionEventLsfj{}).Create(&item).Error
			} else if actionType == 2 { //Follow
				item.IsFollow = &value
				err = db.Debug().Model(&ReligionEventLsfj{}).Create(&item).Error
			}

			if err != nil {
				return err
			}
			

		}
		return nil
	}



//FindALL
func (p *ReligionEventLsfj) GetTotalEventLsfjs(db *gorm.DB, event_id uint64) (int64, int64, int64, error) {
	var err error

	rows, err := db.Debug().Model(&ReligionEventLsfj{}).Select("sum(is_join) as total_join,sum(is_share) as total_share,sum(is_follow) as total_follow").Where("event_id = ?", event_id).Group("event_id").Rows()
	if err != nil {
		return 0, 0, 0, err
	}
	var total_join int64
	var total_share int64
	var total_follow int64
	for rows.Next() {

		rows.Scan(&total_join, &total_share, &total_follow)
		fmt.Println(total_join, &total_share, total_follow)

	}

	return total_join, total_share, total_follow, nil
}

func (p *ReligionEventLsfj) FindEventLsfjFollow(db *gorm.DB, uid uint64, page uint64) ([]ReligionEventLsfj, error) {
	var err error
	items := []ReligionEventLsfj{}

	err = db.Debug().Model(&Post{}).Where("user_id = ?", uid).Limit(20).Offset(page * 20).Find(&items).Error
	if err != nil {
		return []ReligionEventLsfj{}, err
	}

	if len(items) > 0 {
		for i, _ := range items {
			err := db.Debug().Model(&ReligionEvent{}).Where("id = ?", items[i].EventID).Take(&items[i].EventDetails).Error
			if err != nil {
				return []ReligionEventLsfj{}, err
			}
		}
	}

	return items, nil
}

func (p *ReligionEventLsfj) FindReligionEventLsfjByID(db *gorm.DB, uid uint64, event_id uint64) (*ReligionEventLsfj, error) {
	var err error
	err = db.Debug().Model(&ReligionEventLsfj{}).Where("user_id = ? and event_id = ?", uid, event_id).Take(&p).Error
	if err != nil {
		return &ReligionEventLsfj{}, err
	}

	return p, nil
}

func (p *ReligionEventLsfj) UpdateAReligionEventLsfj(db *gorm.DB) (*ReligionEventLsfj, error) {

	var err error

	err = db.Debug().Model(&ReligionEventLsfj{}).Where("id = ?", p.ID).Updates(ReligionEventLsfj{IsFollow: p.IsFollow, IsShare: p.IsShare, IsJoin: p.IsJoin, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &ReligionEventLsfj{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &ReligionEventLsfj{}, err
	// 	}
	// }
	return p, nil
}

func (p *ReligionEventLsfj) DeleteAReligionEventLsfj(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&ReligionEventLsfj{}).Where("id = ?", id).Take(&ReligionEventLsfj{}).Delete(&ReligionEventLsfj{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

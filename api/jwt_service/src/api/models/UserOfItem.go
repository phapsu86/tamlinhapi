package models

import (
	"errors"
	
	"time"
	"github.com/jinzhu/gorm"
	
)
//Tb UserOfItem
type UserOfItem struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname     string    `gorm:"size:255;not null;unique" json:"nick_name"`
	RoomID     string    `gorm:"size:255;not null;unique" json:"room_id"`
	LinkToServer     string    `gorm:"size:255;not null;unique" json:"link_to_server"`
	UserID   uint64    `gorm:"size:255;not null;" json:"user_id"`
	ItemID    int     `gorm:"size:255;not null;" json:"item_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
// Person is a representation of a person
func (p *UserOfItem) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *UserOfItem) Validate() error {

	// if p.Name == "" {
	// 	return errors.New("Required Title")
	// }
	// if p.Image == "" {
	// 	return errors.New("Required Content")
	// }
	// if p.Notes == "" {
	// 	return errors.New("Required Author")
	// }
	return nil
}
//Save
func (p *UserOfItem) SaveUserOfItem(db *gorm.DB) (*UserOfItem, error) {
	var err error
	err = db.Debug().Model(&UserOfItem{}).Create(&p).Error
	if err != nil {
		return &UserOfItem{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.ID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &Sp_religion_list{}, err
	// 	}
	// }
	return p, nil
}

 

 

//FindALL
func (p *UserOfItem) FindAllUserOfItems(db *gorm.DB) ([]UserOfItem, error) {
	var err error
	rows := []UserOfItem{}
	err = db.Debug().Model(&UserOfItem{}).Limit(10).Find(&rows).Error
	if err != nil {
		return []UserOfItem{}, err
	}
	// if len(rows) > 0 {
	// 	for i, _ := range rows {
			
	// 		err := db.Debug().Model(&User{}).Where("id = ?", UserOfItem[i].AuthorID).Take(&UserOfItem[i].Author).Error
	// 		if err != nil {
	// 			return &[]UserOfItem{}, err
	// 		}
	// 	}
	// }
 


	return rows, nil
}

func (p *UserOfItem) FindUserOfItemByID(db *gorm.DB, pid uint64,item_id uint64) (*UserOfItem, error) {
	var err error
	err = db.Debug().Model(&UserOfItem{}).Where("user_id = ? and item_id = ? and status = 1", pid,item_id).Take(&p).Error
	if err != nil {
		return &UserOfItem{}, err
	}
	return p, nil
}



func (p *UserOfItem) DeleteAUserOfItem(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&UserOfItem{}).Where("id = ?", id).Take(&UserOfItem{}).Delete(&UserOfItem{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

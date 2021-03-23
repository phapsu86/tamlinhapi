package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/phapsu86/tamlinh/api/utils/strconvert"
)

//Tb ReligionItem
type ReligionItem struct {
	ID           uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Name         string `gorm:"size:255;not null;unique" json:"name"`
	Image        string `gorm:"size:255;not null;" json:"image"`
	Lat          string `gorm:"size:255;not null;" json:"lat"`
	Lon          string `gorm:"size:255;not null;" json:"lon"`
	Address      string `gorm:"size:255;not null;" json:"address"`
	Descriptions string `gorm:"not null;" json:"descriptions"`
	Phone        string `gorm:"size:255;not null;" json:"phone"`
	Code         string `gorm:"size:255;not null;" json:"code"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Person is a representation of a person
func (p *ReligionItem) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	//p.Notes = p.Notes
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *ReligionItem) Validate() error {

	if p.Name == "" {
		return errors.New("Required Title")
	}
	if p.Image == "" {
		return errors.New("Required Content")
	}
	// if p.Notes == "" {
	// 	return errors.New("Required Author")
	// }
	return nil
}

//Save
func (p *ReligionItem) SaveReligionItem(db *gorm.DB) (*ReligionItem, error) {
	var err error
	err = db.Debug().Model(&ReligionList{}).Create(&p).Error
	if err != nil {
		return &ReligionItem{}, err
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
func (p *ReligionItem) FindAllReligionItems(db *gorm.DB) ([]ReligionItem, error) {
	var err error
	rows := []ReligionItem{}
	err = db.Debug().Model(&ReligionItem{}).Limit(100).Find(&rows).Error
	if err != nil {
		return []ReligionItem{}, err
	}

	return rows, nil
}

func (p *ReligionItem) FindReligionItemDetails(db *gorm.DB, id uint64) (*ReligionItem, error) {
	var err error
	item := ReligionItem{}
	err = db.Debug().Model(&ReligionItem{}).Where("id = ?", id).First(&item).Error
	if err != nil {
		return &ReligionItem{}, err
	}

	return &item, nil
}

func (p *ReligionItem) CheckReligionItem(db *gorm.DB, id int) (*ReligionItem, error) {
	var err error
	item := ReligionItem{}
	err = db.Debug().Model(&ReligionItem{}).Where("id = ? and status = 1", id).First(&item).Error
	if err != nil {
		return &ReligionItem{}, err
	}

	return &item, nil
}

func (p *ReligionItem) FindReligionItemtByReligionID(db *gorm.DB, pid uint64, page uint64) ([]ReligionItem, error) {
	var err error
	rows := []ReligionItem{}
	err = db.Debug().Model(&ReligionItem{}).Where("religion_id = ?", pid).Limit(10).Offset(page * 10).Find(&rows).Error
	if err != nil {
		return []ReligionItem{}, err
	}

	return rows, nil
}

func (p *ReligionItem) SearchItemByName(db *gorm.DB, ReligionID int, Keyword string, locID string, page int) ([]ReligionItem, error) {
	var err error
	rows := []ReligionItem{}
	txtSearch := strconvert.ConvertVitoEn(Keyword)
	fmt.Printf("xxxxxxxx%v", txtSearch)
	if Keyword == "" && locID != "" {
		err = db.Debug().Model(&ReligionItem{}).Where("religion_id = ? and status = ? and province_id = ? ", ReligionID, 1, locID).Limit(20).Offset(20 * page).Find(&rows).Error

	} else if locID == "" && Keyword != "" {
		err = db.Debug().Model(&ReligionItem{}).Where("religion_id = ? and keyword like ? and status = ? ", ReligionID, "%"+txtSearch+"%", 1).Limit(20).Offset(20 * page).Find(&rows).Error

	} else if locID != "" && Keyword != "" {
		err = db.Debug().Model(&ReligionItem{}).Where("religion_id = ? and keyword like ? and status = ? and province_id = ? ", ReligionID, "%"+txtSearch+"%", 1, locID).Limit(20).Offset(20 * page).Find(&rows).Error

	} else {

		err = db.Debug().Model(&ReligionItem{}).Where("religion_id = ? and status = ?", ReligionID, 1).Limit(20).Offset(20 * page).Find(&rows).Error

	}
	if err != nil {
		return []ReligionItem{}, err
	}

	return rows, nil
}

package models

import (
	"errors"
	"html"
	"strings"
	"time"
	"github.com/jinzhu/gorm"
	
)
//Tb Religionlist
type ReligionList struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name     string    `gorm:"size:255;not null;unique" json:"name"`
	Image   string    `gorm:"size:255;not null;" json:"image"`
	Notes    string     `gorm:"size:255;not null;" json:"notes"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
// Person is a representation of a person
func (p *ReligionList) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Notes))
	p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	p.Notes = p.Notes
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *ReligionList) Validate() error {

	if p.Name == "" {
		return errors.New("Required Title")
	}
	if p.Image == "" {
		return errors.New("Required Content")
	}
	if p.Notes == "" {
		return errors.New("Required Author")
	}
	return nil
}
//Save
func (p *ReligionList) SaveReligionList(db *gorm.DB) (*ReligionList, error) {
	var err error
	err = db.Debug().Model(&ReligionList{}).Create(&p).Error
	if err != nil {
		return &ReligionList{}, err
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
func (p *ReligionList) FindAllReligionLists(db *gorm.DB) ([]ReligionList, error) {
	var err error
	rows := []ReligionList{}
	err = db.Debug().Model(&ReligionList{}).Limit(10).Find(&rows).Error
	if err != nil {
		return []ReligionList{}, err
	}
	// if len(rows) > 0 {
	// 	for i, _ := range rows {
			
	// 		err := db.Debug().Model(&User{}).Where("id = ?", ReligionList[i].AuthorID).Take(&ReligionList[i].Author).Error
	// 		if err != nil {
	// 			return &[]ReligionList{}, err
	// 		}
	// 	}
	// }
 


	return rows, nil
}

func (p *ReligionList) FindReligionListByID(db *gorm.DB, pid uint64) (*ReligionList, error) {
	var err error
	err = db.Debug().Model(&ReligionList{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &ReligionList{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &ReligionList{}, err
	// 	}
	// }
	return p, nil
}

func (p *ReligionList) UpdateAReligionList(db *gorm.DB) (*ReligionList, error) {

	var err error
	
	err = db.Debug().Model(&ReligionList{}).Where("id = ?", p.ID).Updates(ReligionList{Name: p.Name, Image: p.Image,Notes:p.Notes, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &ReligionList{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &ReligionList{}, err
	// 	}
	// }
	return p, nil
}

func (p *ReligionList) DeleteAReligionList(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&ReligionList{}).Where("id = ?", id).Take(&ReligionList{}).Delete(&ReligionList{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

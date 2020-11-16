package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/phapsu86/tamlinhapi/api/utils/strconvert"
)

//Tb ReligionEvent
type ReligionEvent struct {
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name           string    `gorm:"size:255;not null;unique" json:"name"`
	Image          string    `gorm:"size:255;not null;" json:"image"`
	Intro          string    `gorm:"size:255;not null;" json:"intro"`
	Description    string    `gorm:"size:255;not null;" json:"description"`
	ObjectName     string    `gorm:"size:255;not null;" json:"object_name"`
	BeginDate      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"begin_date"`
	EndDate        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"end_date"`
	ReligionItemID int       `gorm:"size:255;not null;" json:"religion_item_id"`
	ReligionID     int       `gorm:"size:255;not null;" json:"religion_id"`
	IsPublic       int       `json:"is_public"`
	IsRemember     int       `json:"is_remember"`
	OrganizationBy string    `gorm:"size:255;not null;" json:"organization_by"`
	Phone          string    `gorm:"size:255;not null;" json:"phone"`
	Address        string    `gorm:"size:255;not null;" json:"address"`
	Media          string    `gorm:"size:255;not null;" json:"media"`
	KeyYoutube     string    `gorm:"size:255;not null;" json:"key_youtube"`
	IsYoutube      int       `gorm:"size:255;not null;" json:"is_youtube"`
	ShareLink      string    `gorm:"size:255;not null;" json:"share_link"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Status         string    `gorm:"size:255;not null;" json:"status"`
}

// Person is a representation of a person
func (p *ReligionEvent) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	//p.Notes = p.Notes
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *ReligionEvent) Validate() error {

	if p.Name == "" {
		return errors.New("Required Title")
	}
	if p.Image == "" {
		return errors.New("Required Content")
	}

	return nil
}

//Save
func (p *ReligionEvent) SaveReligionEvent(db *gorm.DB) (*ReligionEvent, error) {
	var err error
	err = db.Debug().Model(&ReligionList{}).Create(&p).Error
	if err != nil {
		return &ReligionEvent{}, err
	}

	return p, nil
}

//FindALL
func (p *ReligionEvent) FindAllReligionEvents(db *gorm.DB) (*[]ReligionEvent, error) {
	var err error
	rows := []ReligionEvent{}
	err = db.Debug().Model(&ReligionEvent{}).Limit(100).Offset(1).Find(&rows).Error
	if err != nil {
		return &[]ReligionEvent{}, err
	}
	// if len(rows) > 0 {
	// 	for i, _ := range rows {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", ReligionList[i].AuthorID).Take(&ReligionList[i].Author).Error
	// 		if err != nil {
	// 			return &[]ReligionList{}, err
	// 		}
	// 	}
	// }

	return &rows, nil
}

func (p *ReligionEvent) FindReligionEventByID(db *gorm.DB, pid uint64, page uint64) ([]ReligionEvent, error) {
	var err error
	rows := []ReligionEvent{}
	err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ?", pid).Limit(10).Offset(10 * page).Find(&rows).Error
	if err != nil {
		return []ReligionEvent{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &ReligionList{}, err
	// 	}
	// }
	return rows, nil
}

//Lấy danh sách sư kiện trong 1 chùa
func (p *ReligionEvent) FindReligionItemEventByID(db *gorm.DB, pid uint64, status uint64, page uint64) ([]ReligionEvent, error) {
	var err error
	rows := []ReligionEvent{}
	err = db.Debug().Model(&ReligionEvent{}).Where("religion_item_id = ? and status = ?", pid, status).Limit(10).Offset(10 * page).Find(&rows).Error
	if err != nil {
		return []ReligionEvent{}, err
	}

	return rows, nil
}

func (p *ReligionEvent) SearchEventByName(db *gorm.DB, ReligionID int, Keyword string, EventType int, locID string, page int) ([]ReligionEvent, error) {
	var err error
	rows := []ReligionEvent{}

	txtSearch := strconvert.ConvertVitoEn(Keyword)
	if locID != "" {
		if Keyword == "" {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and status = ? and province_id = ? ", ReligionID, EventType, locID).Limit(10).Offset(10 * page).Find(&rows).Error

		} else {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and keyword like ? and status = ? and province_id = ? ", ReligionID, "%"+txtSearch+"%", EventType, locID).Limit(10).Offset(10 * page).Find(&rows).Error

		}

	} else {
		if Keyword == "" {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and status = ? ", ReligionID, EventType).Limit(10).Offset(10 * page).Find(&rows).Error

		} else {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and keyword like ? and status = ?", ReligionID, "%"+txtSearch+"%", EventType).Limit(10).Offset(10 * page).Find(&rows).Error

		}

	}

	if err != nil {
		return []ReligionEvent{}, err
	}

	return rows, nil
}

func (p *ReligionEvent) SearchEventByDate(db *gorm.DB, ReligionID int, FromDate string, ToDate string, EventType int, locID string, page int) ([]ReligionEvent, error) {
	var err error
	rows := []ReligionEvent{}

	if locID != "" {
		if FromDate == "" || ToDate == "" {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and status = ? and province_id = ? ", ReligionID, EventType, locID).Limit(10).Offset(10 * page).Find(&rows).Error

		} else {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and begin_date >= ? and end_date <= ? and status = ? and province_id = ? ", ReligionID, FromDate, ToDate, EventType, locID).Limit(10).Offset(10 * page).Find(&rows).Error

		}

	} else {
		if FromDate == "" || ToDate == "" {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and status = ? ", ReligionID, EventType).Limit(10).Offset(10 * page).Find(&rows).Error

		} else {
			err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and begin_date >= ? and end_date <= ? and status = ?", ReligionID, FromDate, ToDate, EventType).Limit(10).Offset(10 * page).Find(&rows).Error

		}
	}

	if err != nil {
		return []ReligionEvent{}, err
	}

	return rows, nil
}

func (p *ReligionEvent) FindReligionEventDetail(db *gorm.DB, pid uint64) (ReligionEvent, error) {
	var err error
	item := ReligionEvent{}
	err = db.Debug().Model(&ReligionEvent{}).Where("id = ?", pid).Take(&item).Error
	if err != nil {
		return ReligionEvent{}, err
	}
	return item, nil
}

func (p *ReligionEvent) CheckEventForMerit(db *gorm.DB, pid int) (ReligionEvent, error) {
	var err error
	item := ReligionEvent{}
	err = db.Debug().Model(&ReligionEvent{}).Where("id = ? and (status = 1 or status = 0)", pid).Take(&item).Error
	if err != nil {
		return ReligionEvent{}, err
	}
	return item, nil
}

func (p *ReligionEvent) UpdateAReligionEvent(db *gorm.DB) (*ReligionEvent, error) {

	var err error

	err = db.Debug().Model(&ReligionEvent{}).Where("id = ?", p.ID).Updates(ReligionEvent{Name: p.Name, Image: p.Image, Intro: p.Intro, Description: p.Description, BeginDate: p.BeginDate, EndDate: p.EndDate, Status: p.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &ReligionEvent{}, err
	}
	// if p.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	// 	if err != nil {
	// 		return &ReligionList{}, err
	// 	}
	// }
	return p, nil
}

func (p *ReligionEvent) DeleteAReligionEvent(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&ReligionEvent{}).Where("id = ?", id).Take(&ReligionEvent{}).Delete(&ReligionEvent{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Sp_religion_list not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//Đinh nghĩa search reponse
type ResponseSearchEventForMerit struct {
	Name           string    `json:"event_name"`
	BeginDate      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"begin_date"`
	EndDate        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"end_date"`
	Image          string    `gorm:"size:255;not null;" json:"image"`
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	OrganizationBy string    `gorm:"size:255;not null;" json:"organization_by"`
}

func (p *ReligionEvent) SearchEventForMeritByName(db *gorm.DB, ReligionID int, Keyword string, page int) ([]ResponseSearchEventForMerit, error) {
	var err error
	rows := []ReligionEvent{}
	txtSearch := strconvert.ConvertVitoEn(Keyword)
	if Keyword == "" {
		err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and (status = ? or status = ?) ", ReligionID, 0, 1).Limit(50).Offset(10 * page).Find(&rows).Error

	} else {
		err = db.Debug().Model(&ReligionEvent{}).Where("religion_id = ? and keyword like ? and (status = ? or status = ?) ", ReligionID, "%"+txtSearch+"%", 0, 1).Limit(10).Offset(10 * page).Find(&rows).Error
	}
	if err != nil {
		return []ResponseSearchEventForMerit{}, err
	}

	rs := []ResponseSearchEventForMerit{}
	if len(rows) > 0 {
		for i, _ := range rows {

			//rows[i].Meritter = Meritter{ Name:u.Nickname,Mobile: u.Mobile}
			item := ResponseSearchEventForMerit{Name: rows[i].Name, ID: rows[i].ID, BeginDate: rows[i].BeginDate, EndDate: rows[i].EndDate, Image: rows[i].Image, OrganizationBy: rows[i].OrganizationBy}
			rs = append(rs, item)
		}
	}

	return rs, nil
}

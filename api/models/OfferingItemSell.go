package models

import (
	"errors"
	//"html"
	//"strings"
	//"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/phapsu86/tamlinh/api/utils/strconvert"
)

type OfferingItemSell struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	ObjectID   uint64    `gorm:"not null;unique" json:"object_id"`
	ObjectType int       `gorm:"not null;unique" json:"object_type"`
	ItemID     uint64    `gorm:"not null;unique" json:"item_id"`
	ItemDetail  OfferingItem      `json:"item_detail"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *OfferingItemSell) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()

}

func (p *OfferingItemSell) Validate() error {

	if p.ObjectID == 0 {
		return errors.New("Required UserID")
	}
	if p.ItemID == 0 {
		return errors.New("Required Số Lượng")
	}

	return nil
}

func (p *OfferingItemSell) CheckItemSellExist(db *gorm.DB, pid uint64, obj_id int, object_type int8) (*OfferingItem, error) {
	var err error
	err = db.Debug().Model(&OfferingItemSell{}).Where("item_id =? and object_id =? and object_type=? and status = 1", pid, obj_id, object_type).Take(&p).Error
	if err != nil {
		return &OfferingItem{}, err
	}

	// Lấy thông tin sản phẩm
	 var item = OfferingItem{}
	err = db.Debug().Model(&OfferingItem{}).Where("id =? and status =1", p.ItemID).Take(&item).Error
	if err != nil {
		return &OfferingItem{}, err
	}
	return &item, nil
}


func (p *OfferingItemSell) GetOfferingItemSellForObject(db *gorm.DB, obj_id uint64, obj_type int, page int) ([]OfferingItemSell, error) {
	var err error
	rows := []OfferingItemSell{}
	err = db.Debug().Model(&MeritList{}).Where("object_id = ? and object_type = ?", obj_id, obj_type).Limit(10).Offset(10 * page).Find(&rows).Error
	if err != nil {
		return []OfferingItemSell{}, err
	}

	if len(rows) > 0 {
		for i, _ := range rows {
			err := db.Debug().Model(&OfferingItem{}).Where("id = ?", rows[i].ItemID).Take(&rows[i].ItemDetail).Error
			if err != nil {
				return []OfferingItemSell{}, err
			}
		}
	}

	return rows, nil
}



func (p *OfferingItemSell) SaveOfferingItemSell(db *gorm.DB) (*OfferingItemSell, error) {

	var err error
	err = db.Debug().Model(&OfferingItemSell{}).Create(&p).Error
	if err != nil {
		return &OfferingItemSell{}, err
	}

	return p, nil
}

func (p *OfferingItemSell) FindOfferingItemSellByID(db *gorm.DB, pid uint64) (*OfferingItemSell, error) {
	var err error

	err = db.Debug().Model(&OfferingItemSell{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &OfferingItemSell{}, err
	}

	return p, nil

}

func (p *OfferingItemSell) FindOfferingItemSellByKeyword(db *gorm.DB,obj_id uint64,obj_type int8, keyword string,page int) ([]OfferingItem, error) {
	var err error
	rows := []OfferingItemSell{}
	txtSearch := strconvert.ConvertVitoEn(keyword)
//	fmt.Printf("xxxxxxxx%v", txtSearch)
if txtSearch == "" {
	err = db.Debug().Model(&OfferingItemSell{}).Where("object_id= ? and object_type =?", obj_id,obj_type).Limit(50).Offset(50*page).Find(&rows).Error

}else {
	err = db.Debug().Model(&OfferingItemSell{}).Where("keyword like ? and object_id= ? and object_type =?", "%"+txtSearch+"%",obj_id,obj_type).Limit(50).Offset(50*page).Find(&rows).Error

}

	
	if err != nil {
		return []OfferingItem{}, err
	}
	var itemArr [] uint64
	if len(rows) > 0 {
	
		for i, _ := range rows {
			itemArr = append(itemArr, rows[i].ItemID)
		}
	}

	data := []OfferingItem{}
	err = db.Debug().Model(&OfferingItem{}).Where(itemArr).Find(&data).Error
	if err != nil {
		return []OfferingItem{}, err
	}


	return data, nil

}






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

//Tb PostLsfc
type PostLsfc struct {
	
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	PostID     uint64    `gorm:"size:11;not null;unique" json:"post_id"`
	UserID   uint64    `gorm:"size:11;not null;" json:"user_id"`
	IsLike    *int     `gorm:"size:1;not null;" json:"is_like"`
	IsShare    *int     `gorm:"size:1;not null;" json:"is_share"`
	IsFollow    *int     `gorm:"size:1;not null;" json:"is_follow"`
	IsComment    *int     `gorm:"size:1;not null;" json:"is_comment"`
	Comment    string     `gorm:"size:1;not null;" json:"comment"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}




// Person is a representation of a person
func (p *PostLsfc) Prepare() {
	p.ID = 0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *PostLsfc) Validate() error {

	if p.PostID == 0 {
		return errors.New("Required POST_ID")
	}
	
	return nil
}
//Save
func (p *PostLsfc) SavePostLsfc(db *gorm.DB,uid uint64,pID uint64, value int, actionType int) (error) {
	var err error

	item := PostLsfc{}

	err = db.Debug().Model(&PostLsfc{}).Where("post_id = ? and user_id = ?", pID, uid).First(&item).Error
	if item.UserID != 0 {
	
		if actionType == 0 { //Like
			err = db.Debug().Model(&PostLsfc{}).Where("post_id = ? and user_id = ?", pID, uid).Updates(PostLsfc{IsLike: &value, UpdatedAt: time.Now()}).Error
		} else if actionType == 1 { //Comment
			err = db.Debug().Model(&PostLsfc{}).Where("post_id = ? and user_id = ?", pID, uid).Updates(PostLsfc{IsComment: &value, UpdatedAt: time.Now()}).Error
		} else if actionType == 2 { //Follow
			err = db.Debug().Model(&PostLsfc{}).Where("post_id = ? and user_id = ?", pID, uid).Updates(PostLsfc{IsFollow: &value, UpdatedAt: time.Now()}).Error
		}else if actionType == 3 { //Share
			err = db.Debug().Model(&PostLsfc{}).Where("post_id = ? and user_id = ?", pID, uid).Updates(PostLsfc{IsShare: &value, UpdatedAt: time.Now()}).Error
		}
		if err != nil {
			return err
		}

		} else {
			item.Prepare()
			item.UserID = uid
			item.PostID = pID
			
			
			if actionType == 0 { //Like
				item.IsLike = &value
				err = db.Debug().Model(&PostLsfc{}).Create(&item).Error
				
			} else if actionType == 1 { //Commnet
				item.IsComment = &value
				err = db.Debug().Model(&PostLsfc{}).Create(&item).Error
			} else if actionType == 2 { //Follow
				item.IsFollow = &value
				err = db.Debug().Model(&PostLsfc{}).Create(&item).Error
			} else if actionType == 3 { //Share
				item.IsShare = &value
				err = db.Debug().Model(&PostLsfc{}).Create(&item).Error
			}

			if err != nil {
				return err
			}
			

		}
	
// Thưc hiện update vao bảng post

tc,ts,tf,tl,e := p.GetTotalPostLsfjc(db,p.PostID) 
postItem := Post{}
posData,err := postItem.FindPostByID(db,p.PostID)
posData.NumFollow = tf
posData.NumShare = ts
posData.NumLike = tl
posData.NumComment = tc
_,e = posData.UpdateLSFPost(db)
if(e != nil){
	return  e	
}


	return  nil
}

 

 

//FindALL
func (p *PostLsfc) GetTotalPostLsfjc(db *gorm.DB, post_id uint64) (int64,int64,int64,int64, error) {
	var err error
	
	rows,err := db.Debug().Model(&PostLsfc{}).Select("sum(is_comment) as total_comment,sum(is_share) as total_share,sum(is_follow) as total_follow,sum(is_like) as total_like").Where("post_id = ?",post_id ).Group("post_id").Rows()
	if err != nil {
		return 0,0,0,0, err
	}
	var total_comment int64
	var total_share int64
	var total_follow int64
	var total_like int64
	for rows.Next() {

        rows.Scan(&total_comment, &total_share,&total_follow,&total_like)
        fmt.Println(total_comment, total_share,total_follow,total_like)
	
	  }


	return total_comment,total_share,total_follow,total_like, nil
}



func (p *PostLsfc) FindReligionPostLsfjByID(db *gorm.DB, pid uint64,event_id uint64) (*PostLsfc, error) {
	var err error
	err = db.Debug().Model(&PostLsfc{}).Where("user_id = ? and post_id = ?", pid,event_id).Take(&p).Error
	if err != nil {
		return &PostLsfc{}, err
	}
	
	return p, nil
}




func (p *PostLsfc) FindAllPostLsfjByUser(db *gorm.DB, uid uint64, page uint64) ([]PostLsfc, error) {
	var err error
	items := []PostLsfc{}

	err = db.Debug().Model(&PostLsfc{}).Where("user_id = ? and is_like = ?", uid,1).Limit(100).Offset(page * 100).Find(&items).Error
	if err != nil {
		return []PostLsfc{}, err
	}

	return items, nil
}
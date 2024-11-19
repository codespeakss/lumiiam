package model

type User struct {
	Id       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"type:varchar(255);unique_index"`
	CreateAt int64  `json:"create_at"`
	Password string `json:"-"` // 不序列化密码
}

package models

type User struct {
	Id       string `gorm:"size:255;primary_key;not null;unique" json:"id"`
	UserName string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
}

type RegReq struct {
	UserName string `gorm:"size:255;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
}

type LoginReq struct {
	UserName string `gorm:"size:255;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
}

type LoginResult struct {
	UserName string `gorm:"size:255;not null" json:"username"`
	Token    string `gorm:"size:255;not null" json:"token"`
}

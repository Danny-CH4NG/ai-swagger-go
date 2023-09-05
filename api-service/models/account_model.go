package models

type Account struct {
	ID       int64  `json:"id" gorm:"primary_key;auto_increment"`
	Username string `json:"username" gorm:"unique_index;not null"`
	Password string `json:"password" gorm:"not null"`
}

func (Account) TableName() string {
	return "account"
}

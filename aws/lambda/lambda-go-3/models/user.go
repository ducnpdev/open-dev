package models

type UserModel struct {
	ID       uint   `gorm:"primarykey"`
	UserName string `gorm:"column:user_name" bson:"user_name"`
	Name     string `gorm:"column:name" bson:"name"`
	Phone    string `gorm:"column:phone" bson:"phone"`
}

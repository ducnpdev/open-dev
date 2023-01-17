package models

type UserModel struct {
	ID       uint   `gorm:"primarykey"`
	UserName string `gorm:"column:username" bson:"username"`
	Email    string `gorm:"column:email" bson:"email"`
	Phone    string `gorm:"column:phone" bson:"phone"`
}

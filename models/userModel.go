package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string // Храним как hash, не в открытом виде!
}

func (User) TableName() string {
	return "user"
}

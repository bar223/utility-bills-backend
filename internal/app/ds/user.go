package ds

// User — пользователь системы. Таблица users.
type User struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Login       string `gorm:"type:varchar(50);unique;not null" json:"login"`
	Password    string `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator bool   `gorm:"type:boolean;not null;default:false" json:"isModerator"`
}

func (User) TableName() string {
	return "users"
}

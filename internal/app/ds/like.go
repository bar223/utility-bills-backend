package ds

// Like — связь м-м «пользователь — услуга»: отметка «нравится» конкретного
// пользователя на конкретном коммунальном ресурсе. Таблица likes.
type Like struct {
	ID uint `gorm:"primaryKey"`

	UserID             uint `gorm:"not null;uniqueIndex:idx_user_resource"`
	CommunalResourceID uint `gorm:"not null;uniqueIndex:idx_user_resource"`

	User             User             `gorm:"foreignKey:UserID"`
	CommunalResource CommunalResource `gorm:"foreignKey:CommunalResourceID"`
}

func (Like) TableName() string {
	return "likes"
}

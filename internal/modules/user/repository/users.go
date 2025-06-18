package repository

type RepositoryUser struct {
	UserName   string
	FirstName  string
	LastName   string
	Email      string
	Password   string
	Phone      string
	UserStatus int
}

type RepositoryUserArray struct {
	UserArray []RepositoryUser
}

type User struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Username   string `gorm:"type:text;unique;not null"`
	FirstName  string `gorm:"type:text;not null"`
	LastName   string `gorm:"type:text;not null"`
	Email      string `gorm:"type:text;not null"`
	Password   string `gorm:"type:text;not null"`
	Phone      string `gorm:"type:text;not null"`
	UserStatus int    `gorm:"column:userstatus"`
}

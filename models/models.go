package models

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignUp struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
	Hobbie   string `json:"hobbie"`
}

type User struct {
	UserId   uint    `gorm:"foreignKey" json:"user_id"`
	Name     string  `json:"name" binding:"required"`
	Password string  `json:"password" binding:"required,min=8"`
	Email    string  `json:"email" binding:"required,email"`
	Hobbie   string  `json:"hobbie"`
	Notes    []Notes `gorm:"foreignKey:UserId"`
}

type Notes struct {
	NotesID string `gorm:"primaryKey" json:"notes_id"`
	Notes   string `json:"notes"`
	UserId  uint   `json:"user_id"`
}

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
	UserId   uint `gorm:"primaryKey;column:user_id"` // mark as PK
	Name     string
	Password string `json:"-"`
	Email    string `gorm:"uniqueIndex;not null"`
	Hobbie   string
	// Tell GORM which field on Notes is the FK that references this PK
	Notes []Notes `gorm:"foreignKey:UserId;references:UserId;constraint:OnDelete:CASCADE"`
}

type Notes struct {
	NotesID     uint   `gorm:"primaryKey;column:notes_id"` // make it uint to match your int parsing
	Notes_Title string `gorm:"column:notes_title"`
	Notes_Body  string `gorm:"column:notes_body"`
	UserId      uint   `gorm:"column:user_id;index;not null"` // matches User.UserId
}

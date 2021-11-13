package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Token    string `json:"token"`
}

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{db}
}

func (u *Users) Find(users *[]User) error {
	result := u.db.Find(&users)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
	return result.Error
}

func (u *Users) FindByEmail(user *User, email string) error {
	result := u.db.Where("email = ?", email).First(&user)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
	return result.Error
}

func (u *Users) First(user *User, id int) error {
	result := u.db.First(user, id)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
	return result.Error
}

func (u *Users) Create(user *User) error {
	result := u.db.Create(user)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
	return result.Error
}

func (u *Users) Update(user *User, id int) error {
	var tmp User
	result := u.db.First(&tmp, id)
	if result.Error != nil {
		return result.Error
	}

	user.ID = id
	result = u.db.Save(user)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
	return result.Error
}

func (u *Users) Delete(user *User, id int) error {
	result := u.db.Delete(user, id)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
	return result.Error
}

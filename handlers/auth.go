package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mwazovzky/loki/models"
	"gorm.io/gorm"
)

type AuthHandlers struct {
	users *models.Users
}

type UserLogin struct {
	Email    string
	Password string
}

func NewAuthHandlers(db *gorm.DB) *AuthHandlers {
	users := models.NewUsers(db)
	return &AuthHandlers{users}
}

// curl -X POST localhost:3000/api/login -d '{"email":"vasya@example.com","password":"secret"}'
func (ah *AuthHandlers) Login(rw http.ResponseWriter, r *http.Request) {
	log.Println("Auth Login Request")

	userLogin := UserLogin{}
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		http.Error(rw, "Unable to read request data", http.StatusBadRequest)
		return
	}

	user := models.User{}
	err = ah.users.FindByEmail(&user, userLogin.Email)
	if err != nil {
		http.Error(rw, "Wrong emain or password", http.StatusBadRequest)
		return
	}

	err = verifyPassword(userLogin.Password, user.Password)
	if err != nil {
		http.Error(rw, "Wrong emain or password", http.StatusBadRequest)
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(rw, "Can't generate token", http.StatusInternalServerError)
		return
	}

	user.Token = token

	err = ah.users.Update(&user, user.ID)
	if err != nil {
		http.Error(rw, "Model not found", http.StatusNotFound)
		return
	}

	e := json.NewEncoder(rw)
	err = e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func generateToken() (string, error) {
	length := 10

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func verifyPassword(password string, userPassword string) error {
	if password != userPassword {
		return errors.New("wrong emain or password")
	}
	return nil
}

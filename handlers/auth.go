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

type Token struct {
	Value string
}

func NewAuthHandlers(db *gorm.DB) *AuthHandlers {
	users := models.NewUsers(db)
	return &AuthHandlers{users}
}

func (ah *AuthHandlers) Login(rw http.ResponseWriter, r *http.Request) {
	log.Println("Auth Login Request")

	userLogin := models.User{}
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

	if userLogin.Password != user.Password {
		http.Error(rw, "Wrong emain or password", http.StatusBadRequest)
		return
	}

	tokenValue, err := generateToken()
	if err != nil {
		http.Error(rw, "Can't generate token", http.StatusInternalServerError)
		return
	}
	token := Token{tokenValue}
	user.Token = token.Value

	err = ah.users.Update(&user, user.ID)
	if err != nil {
		http.Error(rw, "Model not found", http.StatusNotFound)
		return
	}

	e := json.NewEncoder(rw)
	err = e.Encode(token)
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

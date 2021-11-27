package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"loki/models"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type AuthHandlers struct {
	users    *models.Users
	validate *validator.Validate
}

type UserLogin struct {
	Email    string
	Password string
}

func NewAuthHandlers(db *gorm.DB, v *validator.Validate) *AuthHandlers {
	users := models.NewUsers(db)
	return &AuthHandlers{users, v}
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

func (ah *AuthHandlers) Logout(rw http.ResponseWriter, r *http.Request) {
	log.Println("Auth Logout Request")

	authHeader := r.Header.Get("Authorization")
	auth := strings.Split(authHeader, " ")
	token := auth[1]

	err := ah.users.DeleteToken(token)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	e := json.NewEncoder(rw)
	err = e.Encode("success")
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (ah *AuthHandlers) Register(rw http.ResponseWriter, r *http.Request) {
	log.Println("Users Create Request")

	newUser := models.User{}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(rw, "Missing data", http.StatusBadRequest)
		return
	}

	err = ah.validateUser(&newUser)
	if err != nil {
		http.Error(rw, "Validation", http.StatusUnprocessableEntity)
		return
	}

	var user models.User
	err = ah.users.FindByEmail(&user, newUser.Email)
	if err == nil {
		http.Error(rw, "User with this email already exists", http.StatusBadRequest)
		return
	}

	ah.users.Create(&newUser)

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(newUser)
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

func (ah *AuthHandlers) validateUser(user *models.User) error {
	err := ah.validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Println("Validation error")
			log.Println(err.Namespace())
			log.Println(err.Field())
			log.Println(err.StructNamespace())
			log.Println(err.StructField())
			log.Println(err.Tag())
			log.Println(err.ActualTag())
			log.Println(err.Kind())
			log.Println(err.Type())
			log.Println(err.Value())
			log.Println(err.Param())
		}
	}

	return err
}

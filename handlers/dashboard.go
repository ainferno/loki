package handlers

import (
	"encoding/json"
	"log"
	"loki/models"
	"net/http"

	"gorm.io/gorm"
)

type DashboardHandlers struct {
	users *models.Users
}

type Example struct {
	Content string
}

func NewDashboardHandlers(db *gorm.DB) *DashboardHandlers {
	users := models.NewUsers(db)
	return &DashboardHandlers{users}
}

func (dh *DashboardHandlers) Index(rw http.ResponseWriter, r *http.Request) {
	log.Println("Dashboard index")

	e := json.NewEncoder(rw)
	err := e.Encode(Example{"index"})
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

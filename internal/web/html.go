package web

import (
	"html/template"
	"net/http"

	"github.com/chnmk/order-info-l0/internal/models"
)

func DisplayTemplate(w http.ResponseWriter, r *http.Request) {
	// data, _ := test.ReadModelFile()
	var data models.Order
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}

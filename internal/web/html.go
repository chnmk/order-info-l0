package web

import (
	"html/template"
	"net/http"

	"github.com/chnmk/order-info-l0/test"
)

func DisplayTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, test.E)
}
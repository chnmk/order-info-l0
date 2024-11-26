package transport

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/chnmk/order-info-l0/internal/models"
)

// Возвращает страницу поиска заказа.
func DisplayPage(w http.ResponseWriter, r *http.Request) {
	slog.Info("incoming request to: /")

	slog.Info("executing template...")

	var data models.Order
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		slog.Info(err.Error())
		return
	}

	tmpl.Execute(w, data)

	slog.Info("template successfully executed")

	slog.Info("request successfull")
}

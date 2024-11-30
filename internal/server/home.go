package server

import (
	"html/template"
	"log/slog"
	"net/http"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Возвращает страницу поиска заказа.
func DisplayPage(w http.ResponseWriter, r *http.Request) {
	slog.Info("incoming request to: /")

	slog.Info("executing template...")

	var data models.Order
	tmpl, err := template.ParseFiles(cfg.TemplatePath)
	if err != nil {
		slog.Info(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse html template"))
		return
	}

	tmpl.Execute(w, data)

	slog.Info("template successfully executed")

	slog.Info("request successfull")
}

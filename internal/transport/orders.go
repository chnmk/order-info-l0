package transport

import (
	"log/slog"
	"net/http"
	"text/template"

	"github.com/chnmk/order-info-l0/internal/database"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	slog.Info("incoming request to: /orders")

	if r.Method != http.MethodGet {
		slog.Info("invalid request method")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid requiest method"))
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Info("invalid request: no id")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request: no id"))
		return
	}

	result := database.SelectOrderById(database.DB, id)
	/*
		resp, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}
	*/

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, result)

	slog.Info("request successfull")
}

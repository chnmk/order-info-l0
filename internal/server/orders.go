package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"

	"github.com/chnmk/order-info-l0/internal/config"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	slog.Info("incoming request to: /orders")

	if r.Method != http.MethodGet {
		slog.Info("invalid request method")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid requiest method"))
		return
	}

	// Получает id заказа (1, 2, 3..., т.е. порядковый id из кэша, а не order_uid).
	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Info("invalid request: no id")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request: no id"))
		return
	}

	conv_id, err := strconv.Atoi(id)
	if err != nil {
		slog.Info("invalid request: id should be a number")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request: id should be a number"))
		return
	}

	// Получает сам заказ из памяти.
	result := config.Data.ReadByID(conv_id)
	if result.UID == "" {
		slog.Info("invalid request: order not found")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request: order not found"))
		return
	}

	// Получает параметр format.
	// При format == "html" возвращает страницу с данными.
	// При любом другом значении возвращает JSON.
	format := r.URL.Query().Get("format")

	if format == "html" {
		// Возвращаем страницу.
		slog.Info("executing template...")

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			slog.Info(err.Error())
			return
		}

		tmpl.Execute(w, result)

		slog.Info("template successfully executed")

	} else {
		// Формат html не указан, возвращаем JSON.
		if format == "" {
			slog.Info("response format not defined")
		}

		slog.Info("sending response...")

		resp, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}

	slog.Info("request successfull")
}

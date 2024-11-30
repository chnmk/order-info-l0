package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

// Заказ можно получить либо по id (порядковый номер заказа), либо по order_uid.
// По умолчанию в переменных окружения указано получать его по id (SERVER_GET_ORDER_BY_ID = 1).
func GetOrder(w http.ResponseWriter, r *http.Request) {
	slog.Info("incoming request to: /orders")

	// Проверяет метод запроса.
	if r.Method != http.MethodGet {
		slog.Info("invalid request method")

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid requiest method"))
		return
	}

	// Получает параметр id.
	id := r.URL.Query().Get("id")
	if id == "" {
		slog.Info("invalid request: no id")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request: no id"))
		return
	}

	var order models.OrderStorage

	// Если в окружении указано получать заказ по id...
	if cfg.GetOrderById {
		// Cначала стоит проверить, что он является числом.
		conv_id, err := strconv.Atoi(id)
		if err != nil {
			slog.Info("invalid request: id should be a number")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request: id should be a number"))
			return
		}

		// Получает сам заказ из памяти.
		order = cfg.Data.ReadByID(conv_id)
		if order.UID == "" {
			slog.Info("invalid request: order not found")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request: order not found"))
			return
		}

	} else {
		// Если в окружении указано получать заказ по order_uid...
		order = cfg.Data.ReadByUID(id)
		if order.UID == "" {
			slog.Info("invalid request: order not found")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request: order not found"))
			return
		}
	}

	// Получает параметр format. При format == "html" возвращает страницу с данными.
	// При любом другом значении возвращает JSON.
	format := r.URL.Query().Get("format")

	if format == "html" {
		slog.Info("executing template...")

		// Обрабатывает файл с шаблоном.
		tmpl, err := template.ParseFiles(cfg.TemplatePath)
		if err != nil {
			slog.Info(err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to parse html template"))
			return
		}

		// Декодирует данные для отображения на странице.
		var unmarshalled models.Order

		err = json.Unmarshal(order.Order, &unmarshalled)
		if err != nil {
			slog.Info(err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to parse order data"))
			return
		}

		tmpl.Execute(w, unmarshalled)

		slog.Info("template successfully executed")

	} else {
		// Формат html не указан, возвращает JSON.
		if format == "" {
			slog.Info("response format not defined")
		}

		slog.Info("sending response...")

		w.Write(order.Order)
	}

	slog.Info("request successfull")
}

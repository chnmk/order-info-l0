package transport

import (
	"encoding/json"
	"net/http"

	"github.com/chnmk/order-info-l0/internal/database"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid method"))
		return
	}

	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected id"))
		return
	}

	result := database.SelectOrderById(database.DB, id)

	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

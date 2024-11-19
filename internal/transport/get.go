package transport

import (
	"encoding/json"
	"net/http"

	"github.com/chnmk/order-info-l0/test"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	/*
		id := r.FormValue("order_uid")
		if id == "" {
			return
		}
	*/

	resp, err := json.Marshal(test.E)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

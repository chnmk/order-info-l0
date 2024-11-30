package consumer

import (
	"encoding/json"
	"testing"

	"github.com/chnmk/order-info-l0/internal/models"
)

func TestGoFakeGenerator(t *testing.T) {
	var order models.Order

	bytes := goFake()

	err := json.Unmarshal(bytes, &order)
	if err != nil {
		t.Fatalf("unexpected error while unmarshalling models.Order struct: %s", err.Error())
	}

	if order.Order_uid == "" {
		t.Fatalf("expected non-empty values in generated order")

	}

}

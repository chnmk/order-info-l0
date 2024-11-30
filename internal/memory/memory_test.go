package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
)

var M models.Storage

func TestNewMemory(t *testing.T) {
	M = NewStorage(context.Background(), M)
	if M == nil {
		t.Fatalf("created memory storage shouldn't be nil")
	}

	var m2 models.Storage
	m2 = NewStorage(context.Background(), m2)
	if m2 != nil {
		t.Fatalf("memory storage should only be created once")
	}
}

func TestAddMessage(t *testing.T) {
	M = NewStorage(context.Background(), M)

	var data models.Order
	gofakeit.Struct(&data)

	msg, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("error marshalling message: %s", err.Error())
	}

	order1 := models.OrderStorage{ID: 1, UID: "1", Date_created: "1", Order: msg}
	_ = order1.ID

	order2 := M.AddOrder("1", "1", msg)

	if order1.UID != order2.UID || order1.Date_created != order2.Date_created || !bytes.Equal(order1.Order, order2.Order) {
		t.Errorf("order insert failed: data doesn't match")
	}

	order3 := M.AddOrder("2", "2", msg)

	if order2.ID >= order3.ID {
		t.Errorf("newer order should have a bigger id value")
	}
}

func TestReadMessage(t *testing.T) {
	M = NewStorage(context.Background(), M)

	var data models.Order
	gofakeit.Struct(&data)

	msg, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("error marshalling message: %s", err.Error())
	}

	order1 := M.AddOrder("11", "11", msg)
	order2 := M.AddOrder("22", "22", msg)

	if M.ReadByID(order1.ID).ID == M.ReadByID(order2.ID).ID ||
		M.ReadByID(order1.ID).UID == M.ReadByID(order2.ID).UID {
		t.Errorf("expected different messages when reading by id")
	}

	if M.ReadByUID(order1.UID).ID == M.ReadByUID(order2.UID).ID ||
		M.ReadByUID(order1.UID).UID == M.ReadByUID(order2.UID).UID {
		t.Errorf("expected different messages when reading by uid")
	}

	if M.ReadByID(order1.ID).ID != M.ReadByUID(order1.UID).ID {
		t.Errorf("expected to get the same message when reading with different methods")
	}
}

func TestValidateMessage(t *testing.T) {
	var order models.Order
	gofakeit.Struct(&order)
	err := ValidateMsg(order)
	if err != nil {
		t.Errorf("expected order to be valid")
	}

	var emptyOrder models.Order
	err = ValidateMsg(emptyOrder)
	if err == nil {
		t.Errorf("expected empty order to not be valid")
	}

	var emptyOrderWithId models.Order
	emptyOrderWithId.Order_uid = "test"
	err = ValidateMsg(emptyOrderWithId)
	if err == nil {
		t.Errorf("expected order with missing values to not be valid")
	}
}

func TestRestoreData(t *testing.T) {
	M = NewStorage(context.Background(), M)

	cfg.DB = &models.MockDatabase{}

	o1 := models.OrderStorage{ID: 111, UID: "111", Date_created: "111", Order: []byte("test")}
	cfg.DB.InsertOrder(context.Background(), o1)

	o2 := models.OrderStorage{ID: 222, UID: "222", Date_created: "222", Order: []byte("test2")}
	cfg.DB.InsertOrder(context.Background(), o2)

	M.RestoreData(context.Background())

	if M.ReadByID(111).UID == "" {
		t.Fatalf("expected restored orders to contain data")
	}

	if M.ReadByID(111).ID != 111 {
		t.Fatalf("expected valid order data")
	}
}

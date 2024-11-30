package database

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/chnmk/order-info-l0/internal/models"
)

func TestDatabaseExample(t *testing.T) {
	mydb := &models.MockDatabase{}

	var m models.OrderStorage
	gofakeit.Struct(&m)
	mydb.InsertOrder(context.Background(), m)

	newOrder := mydb.RestoreData(context.Background())

	if newOrder[0].ID != m.ID || newOrder[0].Order == nil {
		t.Fatalf("expected to get the same data")
	}
}

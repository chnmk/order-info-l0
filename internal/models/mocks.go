package models

import "context"

type MockDatabase struct {
	Orders []OrderStorage
}

func (m *MockDatabase) Close() {
	// Do nothing
}

func (m *MockDatabase) Ping(ctx context.Context) {
	// Do nothing

}

func (m *MockDatabase) CreateTables(ctx context.Context) {
	// Do nothing
}

func (m *MockDatabase) InsertOrder(ctx context.Context, order OrderStorage) {
	m.Orders = append(m.Orders, order)
}

func (m *MockDatabase) RestoreData(ctx context.Context) []OrderStorage {
	return m.Orders
}

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/memory"
	"github.com/chnmk/order-info-l0/internal/models"
)

func TestGetIndex(t *testing.T) {
	cfg.TemplatePath = "./index.html"

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(DisplayPage)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if rec.Body == nil {
		t.Errorf("expected server response")
	}

	if !strings.Contains(rec.Body.String(), "<html") {
		t.Errorf("expected html file")
	}
}

func TestGetHTML(t *testing.T) {
	cfg.Data = memory.NewStorage(context.Background(), cfg.Data)

	var data models.Order
	gofakeit.Struct(&data)

	msg, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("error marshalling message: %s", err.Error())
	}

	order1 := cfg.Data.AddOrder("10101uid", "1", msg)
	orderid := order1.ID

	// ID

	os.Setenv("SERVER_GET_ORDER_BY_ID", "1")
	cfg1 := cfg.NewConfig()
	cfg.TemplatePath = "./index.html"

	req := httptest.NewRequest("GET", fmt.Sprintf("/order?format=html&id=%d", orderid), nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(GetOrder)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("orderid: %d", orderid)
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if rec.Body == nil {
		t.Errorf("expected server response")
	}

	if !strings.Contains(rec.Body.String(), "Цена:") {
		t.Errorf("expected html file with data")
	}

	// UID

	os.Setenv("SERVER_GET_ORDER_BY_ID", "0")
	cfg1.InitEnv()

	req = httptest.NewRequest("GET", "/order?format=html&id=10101uid", nil)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(GetOrder)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if rec.Body == nil {
		t.Errorf("expected server response")
	}

	if !strings.Contains(rec.Body.String(), "Цена:") {
		t.Errorf("expected html file with data")
	}

	// Invalid

	req = httptest.NewRequest("GET", "/order?format=html&id=10101uidefekjfbjINVALID", nil)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(GetOrder)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got: %d, error: %s", rec.Code, rec.Body)
	}
}

func TestGetJSON(t *testing.T) {
	cfg.Data = memory.NewStorage(context.Background(), cfg.Data)

	var data models.Order
	gofakeit.Struct(&data)

	msg, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("error marshalling message: %s", err.Error())
	}

	order2 := cfg.Data.AddOrder("20202uid", "2", msg)
	orderid := order2.ID

	// ID

	os.Setenv("SERVER_GET_ORDER_BY_ID", "1")
	cfg1 := cfg.NewConfig()
	cfg1.InitEnv()
	cfg.TemplatePath = "./index.html"

	req := httptest.NewRequest("GET", fmt.Sprintf("/order?id=%d", orderid), nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(GetOrder)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if rec.Body == nil {
		t.Errorf("expected server response")
	}

	if !strings.Contains(rec.Body.String(), "{\"order_uid") {
		t.Log(rec.Body.String())
		t.Errorf("expected json with data")
	}

	// UID

	os.Setenv("SERVER_GET_ORDER_BY_ID", "0")
	cfg1.InitEnv()

	req = httptest.NewRequest("GET", "/order?id=20202uid", nil)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(GetOrder)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if rec.Body == nil {
		t.Errorf("expected server response")
	}

	if !strings.Contains(rec.Body.String(), "{\"order_uid") {
		t.Errorf("expected html file with data")
	}

	// Invalid

	req = httptest.NewRequest("GET", "/order?id=20202uidefekjfbjINVALID", nil)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(GetOrder)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got: %d, error: %s", rec.Code, rec.Body)
	}
}

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/chnmk/order-info-l0/internal/models"
	"github.com/segmentio/kafka-go"
)

// Генерирует случайные корректные данные с использованием Gofakeit.
//
// Данные не претендуют на реалистичность (не совпадают цены, места и так далее), но это и не рандомные символы.
// Такой вариант приемлем, поскольку по ТЗ эти данные пока никак не обрабатываются, но отобразить их нужно.
func PublishTestData() {
	slog.Info("creating new kafka writer for fake data...")

	w := &kafka.Writer{
		Addr:  kafka.TCP(fmt.Sprintf("%s:%s", cfg.Env["KAFKA_PROTOCOL"], cfg.Env["KAFKA_PORT"])),
		Topic: cfg.Env["KAFKA_TOPIC"],
	}

	slog.Info("generating fake data...")

	// Запись сообщений работает 1000 секунд. TODO: правильнее отменять через контекст?
	for i := 0; i < 1000; i++ {
		var m models.Order

		gofakeit.Struct(&m)
		data, err := json.Marshal(m)
		if err != nil {
			slog.Error(err.Error())
		}

		err = w.WriteMessages(context.Background(),
			kafka.Message{Value: data},
		)
		if err != nil {
			slog.Error("Failed to write messages: " + err.Error())
		} else {
			slog.Info("writing successful!")
		}

		time.Sleep(time.Second)
	}

	if err := w.Close(); err != nil {
		slog.Error("Failed to close writer: " + err.Error())
	}

	slog.Info("fake data generation stopped")
}

func GofakeInit() {
	// Пример кастомной функции для генерации приближенных к реальности данных.
	gofakeit.AddFuncLookup("wbdate", gofakeit.Info{
		Category:    "custom",
		Description: "random date string",
		Example:     "2021-11-26T06:22:19Z",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			// Отнимает один месяц от сегодняшнего дня, переводит в формат unix.
			min := time.Now().AddDate(0, -1, 0).Unix()
			// Прибавляет к min случайное значение, таким образом получает дату между сегодня и месяц назад.
			unix := min + rand.Int63n(time.Now().Unix()-min)

			time := time.Unix(unix, 0).Format("2006-01-02T15:04:05Z")
			return time, nil
		},
	})
}

// Читает файл model.json, возвращает его в виде структуры и []byte.
func ReadModelFile() (models.Order, []byte) {
	var E models.Order
	content, err := os.ReadFile("test/model.json")
	if err != nil {
		slog.Error(err.Error())
	}

	err = json.Unmarshal(content, &E)
	if err != nil {
		slog.Error(err.Error())
	}

	return E, content
}

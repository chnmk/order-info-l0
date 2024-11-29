package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
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
func publishExampleData(ctx context.Context) {
	defer cfg.ExitWg.Done()

	slog.Info("creating new kafka writer...")

	w := &kafka.Writer{
		Addr:  cfg.KafkaWriterAddr,
		Topic: cfg.KafkaWriterTopic,
	}

	slog.Info("generating fake data...")

	// Запись сообщений работает постоянно.
	for {
		select {

		case <-ctx.Done():
			if err := w.Close(); err != nil {
				slog.Error(
					"failed to close writer",
					"err", err.Error(),
				)
			}

			slog.Info("fake data generation stopped")
			return

		default:
			err := w.WriteMessages(ctx, kafka.Message{Value: goFake()})
			if err != nil {
				slog.Error(
					"failed to write messages",
					"err", err.Error(),
				)
				break
			}

			slog.Info("writing successful!")
		}
	}
}

// Пример кастомной функции для генерации приближенных к реальности данных.
func goFakeInit() {
	gofakeit.AddFuncLookup("wbdate", gofakeit.Info{
		Category:    "custom",
		Description: "random date string",
		Example:     "2021-11-26T06:22:19Z",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			// Отнимает один месяц от сегодняшнего дня, переводит в формат unix.
			min := time.Now().AddDate(0, 0, -7).Unix()
			// Прибавляет к min случайное значение, таким образом получает дату между сегодня и месяц назад.
			unix := min + rand.Int63n(time.Now().Unix()-min)

			time := time.Unix(unix, 0).Format("2006-01-02T15:04:05Z")
			return time, nil
		},
	})
}

// Генерирует один заказ.
func goFake() []byte {
	var m models.Order

	gofakeit.Struct(&m)
	data, err := json.Marshal(m)
	if err != nil {
		slog.Error(err.Error())
	}

	return data
}

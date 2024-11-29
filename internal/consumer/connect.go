package consumer

import (
	"log/slog"
	"time"

	cfg "github.com/chnmk/order-info-l0/internal/config"
	"github.com/segmentio/kafka-go"
)

// Проверяет подключение к Kafka, пытается подключиться указанное количество раз.
// Затем создает горутины для чтения сообщений и, при необходимости, для записи сгенерированных данных.
func Connect() {
	var conn *kafka.Conn
	var err error

	slog.Info("connecting to kafka...")

	// Пытается подключиться к Kafka.
	for i := 0; i < cfg.KafkaReconnectAttempts; i++ {
		conn, err = kafka.DialLeader(cfg.ExitCtx,
			cfg.KafkaInitNework,
			cfg.KafkaInitAddress,
			cfg.KafkaInitTopic,
			cfg.KafkaInitPartition,
		)

		if err != nil {
			slog.Error(err.Error())
		} else {
			slog.Info("kafka connection successful")
			break
		}

		time.Sleep(1 * time.Second)
	}
	if err != nil {
		slog.Error(
			"kafka connection failed",
			"err", err.Error(),
		)
		cfg.Exit()
	}

	conn.Close()

	// Создает горутины для чтения сообщений.
	// TODO: посмотреть, увеличится ли скорость чтения от наличия нескольких горутин.
	for i := 0; i < cfg.KafkaReaderGoroutines; i++ {
		cfg.ExitWg.Add(1)
		go newReader().Read()
	}

	// Создает горутины для записи сгенерированных сообщений.
	if cfg.KafkaWriteExamples {
		goFakeInit()

		for i := 0; i < cfg.KafkaWriterGoroutines; i++ {
			cfg.ExitWg.Add(1)
			go publishExampleData()
		}

	}
}

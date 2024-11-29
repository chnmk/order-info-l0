package config

import (
	"github.com/chnmk/order-info-l0/internal/models"
)

/*
	Устанавливает переменные для пакета memory.
*/

var (
	Data                    models.Storage          // Глобальное хранилище данных.
	RestoreData             bool                    // При значении true (по умолчанию) стоит восстанавливатьы данные из БД при запуске.
	CleanupInterval         int                     // Интервал выполнения функции для очистки кэша, в минутах.
	OrdersLimit             int                     // Количество сообщений, которые можно хранить до их удаления.
	MessagesChan            chan models.MessageData // Канал для отправления сообщений из Kafka в пул обработчиков.
	MemoryHandlerGoroutines int                     // Количество обрабочиков, которые будут читать сообщения из канала.
)

// Получает глобальные переменные для пакета memory.
func getMemoryVars() {
	RestoreData = envToBool("MEMORY_RESTORE_DATA")
	CleanupInterval = envToInt("MEMORY_CLEANUP_MINUTES_INTERVAL")
	OrdersLimit = envToInt("MEMORY_ORDERS_LIMIT")
	OrdersLimit = envToInt("MEMORY_HANDLER_GOROUTINES")

	MessagesChan = make(chan models.MessageData, MemoryHandlerGoroutines)
}

package config

import (
	"github.com/chnmk/order-info-l0/internal/models"
)

/*
	Устанавливает переменные для пакета memory.
*/

var (
	Data            models.Storage // Глобальное хранилище данных.
	RestoreData     bool           // При значении true (по умолчанию) стоит восстанавливатьы данные из БД при запуске.
	CleanupInterval int
	OrdersLimit     int
)

// Получает глобальные переменные для пакета memory.
func getMemoryVars() {
	RestoreData = envToBool("MEMORY_RESTORE_DATA")
	CleanupInterval = envToInt("MEMORY_CLEANUP_MINUTES_INTERVAL")
	OrdersLimit = envToInt("MEMORY_ORDERS_LIMIT")
}

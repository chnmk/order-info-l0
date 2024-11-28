package config

import "github.com/chnmk/order-info-l0/internal/models"

/*
	Устанавливает переменные для пакета memory.
*/

var (
	Data        models.Storage // Глобальное хранилище данных.
	RestoreData bool           // При значении true (по умолчанию) стоит восстанавливатьы данные из БД при запуске.
)

// Получает глобальные переменные для пакета memory.
func getMemoryVars() {
	RestoreData = envToBool("MEMORY_RESTORE_DATA")
}

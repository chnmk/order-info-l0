package models

type Config interface {
	InitEnv()                      // Устанавливает стандартные значения переменных окружения.
	ReadEnv()                      // Получает переменные из окружения.
	Get(key string) (value string) // Получает значение из мапы переменных по названию.

}

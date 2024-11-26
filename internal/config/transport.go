package config

/*
	Устанавливает переменные для пакета transport.
*/

var ServerPort string

// Получает глобальные переменные для пакета transport.
func getTransportVars() {
	ServerPort = Env.Get("SERVER_PORT")
}

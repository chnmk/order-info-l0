package config

/*
	Устанавливает переменные для пакета server.
*/

var ServerPort string

// Получает глобальные переменные для пакета server.
func getTransportVars() {
	ServerPort = Env.Get("SERVER_PORT")
}

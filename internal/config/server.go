package config

/*
	Устанавливает переменные для пакета server.
*/

var (
	ServerPort   string
	GetOrderById bool // Если true, то заказ нужно запрашивать по его ID (порядковый номер). Иначе - по order_uid.
	TemplatePath = "internal/server/index.html"
)

// Получает глобальные переменные для пакета server.
func getServerVars() {
	ServerPort = Env.Get("SERVER_PORT")
	GetOrderById = envToBool("SERVER_GET_ORDER_BY_ID")
}

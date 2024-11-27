# internal/consumer

Пакет для работы с Kafka. 

Использует consumer group для чтения сообщений. Вот несколько причин:
- Возможность легко использовать коммиты и не перечитывать старые сообщения, а также не терять их при падении сервиса, поскольку коммит происходит только после добавления в БД.
- Возможность при желании работать с несколькими reader в горутинах.

## connect.go

Проверяет подключение, пытается подключиться указанное количество раз, при неудаче завершает работу сервиса с ошибкой. Создает горутины для чтения сообщений (`Read()` из файла `reader.go`) и, при необходимости, для записи сгенерированных данных (`publishExampleData()` из файла `example.go`).

## consumer.go

Имплементирует интерфейс конфига из пакета `models`, возвращает новый reader.

## example.go

Генерирует данные для проверки работы сервиса и записывает их в Kafka.

## reader.go

Читает сообщения из Kafka и отправляет их хендлерам из пакета memory.

## Тесты
Проверяют работу генератора сообщений. Помимо него здесь используется только базовый функционал сторонних библиотек, так что писать под них дополнительные тесты не считаю приоритетной задачей.
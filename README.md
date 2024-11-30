# Демонстрационный сервис Go + Kafka + PostgreSQL
Получает сообщения из Kafka и хранит их в памяти. Использует PostgreSQL для резервного хранения и восстановления данных в случае проблем с сервисом.

## Содержание
1. [Логика работы](#логика-работы)
2. [Сторонние библиотеки](#сторонние-библиотеки)
3. [Docker Compose](#docker-compose)
4. [Запуск без Docker](#запуск-без-docker)
5. [Генерация данных](#генерация-данных)
6. [Тестирование](#тестирование)
7. [Переменные окружения](#переменные-окружения)
8. [Бенчмарки](#бенчмарки)

## Логика работы
Получает сообщения из Kafka с использованием consumer group. Может создавать несколько читателей сообщений в горутинах (переменная окружения `KAFKA_READER_GOROUTINES`).

Сообщения отправляются через канал в пул обработчиков (переменная `MEMORY_HANDLER_GOROUTINES`). Сообщения проверяются на соответствие ожидаемой модели данных, записываются в память и в базу данных, после чего коммитятся в Kafka. 

Один раз в указанный промежуток времени (`MEMORY_CLEANUP_MINUTES_INTERVAL`) происходит очистка памяти: старые сообщения удаляются, если превышен лимит сообщений (`MEMORY_ORDERS_LIMIT`), или если заказ определен как устаревший (по полю date_created).

В базе данных заказы хранятся в формате JSONB, поскольку на текущий момент они используются только для восстановления данных в память.

HTTP-сервер возвращает заказ из памяти либо по его ID (порядковый номер заказа), либо по полю order_uid, в зависимости от переменной окружения `SERVER_GET_ORDER_BY_ID` (тип bool, по умолчанию `1`).

Логирование реализовано через slog.

Есть возможность генерировать случайные данные с помощью gofakeit и отправлять их в Kafka для проверки работы сервиса (переменные окружения `KAFKA_WRITE_EXAMPLES` (тип bool, по умолчанию `0`), `KAFKA_WRITER_GOROUTINES`)

## Сторонние библиотеки
- [pgx v5.7.1](https://github.com/jackc/pgx/v5) - для подключения к PostgreSQL
- [segmentio/kafka-go v0.4.47](https://github.com/segmentio/kafka-go) - для подключения к Kafka
- [gofakeit 7.1.2](https://github.com/brianvoe/gofakeit) - для генерации данных для проверки сервиса
- [godotenv v1.5.1](https://github.com/joho/godotenv) -  для получения переменных окружения из .env файла

## Docker Compose

Сервис запускается вместе с PostgreSQL и Kafka (файл ./compose.yaml). Для запуска рекомендуется использовать команду `docker compose up --attach orders`. Необходимые переменные окружения:

```
# Указаны значения по умолчанию
POSTGRES_USER=user
POSTGRES_PASSWORD=12345
POSTGRES_DB=orders
```

Работу веб-интерфейса можно проверить по адресу `localhost:3000` (порт по умолчанию).

## Запуск без Docker
Запуск производится командой `go run .`

Список переменных окружения для подключения к Kafka и PostgreSQL можно посмотреть в разделе [Переменные окружения](#переменные-окружения).


## Генерация данных
Проверить работу сервиса можно с помощью сгенерированных данных. Необходимо указать следующие переменные окружения:

```
KAFKA_WRITE_EXAMPLES=1 # true. По умолчанию 0.
KAFKA_WRITER_GOROUTINES=1 # значение по умолчанию. Чем больше, тем быстрее будут генерироваться данные.
```

## Тестирование

Написаны автотесты на пакеты config, consumer, database, memory, server. Запуск тестов:

```
go test ./...
```

## Переменные окружения
Список переменных окружения со стандартными значениями:

```
# Переменные DB необходимо указать для запуска через Docker Compose
DB_PROTOCOL=postgres
DB_HOST=postgres
DB_PORT=5432 

# Строка подключения к БД формируется следующим образом:
# DB_PROTOCOL://POSTGRES_USER:POSTGRES_PASSWORD@DB_HOST:DB_PORT/POSTGRES_DB
POSTGRES_USER=user
POSTGRES_PASSWORD=12345
POSTGRES_DB=orders

KAFKA_NETWORK=tcp
KAFKA_PROTOCOL=kafka
KAFKA_PORT=9092

KAFKA_TOPIC=go-orders
KAFKA_GROUP_ID=go-orders-1
KAFKA_MAX_BYTES=20000 # 20kb
KAFKA_RECONNECT_ATTEMPTS=20

KAFKA_READER_GOROUTINES=1
KAFKA_WRITE_EXAMPLES=0 # при 1 (true) генерирует данные для записи в Kafka
KAFKA_WRITER_GOROUTINES=1

MEMORY_HANDLER_GOROUTINES=1

MEMORY_RESTORE_DATA=1
MEMORY_CLEANUP_MINUTES_INTERVAL=10
MEMORY_ORDERS_LIMIT=100000 # исходя из максимального размера сообщения 20kb, 100 тысяч сообщений это не более 2гб памяти

SERVER_PORT=3000
SERVER_GET_ORDER_BY_ID=1 # при 0 (false) заказ можно получить по его order_uid
```

## Бенчмарки

С увеличением количества горутин для чтения из Kafka (`KAFKA_READER_GOROUTINES`) и обработки сообщений (`MEMORY_HANDLER_GOROUTINES`) скорость работы сервиса увеличивается:

- Значения по умолчанию (Readers: 1, Handlers: 1): ~ 600 сообщений в секунду.
- Несколько обработчиков (Readers: 1, Handlers: 20): ~ 2500 сообщений в секунду.
- Параллельное чтение (Readers: 7, Handlers: 20): ~ 2750 сообщений в секунду.

Результаты теста HTTP-сервера через WRK (стандартный запрос):

```
Running 30s test @ http://127.0.0.1:3000/orders?id=1
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    49.74ms   48.24ms 224.63ms   77.35%
    Req/Sec     0.89k   474.99     3.30k    56.54%
  310434 requests in 30.06s, 551.55MB read
Requests/sec:  10327.50
Transfer/sec:     18.35MB
```

Результаты теста HTTP-сервера через Vegeta (стандартный запрос):

```
echo "GET http://localhost:3000/orders?id=1" | vegeta attack -duration=30s | tee results.bin | vegeta report
Requests      [total, rate, throughput]         1500, 50.03, 50.03
Duration      [total, attack, wait]             29.981s, 29.981s, 628.023µs
Latencies     [min, mean, 50, 90, 95, 99, max]  320.79µs, 620.393µs, 640.242µs, 780.167µs, 819.74µs, 1.078ms, 4.095ms
Bytes In      [total, mean]                     1621500, 1081.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:1500  
Error Set:
```

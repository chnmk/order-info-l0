package memory

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

// Один раз в указанный промежуток времени блокирует хранилище и удаляет из него ненужные данные.
func (m *MemStore) ClearData(ctx context.Context) {
	defer cfg.ExitWg.Done()

	for {
		select {

		case <-ctx.Done():
			slog.Info("stopping data cleaner...")
			return

		default:
			// time.Sleep(time.Duration(cfg.CleanupInterval) * time.Minute)
			time.Sleep(10 * time.Second)

			slog.Info(
				"data cleaner started...",
				"len", len(m.orders),
			)

			i := 0
			m.mu.Lock()

			// Если превышен лимит заказов, оставляет только самые новые.
			if len(m.orders) > cfg.OrdersLimit {
				removing := len(m.orders) - cfg.OrdersLimit
				m.orders = m.orders[removing-1:]
			}

			for _, order := range m.orders {

				// Будем считать устаревшими заказы, сделанные более 14 дней назад.
				expDate := time.Now().AddDate(0, 0, -14)
				dateConv, err := time.Parse("2006-01-02T15:04:05Z", order.Date_created)

				if dateConv.Before(expDate) || err != nil {
					// Если заказ попадает под условие удаления, нужно сообщить, что заказ будет удалён.
					fmt.Println("ERRRRRRRRRR: " + err.Error())
					if err != nil {
						slog.Error(
							"invalid date string for order, deleting...",
							"id", order.ID,
						)
					} else {
						slog.Error(
							"order expired, deleting...",
							"id", order.ID,
						)
					}
				} else {
					// Если заказ не попадает под условие удаления, он записывается назад в тот же массив.
					m.orders[i] = order
					i++
				}

			}

			// Оставляет массив того размера, сколько было записано заказов, не попавших под удаление.
			m.orders = m.orders[:i]

			m.mu.Unlock()

			slog.Info(
				"data cleaner finished working",
				"len", len(m.orders),
			)
		}
	}

}

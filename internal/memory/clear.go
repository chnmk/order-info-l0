package memory

import (
	"log/slog"
	"time"

	cfg "github.com/chnmk/order-info-l0/internal/config"
)

// Один раз в указанный промежуток времени блокирует хранилище и удаляет из него ненужные данные.
func (m *MemStore) ClearData() {
	defer cfg.ExitWg.Done()

	for {
		select {

		case <-cfg.ExitCtx.Done():
			slog.Info("stopping data cleaner...")
			return

		default:
			// TODO: добавить переменную env - раз во сколько секунд запускать сборщик.
			time.Sleep(1 * time.Minute)

			slog.Info(
				"data cleaner started...",
				"len", len(m.orders),
			)

			i := 0
			m.mu.Lock()

			// Если превышен лимит заказов, оставляет только самые новые.
			// TODO: по переменной окружения
			if len(m.orders) > 15 {
				removing := len(m.orders) - 15
				m.orders = m.orders[removing-1:]
			}

			for _, order := range m.orders {

				// TODO: добавить в env число дней
				expDate := time.Now().AddDate(0, 0, -14)
				dateConv, err := time.Parse(time.UnixDate, order.Date_created)

				if dateConv.Before(expDate) || err != nil {
					// Если заказ попадает под условие удаления, нужно сообщить, что заказ будет удалён.
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
package memory

import "github.com/chnmk/order-info-l0/internal/models"

// Проверяет что нужные поля не пустые и соответствуют нашим требованиям.
//
// Пока что нам точно нужны те данные, которые выводятся в веб-интерфейсе.
func ValidateMsg(order models.Order) bool {
	if order.Order_uid == "" ||
		order.Delivery.Name == "" ||
		order.Delivery.City == "" ||
		order.Delivery.Address == "" ||
		order.Delivery.Phone == "" ||
		len(order.Items) < 1 {

		return false
	}

	for _, i := range order.Items {
		if i.Chrt_id == 0 ||
			i.Name == "" ||
			i.Total_price == 0 {
			return false
		}
	}

	return true
}

package models

type Order struct {
	Order_uid          string   `json:"order_uid"`
	Track_number       string   `json:"track_number"`
	Entry              string   `json:"entry"`
	Delivery           Delivery `json:"delivery"`
	Payment            Payment  `json:"payment"`
	Items              []Item   `json:"items"`
	Locale             string   `json:"locale"`
	Internal_signature string   `json:"internal_signature"`
	Customer_id        string   `json:"customer_id"`
	Delivery_service   string   `json:"delivery_service"`
	Shardkey           string   `json:"shardkey"`
	Sm_id              int32    `json:"sm_id"`
	Date_created       string   `json:"date_created"`
	Oof_shard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        int32  `json:"amount"`
	Payment_dt    int32  `json:"payment_dt"`
	Bank          string `json:"bank"`
	Delivery_cost int32  `json:"delivery_cost"`
	Goods_total   int32  `json:"goods_total"`
	Custom_fee    int32  `json:"custom_fee"`
}

type Item struct {
	Chrt_id      int32  `json:"chrt_id"`
	Track_number string `json:"track_number"`
	Price        int32  `json:"price"`
	Rid          string `json:"rid"`
	Name         string `json:"name"`
	Sale         int32  `json:"sale"`
	Size         string `json:"size"`
	Total_price  int32  `json:"total_price"`
	Nm_id        int32  `json:"nm_id"`
	Brand        string `json:"brand"`
	Status       int32  `json:"status"`
}

package models

type Order struct {
	Order_uid          string   `json:"order_uid"`
	Track_number       string   `json:"track_number"`
	Entry              string   `json:"entry"`
	Delivery           Delivery `json:"delivery"`
	Payment            Payment  `json:"payment"`
	Items              []Item   `json:"items" fakesize:"1,10"`
	Locale             string   `json:"locale" fake:"{languageabbreviation}"`
	Internal_signature string   `json:"internal_signature"`
	Customer_id        string   `json:"customer_id"`
	Delivery_service   string   `json:"delivery_service" fake:"{word}"`
	Shardkey           string   `json:"shardkey" fake:"{digit}"`
	Sm_id              int32    `json:"sm_id" fake:"{number: 1, 1000}"`
	Date_created       string   `json:"date_created" fake:"{wbtime}"`
	Oof_shard          string   `json:"oof_shard" fake:"{digit}"`
}

type Delivery struct {
	Name    string `json:"name" fake:"{name}"`
	Phone   string `json:"phone" fake:"{phone}"`
	Zip     string `json:"zip" fake:"{zip}"`
	City    string `json:"city" fake:"{city}"`
	Address string `json:"address" fake:"{streetname}"`
	Region  string `json:"region" fake:"{state}"`
	Email   string `json:"email" fake:"{email}"`
}

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency" fake:"{currencyshort}"`
	Provider      string `json:"provider" fake:"{word}"`
	Amount        int32  `json:"amount" fake:"{number: 1, 50000}"`
	Payment_dt    int32  `json:"payment_dt"`
	Bank          string `json:"bank" fake:"{word}"`
	Delivery_cost int32  `json:"delivery_cost" fake:"{number: 1, 5000}"`
	Goods_total   int32  `json:"goods_total" fake:"{number: 1, 100000}"`
	Custom_fee    int32  `json:"custom_fee" fake:"{number: 0, 1000}"`
}

type Item struct {
	Chrt_id      int32  `json:"chrt_id" fake:"{number: 1, 10000000}"`
	Track_number string `json:"track_number"`
	Price        int32  `json:"price" fake:"{number: 1, 5000}"`
	Rid          string `json:"rid"`
	Name         string `json:"name" fake:"{productname}"`
	Sale         int32  `json:"sale" fake:"{number: 0, 1000}"`
	Size         string `json:"size" fake:"{digit}"`
	Total_price  int32  `json:"total_price" fake:"{number: 10, 5000}"`
	Nm_id        int32  `json:"nm_id" fake:"{number: 1, 10000000}"`
	Brand        string `json:"brand" fake:"{word}"`
	Status       int32  `json:"status" fake:"{number: 100, 999}"`
}

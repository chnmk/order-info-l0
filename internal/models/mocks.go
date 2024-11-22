package models

type MockEmptyOrder struct {
	Order_uid string `json:"order_uid"`
}

type MockEmptyOrderWithNumber struct {
	Order_uid    string `json:"order_uid"`
	Track_number string `json:"track_number"`
}

type MockNotOrder struct {
	Order_fake string `json:"order_fake"`
}

type MockNotOrderWithNumber struct {
	Order_fake   string `json:"order_fake"`
	Track_number string `json:"track_number"`
}

type MockOrderIntIsString struct {
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
	Sm_id              string   `json:"sm_id"`
	Date_created       string   `json:"date_created" fake:"{wbdate}"`
	Oof_shard          string   `json:"oof_shard" fake:"{digit}"`
}

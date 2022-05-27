package main

type Data struct {
	OrderUid          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json: "payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required,required"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	ShardKey          string   `json:"shardkey"`
	Sm_id             int      `json:"sm_id" validate:"required"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestId    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
	PaymentDt    int     `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   uint    `json:"goods_total" validate:"required"`
	CustomFee    float64 `json:"custom_fee"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required"`
}

type Item struct {
	ChrtId      int     `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Rid         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        float64 `json:"sale"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float64 `json:total_price validate:"required"`
	Nm_id       int     `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand"`
	Status      int  `json:"status"`
}

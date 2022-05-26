package main

type Data struct {
	Order_uid          string
	Track_number       string
	Entry              string
	Delivery           Delivery
	Payment            Payment
	Items              []Item
	Locale             string
	Internal_signature string
	Customer_id        string
	Delivery_service   string
	ShardKey           string
	Sm_id              int
	Date_created       string
	Oof_shard          string
}

type Payment struct {
	Transaction   string
	Request_id    string
	Currency      string
	Provider      string
	Amount        float64
	Payment_dt    int
	Bank          string
	Delivery_cost float64
	Goods_total   uint
	Custom_fee    float64
}

type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type Item struct {
	Chrt_id      int
	Track_number string
	Price        float64
	Rid          string
	Name         string
	Sale         float64
	Size         string
	Total_price  float64
	Nm_id        int
	Brand        string
}

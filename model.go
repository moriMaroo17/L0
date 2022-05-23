package main

type Data struct {
	Payment Payment
}

type Payment struct {
	Transaction   string
	Request_id    string
	Currency      string
	Provider      string
	Amount        uint
	Payment_dt    int
	Bank          string
	Delivery_cost uint
	Goods_total   uint
	Custom_fee    int
}

package main

import (
	"github.com/nats-io/stan.go"
)

const jsonData = `
{
  "order_uid": "993somemoreuid666",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Ghoul Ghoulevich",
    "phone": "+9936661007",
    "zip": "123456",
    "city": "Tokyo Shibuya",
    "address": "Pushkina Dom Kolotushkina",
    "region": "North Africa",
    "email": "zxc234@gmail.com"
  },
  "payment": {
    "transaction": "993somemoreuid666",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 666222,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 10000,
      "rid": "zxc12190asdzxce0btest",
      "name": "roxy migurdia full size figure",
      "sale": 30,
      "size": "0",
      "total_price": 9999,
      "nm_id": 2389212,
      "brand": "noname",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2023-11-03T06:22:19Z",
  "oof_shard": "1"
}
`

func main() {
	sc, _ := stan.Connect("test-cluster", "publisher-id", stan.NatsURL("nats://localhost:4222"))
	defer sc.Close()
	sc.Publish("test-subject", []byte(jsonData))
}

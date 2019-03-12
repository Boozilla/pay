package pay

import "time"

type product struct {
	productId  string
	quantity   int
	orderId    string
	purchaseAt time.Time
}

func (p product) ProductId() string {
	return p.productId
}

func (p product) Quantity() int {
	return p.quantity
}

func (p product) OrderId() string {
	return p.orderId
}

func (p product) PurchaseAt() time.Time {
	return p.purchaseAt
}

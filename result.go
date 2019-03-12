package pay

type result struct {
	purchaseToken string
	status        status
	isTest        bool
	products      []product
}

func (r result) PurchaseToken() string {
	return r.purchaseToken
}

func (r result) Status() status {
	return r.status
}

func (r result) IsTest() bool {
	return r.isTest
}

func (r result) Products() []product {
	return r.products
}

package pay

import (
	"runtime"
	"testing"
	"time"
)

func TestGooglePlay(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	try := 10
	done := make(chan bool)

	for i := 0; i < try; i++ {
		go exec(t, done, `{"orderId":"","packageName":"","productId":"","purchaseTime":0,"purchaseState":0,"developerPayload":"","purchaseToken":""}`,
			GooglePlayConfig{
				JsonKeyPath: "key.json",
				Timeout:     10 * time.Second,
			})
	}

	for i := 0; i < try; i++ {
		<-done
	}
}

func TestAppStore(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	try := 10
	done := make(chan bool)

	for i := 0; i < try; i++ {
		go exec(t, done, "MIITxQYJKoZIhvcNAQcCoIITtjCCE7ICAQExCzAJBgUrDgMCGgUAMIIDZgYJKoZIhvc ... ",
			AppStoreConfig{
				Sandbox: false,
				Timeout: 10 * time.Second,
			})
	}

	for i := 0; i < try; i++ {
		<-done
	}
}

func exec(t *testing.T, done chan<- bool, receipt string, conf interface{}) {
	validator, err := NewValidator(conf)

	if err != nil {
		panic(err)
	}

	ret, err := validator.Do(receipt)

	if err != nil {
		println(err.Error())

		done <- true
		return
	}

	t.Log("PurchaseToken:", ret.PurchaseToken())
	t.Log("Status:", ret.Status().String())
	t.Log("IsTest:", ret.IsTest())

	for _, p := range ret.Products() {
		t.Log("ProductId:", p.ProductId())
		t.Log("Quantity:", p.Quantity())
		t.Log("OrderId:", p.OrderId())
		t.Log("PurchaseAt:", p.PurchaseAt().String())
	}

	if done != nil {
		done <- true
	}
}

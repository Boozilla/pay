package pay

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/androidpublisher/v3"
	"io/ioutil"
	"sync"
	"time"
)

type GooglePlayConfig struct {
	JsonKeyPath string
	Retry       int
	Timeout     time.Duration
}

type googlePlayReceipt struct {
	PackageName string `json:"packageName"`
	ProductId   string `json:"productId"`
	Token       string `json:"purchaseToken"`
}

var gpConfig *jwt.Config
var gpOnce sync.Once

func (c GooglePlayConfig) service() *androidpublisher.Service {
	gpOnce.Do(func() {
		key, err := ioutil.ReadFile(c.JsonKeyPath)
		if err != nil {
			panic(err)
		}

		conf, err := google.JWTConfigFromJSON(key, androidpublisher.AndroidpublisherScope)
		if err != nil {
			panic(err)
		}

		gpConfig = conf
	})

	client := gpConfig.Client(context.TODO())
	client.Timeout = c.Timeout

	androidPublisherService, err := androidpublisher.New(client)

	if err != nil {
		panic(err)
	}

	return androidPublisherService
}

func (c GooglePlayConfig) requestHandler(data string) (*result, error) {
	status := Valid
	receipt, err := unmarshalReceipt(data)

	if err != nil {
		return nil, err
	}

	publisher := c.service()
	req := publisher.Purchases.Products.Get(
		receipt.PackageName,
		receipt.ProductId,
		receipt.Token)

	ret, err := req.Do()

	if err != nil {
		return nil, err
	}

	if ret.ConsumptionState == 1 /* Consumed */ {
		status = ConsumedReceipt
	} else if ret.PurchaseState == 1 /* Canceled */ {
		status = CanceledReceipt
	}

	return &result{
		purchaseToken: receipt.Token,
		status:        status,
		isTest:        *ret.PurchaseType == 0,
		products: []product{
			{receipt.ProductId, 1, ret.OrderId, time.Unix(0, ret.PurchaseTimeMillis*int64(time.Millisecond))},
		},
	}, nil
}

func unmarshalReceipt(data string) (*googlePlayReceipt, error) {
	receipt := googlePlayReceipt{}
	err := json.Unmarshal([]byte(data), &receipt)

	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

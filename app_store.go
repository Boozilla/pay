package pay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type appStoreMaintenanceError struct {
	receiptData string
}

func (e *appStoreMaintenanceError) Error() string {
	return fmt.Sprintf("pay: AppStore: receipt=%s: %s",
		e.receiptData,
		"The App store server is temporarily unavailable.")
}

type AppStoreConfig struct {
	Sandbox bool
	Retry   int
	Timeout time.Duration
}

type appStoreBaseURL struct {
	prod    *url.URL
	sandbox *url.URL
}

type appStoreRequest struct {
	ReceiptData           string `json:"receipt-data"`
	Password              string `json:"password"`
	ExcludeOldTransaction bool   `json:"exclude-old-transactions"`
}

type appStoreResponse struct {
	Status      int    `json:"status"`
	Environment string `json:"environment"`

	Receipt struct {
		BundleId   string `json:"bundle_id"`
		AppVersion string `json:"application_version"`
		PurchaseAt string `json:"original_purchase_date_ms"`
		RequestAt  string `json:"request_date_ms"`
		CreatedAt  string `json:"receipt_creation_date_ms"`

		Purchases []struct {
			Quantity              string `json:"quantity"`
			ProductId             string `json:"product_id"`
			TransactionId         string `json:"transaction_id"`
			OriginalTransactionId string `json:"original_transaction_id"`
			PurchaseAt            string `json:"purchase_date_ms"`
			OriginalPurchaseAt    string `json:"original_purchase_date_ms"`
			IsTrialPeriod         string `json:"is_trial_period"`
		} `json:"in_app"`
	} `json:"receipt"`
}

var asBaseURL *appStoreBaseURL
var asOnce sync.Once

func (c AppStoreConfig) baseURL() *appStoreBaseURL {
	asOnce.Do(func() {
		prodURL, err := url.Parse("https://buy.itunes.apple.com/verifyReceipt")
		if err != nil {
			panic(err)
		}

		sandboxURL, err := url.Parse("https://sandbox.itunes.apple.com/verifyReceipt")
		if err != nil {
			panic(err)
		}

		asBaseURL = &appStoreBaseURL{
			prod:    prodURL,
			sandbox: sandboxURL,
		}
	})

	return asBaseURL
}

func marshalRequest(req *appStoreRequest) (*[]byte, error) {
	json, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	return &json, nil
}

func unmarshalResponse(respBody *[]byte) (*appStoreResponse, error) {
	resp := &appStoreResponse{}
	err := json.Unmarshal(*respBody, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c AppStoreConfig) send(url *url.URL, data *[]byte) (*appStoreResponse, error) {
	client := http.DefaultClient
	client.Timeout = c.Timeout

	resp, err := client.Post(url.String(), "application/json", bytes.NewReader(*data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj, err := unmarshalResponse(&respBody)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func (c AppStoreConfig) requestHandler(data string) (*result, error) {
	req, err := marshalRequest(&appStoreRequest{
		ReceiptData:           data,
		Password:              "",
		ExcludeOldTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	baseURL := c.baseURL()

	var endpoint *url.URL
	if c.Sandbox {
		endpoint = baseURL.sandbox
	} else {
		endpoint = baseURL.prod
	}

	resp, err := c.send(endpoint, req)
	if err != nil {
		return nil, err
	}

	if resp.Status == 21007 {
		/**
		 * This receipt is from the test environment,
		 * but it was sent to the production environment for verification.
		 * Send it to the test environment instead.
		 */
		resp, err = c.send(baseURL.sandbox, req)
		if err != nil {
			return nil, err
		}
	}

	status := Valid

	var productList []product
	switch resp.Status {
	case 0:
		productList = make([]product, len(resp.Receipt.Purchases))

		for i, e := range resp.Receipt.Purchases {
			var purchaseAt time.Time
			purchaseMs, err := strconv.ParseInt(e.PurchaseAt, 0, 64)

			if err != nil {
				purchaseAt = time.Now()
			}

			purchaseAt = time.Unix(0, purchaseMs*int64(time.Millisecond))

			qty, _ := strconv.ParseInt(e.Quantity, 0, 32)

			productList[i] = product{
				productId:  e.ProductId,
				quantity:   int(qty),
				orderId:    e.TransactionId,
				purchaseAt: purchaseAt,
			}
		}
		break
	case 21005:
		return nil, &appStoreMaintenanceError{
			data,
		}
		break
	default:
		status = Invalid
		break
	}

	return &result{
		purchaseToken: data,
		status:        status,
		isTest:        resp.Environment == "Sandbox",
		products:      productList,
	}, nil
}

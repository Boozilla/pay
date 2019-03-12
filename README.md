### Usage example (Google Play)

```go
package main

import "github.com/Boozilla/pay"

func main() {
	conf := pay.GooglePlayConfig{
		JsonKeyPath: "<GOOGLE_API_JSON_KEY_PATH>",
		Timeout:     10 * time.Second,
	}

	done := make(chan bool)

	go func() {
		validator, err := pay.NewValidator(conf)
		if err != nil {
			// Validator initialize error handling
		}

		ret, err := validator.Do(`{"orderId":"","packageName":"","productId":"","purchaseTime":0,"purchaseState":0,"developerPayload":"","purchaseToken":""}`)
		if err != nil {
			// Validation request error handling
		}

		println(ret) // Validation request result

		done <- true
	}()
	
	<-done
}
```

### Usage example (App Store)

```go
package main

import "github.com/Boozilla/pay"

func main() {
	conf := pay.AppStoreConfig{
		Sandbox: false,
		Timeout: 10 * time.Second,
	}

	done := make(chan bool)

	go func() {
		validator, err := pay.NewValidator(conf)
		if err != nil {
			// Validator initialize error handling
		}

		ret, err := validator.Do("MIITxQYJKoZIhvcNAQcCoIITtjCCE7ICAQExCzAJBgUrDgMCGgUAMIIDZgYJKoZIhvc ... ")
		if err != nil {
			// Validation request error handling
		}

		println(ret) // Validation request result

		done <- true
	}()
	
	<-done
}
```
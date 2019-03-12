package pay

type requestHandler = func(receipt string) (*result, error)

type validator struct {
	conf    interface{}
	request requestHandler
}

func (v validator) Do(data string) (*result, error) {
	return v.request(data)
}

package pay

import "fmt"

type unknownType struct {
	typeName string
}

func (e *unknownType) Error() string {
	return fmt.Sprintf("pay: type=%s: %s",
		e.typeName,
		"Unknown receipt verification channel config")
}

func NewValidator(conf interface{}) (*validator, error) {
	reqHandler, err := resolveRequestHandler(conf)

	if err != nil {
		return nil, err
	}

	return &validator{
		request: reqHandler,
		conf:    conf,
	}, nil
}

func resolveRequestHandler(conf interface{}) (requestHandler, error) {
	switch t := conf.(type) {
	case GooglePlayConfig:
		return conf.(GooglePlayConfig).requestHandler, nil
	case AppStoreConfig:
		return conf.(AppStoreConfig).requestHandler, nil
	default:
		return nil, &unknownType{
			fmt.Sprintf("%T", t),
		}
	}
}

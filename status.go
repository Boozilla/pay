package pay

type status int

const (
	Valid status = 0 + iota
	Invalid
	ConsumedReceipt
	CanceledReceipt
)

var messages = [...]string{
	"Valid",
	"Invalid",
	"Consumed Receipt",
}

func (s status) String() string {
	return messages[s]
}

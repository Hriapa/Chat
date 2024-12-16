package protocol

type MessageDiscriminator uint8

const (
	ControlMessageTitle MessageDiscriminator = iota
	DataMessageTitle
	AckMessageTitle
	ErrorMessageTitle
)

type Title uint8

// 0 - reserved
const (
	userIdTitle Title = iota + 1
	roomIdTitle
	nameTitle
	surnameTitle
	familyNameTitle
	birthDateTitle
	userNameTitle
	listElementsTitle
	numberListElementsTitle
	lastMessageNumberTitle
	userDataTitle
	timeTitle
)

// for Fragmentation // 10xx_xxxx - first, 00xx_xxxx - middle, 01xx_xxxx - last

type Fragment uint8

const (
	FirstFragment  Fragment = 0x80
	MiidleFragment Fragment = 0x00
	LastFragment   Fragment = 0x40
)

type Message interface {
	Clear()
	Code() []byte
	Decode([]byte) error
}

func Cleaner(m Message) {
	m.Clear()
}

func Coder(m Message) []byte {
	return m.Code()
}

func Decoder(m Message, in []byte) error {
	return m.Decode(in)
}

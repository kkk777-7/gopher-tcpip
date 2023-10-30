package ethernet

type Address [6]byte

type Type uint16

type Header struct {
	Src  Address
	Dst  Address
	Type Type
}

type Frame struct {
	Header  Header
	Payload []byte
}

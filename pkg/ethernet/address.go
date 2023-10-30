package ethernet

import "fmt"

func ToStringFromByte(data []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", data[0], data[1], data[2], data[3], data[4], data[5])
}

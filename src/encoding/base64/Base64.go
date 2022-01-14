package base64

import "encoding/base64"

func Decode(data string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	} else {
		return decoded
	}
}

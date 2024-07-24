package storage

import "encoding/base64"

func DecodeString(data string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func DecodeBytes(data []byte) (string, error) {
	s := base64.StdEncoding.EncodeToString(data)
	return s, nil
}

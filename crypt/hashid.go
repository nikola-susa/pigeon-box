package crypt

import (
	"encoding/hex"
	"github.com/speps/go-hashids/v2"
)

func HashIDEncodeString(message string, salt string, length int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length
	h, _ := hashids.NewWithData(hd)

	userIDHex := hex.EncodeToString([]byte(message))

	res, err := h.EncodeHex(userIDHex)
	if err != nil {
		return "", err
	}

	return res, nil
}

func HashIDDecodeString(message string, salt string, length int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length
	h, _ := hashids.NewWithData(hd)

	resHex, err := h.DecodeHex(message)
	if err != nil {
		return "", err
	}

	res, err := hex.DecodeString(resHex)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func HashIDEncodeInt(message int, salt string, length int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length
	h, _ := hashids.NewWithData(hd)

	res, err := h.Encode([]int{message})
	if err != nil {
		return "", err
	}

	return res, nil
}

func HashIDDecodeInt(message string, salt string, length int) (int, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = length
	h, _ := hashids.NewWithData(hd)

	res, err := h.DecodeWithError(message)
	if err != nil {
		return 0, err
	}

	return res[0], nil
}

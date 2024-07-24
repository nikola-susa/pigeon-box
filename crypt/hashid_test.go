package crypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashIDEncodeString(t *testing.T) {
	message := "message"
	salt := "salt"
	length := 10

	res, err := HashIDEncodeString(message, salt, length)
	assert.NoError(t, err)

	assert.NotEmpty(t, res)
}

func TestHashIDDecodeString(t *testing.T) {
	message := "message"
	salt := "salt"
	length := 10

	encoded, err := HashIDEncodeString(message, salt, length)
	assert.NoError(t, err)

	res, err := HashIDDecodeString(encoded, salt, length)
	assert.NoError(t, err)

	assert.Equal(t, message, res)
}

func TestHashIDDecodeStringWithString(t *testing.T) {
	message := "message"
	salt := "salt"
	length := 10

	_, err := HashIDDecodeString(message, salt, length)
	assert.Error(t, err)
}

func TestHashIDEncodeInt(t *testing.T) {
	message := 123
	salt := "salt"
	length := 10

	res, err := HashIDEncodeInt(message, salt, length)
	assert.NoError(t, err)

	assert.NotEmpty(t, res)
}

func TestHashIDDecodeInt(t *testing.T) {
	message := 123
	salt := "salt"
	length := 10

	encoded, err := HashIDEncodeInt(message, salt, length)
	assert.NoError(t, err)

	res, err := HashIDDecodeInt(encoded, salt, length)
	assert.NoError(t, err)

	assert.Equal(t, message, res)
}

func TestHashIDDecodeIntWithString(t *testing.T) {
	message := "message"
	salt := "salt"
	length := 10

	_, err := HashIDDecodeInt(message, salt, length)
	assert.Error(t, err)
}

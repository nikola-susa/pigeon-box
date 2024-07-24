package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeStringWithValidDataReturnsDecodedBytes(t *testing.T) {
	data := "SGVsbG8sIFdvcmxkIQ=="
	expectedBytes := []byte("Hello, World!")

	decodedBytes, err := DecodeString(data)

	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, decodedBytes)
}

func TestDecodeStringWithInvalidDataReturnsError(t *testing.T) {
	data := "invalid data"

	_, err := DecodeString(data)

	assert.Error(t, err)
}

func TestDecodeBytesWithValidDataReturnsEncodedString(t *testing.T) {
	data := []byte("Hello, World!")
	expectedString := "SGVsbG8sIFdvcmxkIQ=="

	encodedString, err := DecodeBytes(data)

	assert.NoError(t, err)
	assert.Equal(t, expectedString, encodedString)
}

func TestDecodeBytesWithEmptyDataReturnsEmptyString(t *testing.T) {
	data := []byte("")

	encodedString, err := DecodeBytes(data)

	assert.NoError(t, err)
	assert.Equal(t, "", encodedString)
}

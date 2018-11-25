package bbclib

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"time"
)


func GetIdentifier(seed string, length int) []byte {
	digest := sha256.Sum256([]byte(seed))
	return digest[:length]
}

func GetIdentifierWithTimestamp(seed string, length int) []byte {
	digest := sha256.Sum256([]byte(seed+time.Now().String()))
	return digest[:length]
}

func GetRandomValue(length int) []byte {
	val := make([]byte, length)
	_, err := rand.Read(val)
	if err != nil {
		for i := range val {
			val[i] = 0x00
		}
	}
	return val
}

func Put2byte(buf *bytes.Buffer, val uint16) error{
	return binary.Write(buf, binary.LittleEndian, val)
}

func Get2byte(buf *bytes.Buffer) (uint16, error) {
	var val uint16
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		return 0, err
	}
	return val, nil
}

func Put4byte(buf *bytes.Buffer, val uint32) error{
	return binary.Write(buf, binary.LittleEndian, val)
}

func Get4byte(buf *bytes.Buffer) (uint32, error) {
	var val uint32
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		return 0, err
	}
	return val, nil
}

func Put8byte(buf *bytes.Buffer, val int64) error{
	return binary.Write(buf, binary.LittleEndian, val)
}

func Get8byte(buf *bytes.Buffer) (int64, error) {
	var val int64
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		return 0, err
	}
	return val, nil
}

func PutBigInt(buf *bytes.Buffer, val *[]byte, length int) {
	Put2byte(buf, uint16(length))
	binary.Write(buf, binary.LittleEndian, val)
}

func GetBigInt(buf *bytes.Buffer) ([]byte, error) {
	length, err := Get2byte(buf)
	if err != nil {
		return nil, err
	}
	return GetBytes(buf, int(length))
}

func GetBytes(buf *bytes.Buffer, length int) ([]byte, error) {
	val := make([]byte, length)
	if err := binary.Read(buf, binary.LittleEndian, val); err != nil {
		return nil, err
	}
	return val, nil
}

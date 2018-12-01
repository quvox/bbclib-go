/*
Copyright (c) 2018 Zettant Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 */

package bbclib

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"time"
)


// Get random byte data with specified length (seed string ais used)
func GetIdentifier(seed string, length int) []byte {
	digest := sha256.Sum256([]byte(seed))
	return digest[:length]
}

// Get random byte data with specified length (seed string and timestamp are used)
func GetIdentifierWithTimestamp(seed string, length int) []byte {
	digest := sha256.Sum256([]byte(seed+time.Now().String()))
	return digest[:length]
}

// Get random byte data with specified length
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

// Put uint16 in the buffer for packing
func Put2byte(buf *bytes.Buffer, val uint16) error{
	return binary.Write(buf, binary.LittleEndian, val)
}

// Get uint16 value from the buffer
func Get2byte(buf *bytes.Buffer) (uint16, error) {
	var val uint16
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		return 0, err
	}
	return val, nil
}

// Put uint32 in the buffer for packing
func Put4byte(buf *bytes.Buffer, val uint32) error{
	return binary.Write(buf, binary.LittleEndian, val)
}

// Get uint32 value from the buffer
func Get4byte(buf *bytes.Buffer) (uint32, error) {
	var val uint32
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		return 0, err
	}
	return val, nil
}

// Put int64 in the buffer for packing
func Put8byte(buf *bytes.Buffer, val int64) error{
	return binary.Write(buf, binary.LittleEndian, val)
}

// Get int64 value from the buffer
func Get8byte(buf *bytes.Buffer) (int64, error) {
	var val int64
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		return 0, err
	}
	return val, nil
}

// Put ID data in the buffer for packing
func PutBigInt(buf *bytes.Buffer, val *[]byte, length int) {
	Put2byte(buf, uint16(length))
	binary.Write(buf, binary.LittleEndian, val)
}

// Get ID data from the buffer
func GetBigInt(buf *bytes.Buffer) ([]byte, error) {
	length, err := Get2byte(buf)
	if err != nil {
		return nil, err
	}
	return GetBytes(buf, int(length))
}

// Get binary data with specified length from the buffer
func GetBytes(buf *bytes.Buffer, length int) ([]byte, error) {
	val := make([]byte, length)
	if err := binary.Read(buf, binary.LittleEndian, val); err != nil {
		return nil, err
	}
	return val, nil
}
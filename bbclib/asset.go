package bbclib

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/ugorji/go/codec"
)

type (
	BBcAsset struct {
		IdLength		int
		digestCalculating bool
		AssetId			[]byte
		UserId			[]byte
		Nonce			[]byte
		AssetFileSize 	uint32
		AssetFile 		[]byte
		AssetFileDigest	[]byte
		AssetBodyType	uint16
		AssetBodySize 	uint16
		AssetBody		[]byte
	}
)

var (
	mh codec.MsgpackHandle
)


func encodeMessagePack(values interface{}) ([]byte, error) {
	buf := make([]byte, 0, 65536)

	err := codec.NewEncoderBytes(&buf, &mh).Encode(values)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func decodeMessagePack(buf []byte) (interface{}, error) {
	var values interface{}

	err := codec.NewDecoderBytes(buf, &mh).Decode(&values)
	if err != nil {
		return nil, err
	}
	return values, nil
}


func (p *BBcAsset) Stringer() string {
	ret :=  "  Asset:\n"
	ret += fmt.Sprintf("     asset_id: %x\n", p.AssetId)
	ret += fmt.Sprintf("     user_id: %x\n", p.UserId)
	ret += fmt.Sprintf("     nonce: %x\n", p.Nonce)
	ret += fmt.Sprintf("     file_size: %d\n", p.AssetFileSize)
	if p.AssetFileDigest != nil {
		ret += fmt.Sprintf("     file_digest: %x\n", p.AssetFileDigest)
	} else {
		ret += "     file_digest: None\n"
	}
	ret += fmt.Sprintf("     body_size: %d\n", p.AssetBodySize)
	ret += fmt.Sprintf("     body: %v\n", p.AssetBody)
	return ret
}


func (p *BBcAsset) Add(userId *[]byte) {
	if userId != nil {
		if p.IdLength == 0 {
			p.IdLength = len(*userId)
		}
		p.UserId = make([]byte, p.IdLength)
		copy(p.UserId, (*userId)[:p.IdLength])
	}
	p.Nonce = GetRandomValue(p.IdLength)
}

func (p *BBcAsset) AddFile(fileContent *[]byte) {
	p.AssetFileSize = uint32(binary.Size(fileContent))
	digest := sha256.Sum256(*fileContent)
	p.AssetFileDigest = digest[:]
}

func (p *BBcAsset) AddBodyString(bodyContent string) {
	p.AssetBodyType = 0
	p.AssetBody = []byte(bodyContent)
	p.AssetBodySize = uint16(len(bodyContent))
}

func (p *BBcAsset) AddBodyObject(bodyContent interface{}) error {
	p.AssetBodyType = 1
	var err error
	p.AssetBody, err = encodeMessagePack(bodyContent)
	if err != nil {
		return err
	}
	p.AssetBodySize = uint16(len(p.AssetBody))
	return nil
}

func (p *BBcAsset) GetBodyObject() (interface{}, error) {
	if p.AssetBodyType != 1 {
		return nil, nil
	}
	return decodeMessagePack(p.AssetBody)
}

func (p *BBcAsset) Digest() []byte {
	p.digestCalculating = true
	asset, err := p.Pack()
	if err != nil {
		p.digestCalculating = false
		return nil
	}

	digest := sha256.Sum256(asset)
	if p.AssetId == nil {
		p.AssetId = make([]byte, p.IdLength)
	}
	p.AssetId = digest[:p.IdLength]
	p.digestCalculating = false
	return digest[:]
}

func (p *BBcAsset) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	if ! p.digestCalculating {
		if p.AssetId == nil {
			p.Digest()
		}
		PutBigInt(buf, &p.AssetId, p.IdLength)
	}
	PutBigInt(buf, &p.UserId, p.IdLength)
	PutBigInt(buf, &p.Nonce, len(p.Nonce))
	Put4byte(buf, p.AssetFileSize)
	if p.AssetFileSize > 0 {
		if err := binary.Write(buf, binary.LittleEndian, p.AssetFileDigest); err != nil {
			return nil, err
		}
	}

	Put2byte(buf, p.AssetBodyType)
	Put2byte(buf, p.AssetBodySize)
	if p.AssetBodySize > 0 {
		if err := binary.Write(buf, binary.LittleEndian, p.AssetBody); err != nil {
			return nil, err
		}
	}
	fmt.Printf("packed: %x\n", buf.Bytes())
	return buf.Bytes(), nil
}


func (p *BBcAsset) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.UserId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.Nonce, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.AssetFileSize, err = Get4byte(buf)
	if err != nil {
		return err
	}
	if p.AssetFileSize > 0 {
		p.AssetFileDigest, err = GetBytes(buf, 32)
		if err != nil {
			return err
		}
	}

	p.AssetBodyType, err = Get2byte(buf)
	if err != nil {
		return err
	}
	p.AssetBodySize, err = Get2byte(buf)
	if err != nil {
		return err
	}
	p.AssetBody, err = GetBytes(buf, int(p.AssetBodySize))

	return nil
}

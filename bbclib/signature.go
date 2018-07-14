package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcSignature struct {
		Format_type		int			`bson:"-" json:"-"`
		Key_type 		int			`bson:"key_type" json:"key_type"`
		Signature 		[]byte		`bson:"signature" json:"signature"`
		Signature_len 	int			`bson:"signature_len" json:"signature_len"`
		Pubkey 			[]byte		`bson:"pubkey" json:"pubkey"`
		Pubkey_len 		int			`bson:"pubkey_len" json:"pubkey_len"`
	}
)


func (p *BBcSignature) Stringer() string {
	if p.Key_type == 0 {
		return "  Not initialized\n"
	}
	ret :=  fmt.Sprintf("  key_type: %d\n", p.Key_type)
	ret +=  fmt.Sprintf("  signature: %x\n", p.Signature)
	ret +=  fmt.Sprintf("  pubkey: %x\n", p.Pubkey)
	return ret
}


func (p *BBcSignature) Serialize() ([]byte, error) {
	p.Pubkey_len = len(p.Pubkey) * 8
	p.Signature_len = len(p.Signature) * 8
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcSignature) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcSignatureDeserialize(format_type int, dat []byte) (BBcSignature, error) {
	if format_type != FORMAT_BINARY {
		return bbcSignatureDeserializeObj(format_type, dat)
	}
	obj := BBcSignature{}
	return obj, errors.New("not support the format")
}


func bbcSignatureDeserializeObj(format_type int, dat []byte) (BBcSignature, error) {
	obj := BBcSignature{}
	if format_type == FORMAT_BSON_COMPRESS_ZLIB || format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		if dat2, err := ZlibDecompress(&dat); err != nil {
			return obj, errors.New("failed to deserialize")
		} else {
			dat = dat2
		}
	}
	err := bson.Unmarshal(dat, &obj)
	return obj, err
}


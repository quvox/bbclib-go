package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcCrossRef struct {
		Format_type		int		`bson:"-" json:"-"`
		Id_length 		int		`bson:"-" json:"-"`
		Domain_id 		*[]byte	`bson:"domain_id" json:"domain_id"`
		Transaction_id	*[]byte	`bson:"transaction_id" json:"transaction_id"`
	}
)


func (p *BBcCrossRef) Stringer() string {
	ret := "Cross_Ref:\n"
	ret += fmt.Sprintf("  domain_id: %x\n", p.Domain_id)
	ret += fmt.Sprintf("  transaction_id: %x\n", p.Transaction_id)
	return ret
}


func (p *BBcCrossRef) Serialize() ([]byte, error) {
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcCrossRef) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcCrossRefDeserialize(format_type int, dat []byte) (BBcCrossRef, error) {
	if format_type != FORMAT_BINARY {
		return bbcCrossRefDeserializeObj(format_type, dat)
	}
	obj := BBcCrossRef{}
	return obj, errors.New("not support the format")
}


func bbcCrossRefDeserializeObj(format_type int, dat []byte) (BBcCrossRef, error) {
	obj := BBcCrossRef{}
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


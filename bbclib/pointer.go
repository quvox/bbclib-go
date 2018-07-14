package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcPointer struct {
		Format_type 	int		`bson:"-" json:"-"`
		Id_length 		int		`bson:"-" json:"-"`
		Transaction_id 	*[]byte	`bson:"transaction_id" json:"transaction_id"`
		Asset_id 		*[]byte	`bson:"asset_id" json:"asset_id"`
	}
)


func (p *BBcPointer) Stringer() string {
	ret := fmt.Sprintf("     transaction_id: %x\n", p.Transaction_id)
	ret += fmt.Sprintf("     asset_id: %x\n", p.Asset_id)
	return ret
}


func (p *BBcPointer) Serialize() ([]byte, error) {
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcPointer) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcPointerDeserialize(format_type int, dat []byte) (BBcPointer, error) {
	if format_type != FORMAT_BINARY {
		return bbcPointerDeserializeObj(format_type, dat)
	}
	obj := BBcPointer{}
	return obj, errors.New("not support the format")
}


func bbcPointerDeserializeObj(format_type int, dat []byte) (BBcPointer, error) {
	obj := BBcPointer{}
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

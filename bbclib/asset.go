package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcAsset struct {
		Format_type 	int		`bson:"-" json:"-"`
		Id_length		int		`bson:"-" json:"-"`
		Asset_id		[]byte	`bson:"asset_id" json:"asset_id"`
		User_id			[]byte	`bson:"user_id" json:"user_id"`
		Nonce			[]byte	`bson:"nonce" json:"nonce"`
		Asset_file_size int		`bson:"asset_file_size" json:"asset_file_size"`
		Asset_file 		*[]byte	`bson:"-" json:"-"`
		Asset_file_digest	*[]byte	`bson:"asset_file_digest" json:"asset_file_digest"`
		Asset_body_size int		`bson:"asset_body_size" json:"asset_body_size"`
		Asset_body		*[]byte	`bson:"asset_body" json:"asset_body"`
	}
)


func (p *BBcAsset) Stringer() string {
	ret :=  "  Asset:\n"
	ret += fmt.Sprintf("     asset_id: %x\n", p.Asset_id)
	ret += fmt.Sprintf("     user_id: %x\n", p.User_id)
	ret += fmt.Sprintf("     nonce: %x\n", p.Nonce)
	ret += fmt.Sprintf("     file_size: %d\n", p.Asset_file_size)
	if p.Asset_file_digest != nil {
		ret += fmt.Sprintf("     file_digest: %x\n", p.Asset_file)
	} else {
		ret += "     file_digest: None\n"
	}
	ret += fmt.Sprintf("     body_size: %d\n", p.Asset_body_size)
	ret += fmt.Sprintf("     body: %v\n", p.Asset_body)
	return ret
}


func (p *BBcAsset) Serialize() ([]byte, error) {
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcAsset) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcAssetDeserialize(format_type int, dat []byte) (BBcAsset, error) {
	if format_type != FORMAT_BINARY {
		return bbcAssetDeserializeObj(format_type, dat)
	}
	obj := BBcAsset{}
	return obj, errors.New("not support the format")
}


func bbcAssetDeserializeObj(format_type int, dat []byte) (BBcAsset, error) {
	obj := BBcAsset{}
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


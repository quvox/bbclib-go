package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcReference struct {
		Format_type 		int		`bson:"-" json:"-"`
		Id_length 			int		`bson:"-" json:"-"`
		Asset_group_id 		[]byte	`bson:"asset_group_id" json:"asset_group_id"`
		Transaction_id 		[]byte	`bson:"transaction_id" json:"transaction_id"`
		Event_index_in_ref 	int		`bson:"event_index_in_ref" json:"event_index_in_ref"`
		Sig_indices 		[]int	`bson:"sig_indices" json:"sig_indices"`
	}
)

func (p *BBcReference) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.Asset_group_id)
	ret += fmt.Sprintf("  transaction_id: %x\n", p.Transaction_id)
	ret += fmt.Sprintf("  event_index_in_ref: %v\n", p.Event_index_in_ref)
	ret += fmt.Sprintf("  sig_indices: %v\n", p.Sig_indices)
	return ret
}


func (p *BBcReference) Serialize() ([]byte, error) {
	if len(p.Asset_group_id) > p.Id_length {
		p.Asset_group_id = p.Asset_group_id[:p.Id_length]
	}
	if len(p.Transaction_id) > p.Id_length {
		p.Transaction_id = p.Transaction_id[:p.Id_length]
	}
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcReference) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcReferenceDeserialize(format_type int, dat []byte) (BBcReference, error) {
	if format_type != FORMAT_BINARY {
		return bbcReferenceDeserializeObj(format_type, dat)
	}
	obj := BBcReference{}
	return obj, errors.New("not support the format")
}


func bbcReferenceDeserializeObj(format_type int, dat []byte) (BBcReference, error) {
	obj := BBcReference{}
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


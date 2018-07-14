package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcRelation struct {
		Format_type		int				`bson:"-" json:"-"`
		Id_length		int				`bson:"-" json:"-"`
		Asset_group_id	[]byte			`bson:"asset_group_id" json:"asset_group_id"`
		Pointers 		[]BBcPointer	`bson:"pointers" json:"pointers"`
		Asset 			BBcAsset		`bson:"asset" json:"asset"`
	}
)


func (p *BBcRelation) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.Asset_group_id)
	if p.Pointers != nil {
		ret += fmt.Sprintf("  Pointers[]: %d\n", len(p.Pointers))
		for i := range p.Pointers {
			ret += fmt.Sprintf("   [%d]\n", i)
			ret += p.Pointers[i].Stringer()
		}
	}
	ret += p.Asset.Stringer()
	return ret
}


func (p *BBcRelation) Serialize() ([]byte, error) {
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcRelation) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcRelationDeserialize(format_type int, dat []byte) (BBcRelation, error) {
	if format_type != FORMAT_BINARY {
		return bbcRelationDeserializeObj(dat)
	}
	obj := BBcRelation{}
	return obj, errors.New("not support the format")
}


func bbcRelationDeserializeObj(dat []byte) (BBcRelation, error) {
	obj := BBcRelation{}
	err := bson.Unmarshal(dat, &obj)
	return obj, err
}


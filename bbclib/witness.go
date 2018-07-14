package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcWitness struct {
		Format_type int			`bson:"-" json:"-"`
		Id_length 	int			`bson:"-" json:"-"`
		User_ids 	[][]byte	`bson:"user_ids" json:"user_ids"`
		Sig_indices []int		`bson:"sig_indices" json:"sig_indices"`
	}
)


func (p *BBcWitness) Stringer() string {
	ret := "Witness:\n"
	if p.User_ids != nil {
		for i := range p.User_ids {
			ret += fmt.Sprintf(" [%d]\n", i)
			ret += fmt.Sprintf(" user_id: %x\n", p.User_ids[i])
			ret += fmt.Sprintf(" sig_index: %d\n", p.Sig_indices[i])
		}
	} else {
		ret += "  None (invalid)\n"
	}
	return ret
}


func (p *BBcWitness) Serialize() ([]byte, error) {
	if p.User_ids != nil {
		for i := range p.User_ids {
			if len(p.User_ids[i]) > p.Id_length {
				p.User_ids[i] = p.User_ids[i][:p.Id_length]
			}
		}
	}
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcWitness) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcWitnessDeserialize(format_type int, dat []byte) (BBcWitness, error) {
	if format_type != FORMAT_BINARY {
		return bbcWitnessDeserializeObj(format_type, dat)
	}
	obj := BBcWitness{}
	return obj, errors.New("not support the format")
}


func bbcWitnessDeserializeObj(format_type int, dat []byte) (BBcWitness, error) {
	obj := BBcWitness{}
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


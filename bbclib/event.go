package bbclib

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fmt"
)

type (
	BBcEvent struct {
		Format_type 		int			`bson:"-" json:"-"`
		Id_length 			int			`bson:"-" json:"-"`
		Asset_group_id 		[]byte		`bson:"asset_group_id" json:"asset_group_id"`
		Reference_indices 	[]int		`bson:"reference_indices" json:"reference_indices"`
		Mandatory_approvers [][]byte	`bson:"mandatory_approvers" json:"mandatory_approvers"`
		Option_approver_num_numerator 	int		`bson:"option_approver_num_numerator" json:"option_approver_num_numerator"`
		Option_approver_num_denominator	int		`bson:"option_approver_num_denominator" json:"option_approver_num_denominator"`
		Option_approvers 	[][]byte	`bson:"option_approvers" json:"option_approvers"`
		Asset				BBcAsset	`bson:"asset" json:"asset"`
	}
)


func (p *BBcEvent) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.Asset_group_id)
	ret += fmt.Sprintf("  reference_indices: %v\n", p.Reference_indices)
	ret += "  mandatory_approvers:\n"
	if p.Mandatory_approvers != nil {
		for i := range p.Mandatory_approvers {
			ret += fmt.Sprintf("    - %x\n", p.Mandatory_approvers[i])
		}
	} else {
		ret += "    - None\n"
	}
	ret += "  option_approvers:\n"
	if p.Option_approvers != nil {
		for i := range p.Option_approvers {
			ret += fmt.Sprintf("    - %x\n", p.Option_approvers[i])
		}
	} else {
		ret += "    - None\n"
	}
	ret += fmt.Sprintf("  option_approver_num_numerator: %d\n", p.Option_approver_num_numerator)
	ret += fmt.Sprintf("  option_approver_num_denominator: %d\n", p.Option_approver_num_denominator)
	ret += p.Asset.Stringer()
	return ret
}


func (p *BBcEvent) Serialize() ([]byte, error) {
	if len(p.Asset_group_id) > p.Id_length {
		p.Asset_group_id = p.Asset_group_id[:p.Id_length]
	}
	if p.Mandatory_approvers != nil {
		for i := range p.Mandatory_approvers {
			if len(p.Mandatory_approvers[i]) > p.Id_length {
				p.Mandatory_approvers[i] = p.Mandatory_approvers[i][:p.Id_length]
			}
		}
	}
	if p.Option_approvers != nil {
		for i := range p.Option_approvers {
			if len(p.Option_approvers[i]) > p.Id_length {
				p.Option_approvers[i] = p.Option_approvers[i][:p.Id_length]
			}
		}
	}
	if p.Format_type != FORMAT_BINARY {
		return p.serializeObj()
	}
	return nil, errors.New("not support the format")
}


func (p *BBcEvent) serializeObj() ([]byte, error) {
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcEventDeserialize(format_type int, dat []byte) (BBcEvent, error) {
	if format_type != FORMAT_BINARY {
		return bbcEventDeserializeObj(format_type, dat)
	}
	obj := BBcEvent{}
	return obj, errors.New("not support the format")
}


func bbcEventDeserializeObj(format_type int, dat []byte) (BBcEvent, error) {
	obj := BBcEvent{}
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


package bbclib

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"encoding/binary"
	"crypto/sha256"
	"unsafe"
	"encoding/json"
	"fmt"
)

type (
	BBcHeader struct {
		Version 	int				`bson:"version" json:"version"`
		Timestamp 	int				`bson:"timestamp" json:"timestamp"`
		Id_length 	int				`bson:"id_length" json:"id_length"`
	}

	BBcTransactionBase struct {
		Header		BBcHeader		`bson:"header" json:"header"`
		Events 		[]BBcEvent		`bson:"events" json:"events"`
		References 	[]BBcReference	`bson:"references" json:"references"`
		Relations 	[]BBcRelation	`bson:"relations" json:"relations"`
		Witness 	*BBcWitness		`bson:"witness" json:"witness"`
	}

	BBcTransactionForId struct {
		Tx_base		[]byte			`bson:"tx_base" json:"tx_base"`
		Crossref	*BBcCrossRef	`bson:"cross_ref" json:"cross_ref"`
	}

	BBcTransaction struct {
		Format_type 	int					`bson:"-" json:"-"`
		Transaction_id	[]byte				`bson:"-" json:"-"`
		Tx_base			BBcTransactionBase	`bson:"transaction_base" json:"transaction_base"`
		Crossref		*BBcCrossRef		`bson:"cross_ref" json:"cross_ref"`
		Signatures 		[]BBcSignature		`bson:"signatures" json:"signatures"`
	}
)


const (
	KeyType_NOT_INITIALIZED = 0
	KeyType_ECDSA_SECP256k1 = 1
	KeyType_ECDSA_P256v1 = 2
	FORMAT_BINARY = 0
	FORMAT_BSON = 1
	FORMAT_BSON_COMPRESS_BZ2 = 2
	FORMAT_BSON_COMPRESS_ZLIB = 3
	FORMAT_MSGPACK = 4
	FORMAT_MSGPACK_COMPRESS_BZ2 = 5
	FORMAT_MSGPACK_COMPRESS_ZLIB = 6
)


func (p *BBcTransaction) Stringer() string {
	var ret string
	ret =  "------- Dump of the transaction data ------\n"
	ret += fmt.Sprintf("* transaction_id: %x\n", p.Transaction_id)
	ret += fmt.Sprintf("version: %d\n", p.Tx_base.Header.Version)
	ret += fmt.Sprintf("timestamp: %d\n", p.Tx_base.Header.Timestamp)
	if p.Tx_base.Header.Version != 0 {
		ret += fmt.Sprintf("id_length: %d\n", p.Tx_base.Header.Id_length)
	}

	ret += fmt.Sprintf("Event[]: %d\n", len(p.Tx_base.Events))
	for i := range p.Tx_base.Events {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Tx_base.Events[i].Stringer()
	}

	ret += fmt.Sprintf("Reference[]: %d\n", len(p.Tx_base.References))
	for i := range p.Tx_base.References {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Tx_base.References[i].Stringer()
	}

	ret += fmt.Sprintf("Relation[]: %d\n", len(p.Tx_base.Relations))
	for i := range p.Tx_base.Relations {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Tx_base.Relations[i].Stringer()
	}

	if p.Tx_base.Witness != nil {
		ret += p.Tx_base.Witness.Stringer()
	} else {
		ret += "Witness: None\n"
	}

	if p.Crossref != nil {
		ret += p.Crossref.Stringer()
	} else {
		ret += "Cross_Ref: None\n"
	}

	ret += fmt.Sprintf("Signature[]: %d\n", len(p.Signatures))
	for i := range p.Signatures {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Signatures[i].Stringer()
	}
	return ret
}


func (p *BBcTransaction) Digest() []byte {
	txdata, err := p.Serialize(true)
	if err != nil {
		return nil
	}
	fmt.Printf("DIGESTtarget: %x\n", txdata)
	var dt map[string]interface{}
	bson.Unmarshal(txdata, &dt)
	fmt.Printf("DIGESTtarget(content): %v\n", dt)

	digest := sha256.Sum256(txdata)
	p.Transaction_id = digest[:p.Tx_base.Header.Id_length]
	return digest[:]
}


func (p *BBcTransaction) Jsonify() string {
	dat, _ := json.Marshal(p)
	return *(*string)(unsafe.Pointer(&dat))
}


func (p *BBcTransaction) Serialize(forId bool) ([]byte, error) {
	if p.Tx_base.Header.Version == 0 {
		return nil, errors.New("not support version=0 transaction")
	}
	if p.Format_type != FORMAT_BINARY {
		if forId {
			txCore, err := p.serializeObj(forId)
			if err != nil {
				return nil, err
			}
			coreDigest := sha256.Sum256(txCore)
			fmt.Printf("coreDigest: %x\n", coreDigest)

			return bson.Marshal(&BBcTransactionForId {
				Tx_base: coreDigest[:],
				Crossref: p.Crossref,
			})
		}
		dat, err := p.serializeObj(forId)
		if err != nil {
			return nil, err
		}
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, (uint16)(p.Format_type))
		return append(b, dat...), nil
	}
	return nil, errors.New("not support the format")
}


func (p *BBcTransaction) serializeObj(forId bool) ([]byte, error) {
	if forId {
		dat, err := bson.Marshal((*BBcTransactionBase)(unsafe.Pointer(&p.Tx_base)))
		return dat, err
	}
	dat, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	if p.Format_type == FORMAT_BSON_COMPRESS_ZLIB || p.Format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		return ZlibCompress(&dat), nil
	}
	return dat, err
}


func BBcTransactionDeserialize(dat []byte) (*BBcTransaction, error) {
	format_type := (int)(binary.LittleEndian.Uint16(dat[:2]))
	if format_type != FORMAT_BINARY {
		return bbcTransactionDeserializeObj(format_type, dat[2:])
	}
	obj := BBcTransaction{}
	obj.Format_type = format_type

	return &obj, errors.New("not support the format")
}


func bbcTransactionDeserializeObj(format_type int, dat []byte) (*BBcTransaction, error) {
	var obj BBcTransaction
	if format_type == FORMAT_BSON_COMPRESS_ZLIB || format_type == FORMAT_MSGPACK_COMPRESS_ZLIB {
		if dat2, err := ZlibDecompress(&dat); err != nil {
			return nil, errors.New("failed to deserialize")
		} else {
			dat = dat2
		}
	}
	err := bson.Unmarshal(dat, &obj)
	if err != nil {
		return nil, err
	}
	obj.Format_type = format_type
	idLength := obj.Tx_base.Header.Id_length
	obj.Digest()

	if obj.Tx_base.Events != nil {
		for i := range obj.Tx_base.Events {
			obj.Tx_base.Events[i].Format_type = format_type
			obj.Tx_base.Events[i].Id_length = idLength
			obj.Tx_base.Events[i].Asset.Format_type = format_type
			obj.Tx_base.Events[i].Asset.Id_length = idLength
		}
	}
	if obj.Tx_base.References != nil {
		for i := range obj.Tx_base.References {
			obj.Tx_base.References[i].Format_type = format_type
			obj.Tx_base.References[i].Id_length = idLength
		}
	}
	if obj.Tx_base.Relations != nil {
		for i := range obj.Tx_base.Relations {
			obj.Tx_base.Relations[i].Format_type = format_type
			obj.Tx_base.Relations[i].Id_length = idLength
			obj.Tx_base.Relations[i].Asset.Format_type = format_type
			obj.Tx_base.Relations[i].Asset.Id_length = idLength
			if obj.Tx_base.Relations[i].Pointers != nil {
				for j := range obj.Tx_base.Relations[i].Pointers {
					obj.Tx_base.Relations[i].Pointers[j].Format_type = format_type
					obj.Tx_base.Relations[i].Pointers[j].Id_length = idLength
				}
			}
		}
	}
	if obj.Signatures != nil {
		for i := range obj.Signatures {
			obj.Signatures[i].Format_type = format_type
		}
	}

	return &obj, nil
}


func (p *BBcTransaction) VerifyAll() (bool, int) {
	digest := p.Digest()
	for i := range p.Signatures {
		if p.Signatures[i].Key_type == 0 {
			continue
		}
		if ret := VerifyBBcSignature(digest, &p.Signatures[i]); ! ret {
			return false, i
		}
	}
	return true, -1
}

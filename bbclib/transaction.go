package bbclib

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type (
	BBcTransaction struct {
		digestCalculating bool
		TransactionId	[]byte
		TransactionBaseDigest	[]byte
		TransactionData	[]byte
		SigIndices		[][]byte
		Version 		uint32
		Timestamp 		int64
		IdLength 		int
		Events 			[]*BBcEvent
		References 		[]*BBcReference
		Relations 		[]*BBcRelation
		Witness 		*BBcWitness
		Crossref		*BBcCrossRef
		Signatures		[]*BBcSignature
	}
)


/*
const (
	KeyType_NOT_INITIALIZED = 0
	KeyType_ECDSA_SECP256k1 = 1
	KeyType_ECDSA_P256v1 = 2
	FORMAT_PLAIN = 0x0000
	FORMAT_ZLIB = 0x0010
)
*/


func (p *BBcTransaction) Stringer() string {
	var ret string
	ret =  "------- Dump of the transaction data ------\n"
	ret += fmt.Sprintf("* transaction_id: %x\n", p.TransactionId)
	ret += fmt.Sprintf("version: %d\n", p.Version)
	ret += fmt.Sprintf("timestamp: %d\n", p.Timestamp)
	if p.Version != 0 {
		ret += fmt.Sprintf("id_length: %d\n", p.IdLength)
	}

	ret += fmt.Sprintf("Event[]: %d\n", len(p.Events))
	for i := range p.Events {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Events[i].Stringer()
	}

	ret += fmt.Sprintf("Reference[]: %d\n", len(p.References))
	for i := range p.References {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.References[i].Stringer()
	}

	ret += fmt.Sprintf("Relation[]: %d\n", len(p.Relations))
	for i := range p.Relations {
		ret += fmt.Sprintf(" [%d]\n", i)
		ret += p.Relations[i].Stringer()
	}

	if p.Witness != nil {
		ret += p.Witness.Stringer()
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

func (p *BBcTransaction) AddEvent(obj *BBcEvent) {
	obj.IdLength = p.IdLength
	p.Events = append(p.Events, obj)
}

func (p *BBcTransaction) AddReference(obj *BBcReference) {
	obj.IdLength = p.IdLength
	p.References = append(p.References, obj)
	obj.Transaction = p
}

func (p *BBcTransaction) AddRelation(obj *BBcRelation) {
	obj.IdLength = p.IdLength
	p.Relations = append(p.Relations, obj)
}

func (p *BBcTransaction) AddWitness(obj *BBcWitness) {
	obj.IdLength = p.IdLength
	p.Witness = obj
	obj.Transaction = p
}

func (p *BBcTransaction) AddCrossRef(obj *BBcCrossRef) {
	obj.IdLength = p.IdLength
	p.Crossref = obj
}

func (p *BBcTransaction) AddSignature(userId *[]byte, sig *BBcSignature) {
	for i := range p.SigIndices {
		if reflect.DeepEqual(p.SigIndices[i], userId) {
			p.Signatures[i] = sig
			return
		}
	}
	uid := make([]byte, int(p.IdLength))
	copy(uid, *userId)
	p.SigIndices = append(p.SigIndices, uid)
	p.Signatures = append(p.Signatures, sig)
}

func (p *BBcTransaction) GetSigIndex(uid []byte) int {
	var i = -1
	for i = range p.SigIndices {
		if reflect.DeepEqual(p.SigIndices[i], uid) {
			return i
		}
	}
	p.SigIndices = append(p.SigIndices, uid)
	return i + 1
}

func (p *BBcTransaction) Sign(keypair *KeyPair) ([]byte, error) {
	if p.TransactionId == nil {
		p.Digest()
	}
	signature := keypair.Sign(p.TransactionId)
	if signature == nil {
		return nil, errors.New("fail to sign")
	}
	return signature, nil
}

func (p *BBcTransaction) VerifyAll() (bool, int) {
	digest := p.Digest()
	for i := range p.Signatures {
		if p.Signatures[i].KeyType == KeyType_NOT_INITIALIZED {
			continue
		}
		if ret := VerifyBBcSignature(digest, p.Signatures[i]); ! ret {
			return false, i
		}
	}
	return true, -1
}


func (p *BBcTransaction) Digest() []byte {
	p.digestCalculating = true
	if p.TransactionId == nil {
		p.TransactionId = make([]byte, p.IdLength)
	}
	buf := new(bytes.Buffer)

	err := p.packBase(buf)
	if err != nil {
		p.digestCalculating = false
		return nil
	}

	buf = new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, p.TransactionBaseDigest); err != nil {
		p.digestCalculating = false
		return nil
	}

	err = p.packCrossRef(buf)
	if err != nil {
		p.digestCalculating = false
		return nil
	}

	digest := sha256.Sum256(buf.Bytes())
	p.TransactionId = digest[:p.IdLength]
	p.digestCalculating = false
	return digest[:]
}


func (p *BBcTransaction) packCrossRef(buf *bytes.Buffer) error {
	if p.Crossref != nil {
		dat, err := p.Crossref.Pack()
		if err != nil {
			return err
		}
		Put2byte(buf, 1)
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	} else {
		Put2byte(buf, 0)
	}
	return nil
}

func (p *BBcTransaction) packBase(buf *bytes.Buffer) error {
	Put4byte(buf, p.Version)
	Put8byte(buf, p.Timestamp)
	Put2byte(buf, uint16(p.IdLength))

	Put2byte(buf, uint16(len(p.Events)))
	for _, obj := range p.Events {
		dat, err := obj.Pack()
		if err != nil {
			return err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	}

	Put2byte(buf, uint16(len(p.References)))
	for _, obj := range p.References {
		dat, err := obj.Pack()
		if err != nil {
			return err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	}

	Put2byte(buf, uint16(len(p.Relations)))
	for _, obj := range p.Relations {
		dat, err := obj.Pack()
		if err != nil {
			return err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	}

	if p.Witness != nil {
		dat, err := p.Witness.Pack()
		if err != nil {
			return err
		}
		Put2byte(buf, 1)
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return err
		}
	} else {
		Put2byte(buf, 0)
	}

	digest := sha256.Sum256(buf.Bytes())
	p.TransactionBaseDigest = digest[:]

	return nil
}


func (p *BBcTransaction) Pack() ([]byte, error) {
	if ! p.digestCalculating && p.TransactionId == nil {
		p.Digest()
	}

	if p.Version == 0 {
		return nil, errors.New("not support version=0 transaction")
	}

	buf := new(bytes.Buffer)
	err := p.packBase(buf)
	if err != nil {
		return nil, err
	}
	err = p.packCrossRef(buf)
	if err != nil {
		return nil, err
	}

	Put2byte(buf, uint16(len(p.Signatures)))
	for _, obj := range p.Signatures {
		dat, err := obj.Pack()
		if err != nil {
			return nil, err
		}
		Put4byte(buf, uint32(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return nil, err
		}
	}

	p.TransactionData = buf.Bytes()
	return p.TransactionData, nil
}


func (p *BBcTransaction) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.Version, err = Get4byte(buf)
	if err != nil {
		return err
	}

	p.Timestamp, err = Get8byte(buf)
	if err != nil {
		return err
	}

	idLen, err := Get2byte(buf)
	if err != nil {
		return err
	}
	p.IdLength = int(idLen)

	num, err := Get2byte(buf)
	for i := 0; i < int(num); i++ {
		plen, err := Get4byte(buf)
		if err != nil {
			return err
		}
		dat, err := GetBytes(buf, int(plen))
		obj := BBcEvent{IdLength:p.IdLength}
		obj.Unpack(&dat)
		p.Events = append(p.Events, &obj)
	}

	num, err = Get2byte(buf)
	for i := 0; i < int(num); i++ {
		plen, err := Get4byte(buf)
		if err != nil {
			return err
		}
		dat, err := GetBytes(buf, int(plen))
		obj := BBcReference{IdLength:p.IdLength}
		obj.Unpack(&dat)
		p.References = append(p.References, &obj)
	}

	num, err = Get2byte(buf)
	for i := 0; i < int(num); i++ {
		plen, err := Get4byte(buf)
		if err != nil {
			return err
		}
		dat, err := GetBytes(buf, int(plen))
		obj := BBcRelation{IdLength:p.IdLength}
		obj.Unpack(&dat)
		p.Relations = append(p.Relations, &obj)
	}

	num, err = Get2byte(buf)
	if num > 0 {
		size, err := Get4byte(buf)
		dat, err := GetBytes(buf, int(size))
		if err != nil {
			return err
		}
		p.Witness = &BBcWitness{IdLength:p.IdLength}
		p.Witness.Unpack(&dat)
	}

	num, err = Get2byte(buf)
	if num > 0 {
		size, err := Get4byte(buf)
		dat, err := GetBytes(buf, int(size))
		if err != nil {
			return err
		}
		p.Crossref = &BBcCrossRef{IdLength:p.IdLength}
		p.Crossref.Unpack(&dat)
	}

	num, err = Get2byte(buf)
	for i := 0; i < int(num); i++ {
		plen, err := Get4byte(buf)
		if err != nil {
			return err
		}
		dat, err := GetBytes(buf, int(plen))
		obj := BBcSignature{}
		obj.Unpack(dat)
		p.Signatures = append(p.Signatures, &obj)
	}

	p.Digest()
	return nil
}

/*
Copyright (c) 2018 Zettant Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 */

package bbclib

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

/*
This is the BBcTransaction definition.

BBcTransaction is just a container of various objects.

Events, References, Relations and Signatures are list of BBcEvent, BBcReference, BBcRelation and BBcSignature objects, respectively.
"digestCalculating", "TransactionBaseDigest", "TransactionData" and "SigIndices" are not included in the packed data. They are internal use only.

Calculating TransactionId

How to calculate the TransactionId of the transaction is a little bit complicated, meaning that 2-step manner.
This is because inter-domain transaction authenticity (i.e., CrossReference) can be conducted in secure manner.
By presenting TransactionBaseDigest (see below) to an outer-domain, the domain user can confirm the existence of the transaction in the past.
(no need to present whole transaction data including the asset information).

1st step:
  * Pack info (from version to Witness) by packBase()
  * Calculate SHA256 digest of the packed info. This value is TransactionBaseDigest.

2nd step:
  * Pack BBcCrossRef object to get packed data by packCrossRef()
  * Concatenate TransactionBaseDigest and the packed BBcCrossRef
  * Calculate SHA256 digest of the concatenated data. This value is TransactionId
 */
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


// Output content of the object
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

// Add BBcEvent object in the transaction object
func (p *BBcTransaction) AddEvent(obj *BBcEvent) {
	obj.IdLength = p.IdLength
	p.Events = append(p.Events, obj)
}

// Add BBcReference object in the transaction object
func (p *BBcTransaction) AddReference(obj *BBcReference) {
	obj.IdLength = p.IdLength
	p.References = append(p.References, obj)
	obj.Transaction = p
}

// Add BBcRelation object in the transaction object
func (p *BBcTransaction) AddRelation(obj *BBcRelation) {
	obj.IdLength = p.IdLength
	p.Relations = append(p.Relations, obj)
}

// Add BBcWitness object in the transaction object
func (p *BBcTransaction) AddWitness(obj *BBcWitness) {
	obj.IdLength = p.IdLength
	p.Witness = obj
	obj.Transaction = p
}

// Add BBcCrossRef object in the transaction object
func (p *BBcTransaction) AddCrossRef(obj *BBcCrossRef) {
	obj.IdLength = p.IdLength
	p.Crossref = obj
}

// Add BBcSignature object for the specified userId in the transaction object
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

// Get position (index) of the corespondent userId in the signature list
func (p *BBcTransaction) GetSigIndex(userId []byte) int {
	var i = -1
	for i = range p.SigIndices {
		if reflect.DeepEqual(p.SigIndices[i], userId) {
			return i
		}
	}
	p.SigIndices = append(p.SigIndices, userId)
	return i + 1
}

// Sign TransactionId using private key in the given keypair
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

// Verify TransactionId with all BBcSignature objects in the transaction
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

// Calculate TransactionId of the BBcTransaction object
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

// Pack only BBcCrossRef object in binary data
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

// Pack the base part of BBcTransaction object in binary data (from version to witness)
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

// Pack BBcTransaction object in binary data
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

// Unpack binary data to BBcTransaction object
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

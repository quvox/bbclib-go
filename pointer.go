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
	"fmt"
)

/*
This is the BBcPointer definition.

BBcPointer(s) are included in BBcRelation object. A BBcPointer object includes "TransactionId" and "AssetId" and
declares that the transaction has a certain relationship with the BBcTransaction and BBcAsset object specified by those IDs.

IdLength is not included in a packed data. It is for internal use only.
 */
type (
	BBcPointer struct {
		IdLength 		int
		TransactionId 	[]byte
		AssetId 		[]byte
	}
)

// Output content of the object
func (p *BBcPointer) Stringer() string {
	ret := fmt.Sprintf("     transaction_id: %x\n", p.TransactionId)
	ret += fmt.Sprintf("     asset_id: %x\n", p.AssetId)
	return ret
}

// Add essential information to the BBcPointer object
func (p *BBcPointer) Add(txid *[]byte, asid *[]byte) {
	if txid != nil {
		p.TransactionId = make([]byte, p.IdLength)
		copy(p.TransactionId, (*txid)[:p.IdLength])
	}
	if asid != nil {
		p.AssetId = make([]byte, p.IdLength)
		copy(p.AssetId, (*asid)[:p.IdLength])
	}
}

// Pack BBcPointer object in binary data
func (p *BBcPointer) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.TransactionId, p.IdLength)

	if p.AssetId != nil {
		Put2byte(buf, 1)
	} else {
		Put2byte(buf, 0)
		return buf.Bytes(), nil
	}

	PutBigInt(buf, &p.AssetId, p.IdLength)

	return buf.Bytes(), nil
}

// Unpack binary data to BBcPointer object
func (p *BBcPointer) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.TransactionId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	if val, err := Get2byte(buf); err != nil {
		return err
	} else {
		if val == 0 {
			p.AssetId = nil
			return nil
		}
	}

	p.AssetId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	return nil
}

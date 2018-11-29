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
This is the BBcCrossRef definition.

CrossRef stands for CrossReference, which holds information in other domain for inter-domain collaboration of transaction authenticity.

"IdLength" is not included in a packed data. It is for internal use only.

"DomainId" is the identifier of a domain and the length of the ID must be 256 bits (=32 bytes).
"TransactionId" is that of transaction object in other domain (specified by the DomainId).
 */
type (
	BBcCrossRef struct {
		IdLength 		int
		DomainId 		[]byte
		TransactionId	[]byte
	}
)

// The length of DomainID must be 256-bit in any domain.
const (
	DomainIdLength = 32
)

// Output content of the object
func (p *BBcCrossRef) Stringer() string {
	ret := "Cross_Ref:\n"
	ret += fmt.Sprintf("  domain_id: %x\n", p.DomainId)
	ret += fmt.Sprintf("  transaction_id: %x\n", p.TransactionId)
	return ret
}

// Add essential information to the BBcCrossRef object
func (p *BBcCrossRef) Add(domainId *[]byte, txid *[]byte) {
	if domainId != nil {
		p.DomainId = make([]byte, DomainIdLength)
		copy(p.DomainId, *domainId)
	}
	if txid != nil {
		p.TransactionId = make([]byte, p.IdLength)
		copy(p.TransactionId, (*txid)[:p.IdLength])
	}
}

// Pack BBcCrossRef object in binary data
func (p *BBcCrossRef) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.DomainId, DomainIdLength)
	PutBigInt(buf, &p.TransactionId, 32)

	return buf.Bytes(), nil
}

// Unpack binary data to BBcCrossRef object
func (p *BBcCrossRef) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.DomainId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.TransactionId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	return nil
}

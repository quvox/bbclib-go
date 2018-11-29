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
	"errors"
	"fmt"
)

/*
This is the BBcWitness definition.

The BBcWitness has the mapping info between the userIds and BBcSignature objects.
This object should be used if BBcRelation is used or a certain user wants to sign to the transaction in some reason.

"UserIds" is the list of userId, and "SigIndices" is a mapping info between userId and the position (index) of the signature list in the BBcTransaction object.

"Transaction" is the pointer to the parent BBcTransaction object.

"IdLength" and "Transaction" are not included in a packed data. They are for internal use only.
 */
type (
	BBcWitness struct {
		IdLength 	int
		UserIds 	[][]byte
		SigIndices	[]int
		Transaction *BBcTransaction
	}
)

// Output content of the object
func (p *BBcWitness) Stringer() string {
	ret := "Witness:\n"
	if p.UserIds != nil {
		for i := range p.UserIds {
			ret += fmt.Sprintf(" [%d]\n", i)
			ret += fmt.Sprintf(" user_id: %x\n", p.UserIds[i])
			ret += fmt.Sprintf(" sig_index: %d\n", p.SigIndices[i])
		}
	} else {
		ret += "  None (invalid)\n"
	}
	return ret
}

// Set pointer to the parent transaction object
func (p *BBcWitness) SetTransaction(txobj *BBcTransaction) {
	p.Transaction = txobj
}

// Make a memo for managing signer who sign this BBcTransaction object
// This must be done before AddSignature.
func (p *BBcWitness) AddWitness(userId *[]byte) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.UserIds = append(p.UserIds, (*userId)[:p.IdLength])
	idx := p.Transaction.GetSigIndex(*userId)
	p.SigIndices = append(p.SigIndices, idx)
	return nil
}

// Add BBcSignature to the parent BBcTransaction and the position in the Signatures list in BBcTransaction is based on the UserId
func (p *BBcWitness) AddSignature(userId *[]byte, sig *BBcSignature) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.Transaction.AddSignature(userId, sig)
	return nil
}

// Pack BBcWitness object in binary data
func (p *BBcWitness) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	Put2byte(buf, uint16(len(p.UserIds)))
	for i := 0; i < len(p.UserIds); i++ {
		PutBigInt(buf, &p.UserIds[i], p.IdLength)
		Put2byte(buf, uint16(p.SigIndices[i]))
	}

	return buf.Bytes(), nil
}

// Unpack binary data to BBcWitness object
func (p *BBcWitness) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	userNum, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(userNum); i++ {
		userId := make([]byte, p.IdLength)
		userId, err = GetBigInt(buf)
		if err != nil {
			return err
		}
		p.UserIds = append(p.UserIds, userId)

		idx, err := Get2byte(buf)
		if err != nil {
			return err
		}
		p.SigIndices = append(p.SigIndices, int(idx))
	}

	return nil
}

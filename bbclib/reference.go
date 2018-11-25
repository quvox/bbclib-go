package bbclib

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

type (
	BBcReference struct {
		IdLength 			int
		AssetGroupId 		[]byte
		TransactionId 		[]byte
		EventIndexInRef 	uint16
		SigIndices 			[]int
		Transaction			*BBcTransaction
		RefTransaction		*BBcTransaction
		RefEvent			BBcEvent
	}
)


func (p *BBcReference) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.AssetGroupId)
	ret += fmt.Sprintf("  transaction_id: %x\n", p.TransactionId)
	ret += fmt.Sprintf("  event_index_in_ref: %v\n", p.EventIndexInRef)
	ret += fmt.Sprintf("  sig_indices: %v\n", p.SigIndices)
	return ret
}


func (p *BBcReference) SetTransaction(txobj *BBcTransaction) {
	p.Transaction = txobj
}

func (p *BBcReference) Add(assetGroupId *[]byte, refTransaction *BBcTransaction, eventIdx int) {
	if assetGroupId != nil {
		p.AssetGroupId = make([]byte, p.IdLength)
		copy(p.AssetGroupId, *assetGroupId)
	}
	if eventIdx > -1 {
		p.EventIndexInRef = uint16(eventIdx)
	}
	if refTransaction != nil {
		p.RefTransaction = refTransaction
		p.TransactionId = refTransaction.TransactionId[:p.IdLength]
		p.RefEvent = *p.RefTransaction.Events[p.EventIndexInRef]
	}
}

func (p *BBcReference) AddApprover(userId *[]byte) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	if p.RefTransaction == nil {
		return errors.New("reference_transaction must be set")
	}

	flag := false
	for _, m := range p.RefEvent.MandatoryApprovers {
		if reflect.DeepEqual(m, userId) {
			flag = true
			break
		}
	}
	if ! flag {
		for _, o := range p.RefEvent.OptionApprovers {
			if reflect.DeepEqual(o, userId) {
				flag = true
				break
			}
		}
	}
	if ! flag {
		return errors.New("no user is specified as approver")
	}

	idx := p.Transaction.GetSigIndex(*userId)
	p.SigIndices = append(p.SigIndices, idx)
	return nil
}

func (p *BBcReference) AddSignature(userId *[]byte, sig *BBcSignature) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.Transaction.AddSignature(userId, sig)
	return nil
}


func (p *BBcReference) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.AssetGroupId, p.IdLength)
	PutBigInt(buf, &p.TransactionId, p.IdLength)
	Put2byte(buf, p.EventIndexInRef)
	Put2byte(buf, uint16(len(p.SigIndices)))
	for i := 0; i < len(p.SigIndices); i++ {
		Put2byte(buf, uint16(p.SigIndices[i]))
	}

	return buf.Bytes(), nil
}

func (p *BBcReference) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetGroupId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.TransactionId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	p.EventIndexInRef, err = Get2byte(buf)
	if err != nil {
		return err
	}

	sigNum, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(sigNum); i++ {
		idx, err := Get2byte(buf)
		if err != nil {
			return err
		}
		p.SigIndices = append(p.SigIndices, int(idx))
	}

	return nil
}

package bbclib

import (
	"bytes"
	"errors"
	"fmt"
)

type (
	BBcWitness struct {
		IdLength 	int
		UserIds 	[][]byte
		SigIndices	[]int
		Transaction *BBcTransaction
	}
)


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

func (p *BBcWitness) SetTransaction(txobj *BBcTransaction) {
	p.Transaction = txobj
}

func (p *BBcWitness) AddWitness(userId *[]byte) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.UserIds = append(p.UserIds, (*userId)[:p.IdLength])
	idx := p.Transaction.GetSigIndex(*userId)
	p.SigIndices = append(p.SigIndices, idx)
	return nil
}

func (p *BBcWitness) AddSignature(userId *[]byte, sig *BBcSignature) error {
	if p.Transaction == nil {
		return errors.New("transaction must be set")
	}
	p.Transaction.AddSignature(userId, sig)
	return nil
}


func (p *BBcWitness) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	Put2byte(buf, uint16(len(p.UserIds)))
	for i := 0; i < len(p.UserIds); i++ {
		PutBigInt(buf, &p.UserIds[i], p.IdLength)
		Put2byte(buf, uint16(p.SigIndices[i]))
	}

	return buf.Bytes(), nil
}


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

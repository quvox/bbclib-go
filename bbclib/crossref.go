package bbclib

import (
	"bytes"
	"fmt"
)

type (
	BBcCrossRef struct {
		IdLength 		int
		DomainId 		[]byte
		TransactionId	[]byte
	}
)

var (
	DomainIdLength = 32
)

func (p *BBcCrossRef) Stringer() string {
	ret := "Cross_Ref:\n"
	ret += fmt.Sprintf("  domain_id: %x\n", p.DomainId)
	ret += fmt.Sprintf("  transaction_id: %x\n", p.TransactionId)
	return ret
}


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

func (p *BBcCrossRef) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.DomainId, DomainIdLength)
	PutBigInt(buf, &p.TransactionId, 32)

	return buf.Bytes(), nil
}

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

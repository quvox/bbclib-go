package bbclib

import (
	"bytes"
	"fmt"
)

type (
	BBcPointer struct {
		IdLength 		int
		TransactionId 	[]byte
		AssetId 		[]byte
	}
)


func (p *BBcPointer) Stringer() string {
	ret := fmt.Sprintf("     transaction_id: %x\n", p.TransactionId)
	ret += fmt.Sprintf("     asset_id: %x\n", p.AssetId)
	return ret
}

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

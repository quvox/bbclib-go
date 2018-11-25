package bbclib

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type (
	BBcEvent struct {
		IdLength 			int
		AssetGroupId 		[]byte
		ReferenceIndices 	[]int
		MandatoryApprovers	[][]byte
		OptionApproverNumNumerator 		uint16
		OptionApproverNumDenominator	uint16
		OptionApprovers 	[][]byte
		Asset				*BBcAsset
	}
)


func (p *BBcEvent) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.AssetGroupId)
	if p.ReferenceIndices != nil {
		ret += fmt.Sprintf("  reference_indices: %v\n", p.ReferenceIndices)
	} else {
		ret += fmt.Sprintf("  reference_indices: None\n")
	}
	ret += "  mandatory_approvers:\n"
	if p.MandatoryApprovers != nil {
		for _, a := range p.MandatoryApprovers {
			ret += fmt.Sprintf("    - %x\n", a)
		}
	} else {
		ret += "    - None\n"
	}
	ret += "  option_approvers:\n"
	if p.OptionApprovers != nil {
		for _, o := range p.OptionApprovers {
			ret += fmt.Sprintf("    - %x\n", o)
		}
	} else {
		ret += "    - None\n"
	}
	ret += fmt.Sprintf("  option_approver_num_numerator: %d\n", p.OptionApproverNumNumerator)
	ret += fmt.Sprintf("  option_approver_num_denominator: %d\n", p.OptionApproverNumDenominator)
	if p.Asset != nil {
		ret += p.Asset.Stringer()
	} else {
		ret += fmt.Sprintf("  Asset: None\n")
	}
	return ret
}


func (p *BBcEvent) Add(assetGroupId *[]byte, asset *BBcAsset) {
	if assetGroupId != nil {
		p.AssetGroupId = make([]byte, p.IdLength)
		copy(p.AssetGroupId, (*assetGroupId)[:p.IdLength])
	}
	if asset != nil {
		p.Asset = asset
		p.Asset.IdLength = p.IdLength
	}
}

func (p *BBcEvent) AddReferenceIndex(relIndex int) {
	if relIndex != -1 {
		p.ReferenceIndices = append(p.ReferenceIndices, relIndex)
	}
}

func (p *BBcEvent) AddOptionParams(numerator int, denominator int) {
	p.OptionApproverNumNumerator = uint16(numerator)
	p.OptionApproverNumDenominator = uint16(denominator)
}

func (p *BBcEvent) AddMandatoryApprover(userId *[]byte) {
	uid := make([]byte, p.IdLength)
	copy(uid, *userId)
	p.MandatoryApprovers = append(p.MandatoryApprovers, uid)
}

func (p *BBcEvent) AddOptionApprover(userId *[]byte) {
	uid := make([]byte, p.IdLength)
	copy(uid, *userId)
	p.OptionApprovers = append(p.OptionApprovers, uid)
}


func (p *BBcEvent) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.AssetGroupId, p.IdLength)

	Put2byte(buf, uint16(len(p.ReferenceIndices)))
	for i := 0; i < len(p.ReferenceIndices); i++ {
		Put2byte(buf, uint16(p.ReferenceIndices[i]))
	}

	Put2byte(buf, uint16(len(p.MandatoryApprovers)))
	for i := 0; i < len(p.MandatoryApprovers); i++ {
		PutBigInt(buf, &p.MandatoryApprovers[i], p.IdLength)
	}

	Put2byte(buf, p.OptionApproverNumNumerator)
	Put2byte(buf, p.OptionApproverNumDenominator)
	Put2byte(buf, uint16(len(p.OptionApprovers)))
	for i := 0; i < len(p.OptionApprovers); i++ {
		PutBigInt(buf, &p.OptionApprovers[i], p.IdLength)
	}

	if p.Asset != nil {
		ast, err := p.Asset.Pack()
		if err != nil {
			return nil, err
		}
		Put4byte(buf, uint32(binary.Size(ast)))
		if err := binary.Write(buf, binary.LittleEndian, ast); err != nil {
			return nil, err
		}
	} else {
		Put4byte(buf, 0)
	}

	return buf.Bytes(), nil
}



func (p *BBcEvent) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetGroupId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	numReferences, err := Get2byte(buf)
	if err != nil {
		return err
	}
	for i := 0; i < int(numReferences); i++ {
		idx, err := Get2byte(buf)
		if err != nil {
			return err
		}
		p.ReferenceIndices = append(p.ReferenceIndices, int(idx))
	}

	numMandatory, err := Get2byte(buf)
	for i := 0; i < int(numMandatory); i++ {
		userId := make([]byte, p.IdLength)
		userId, err = GetBigInt(buf)
		if err != nil {
			return err
		}
		p.MandatoryApprovers = append(p.MandatoryApprovers, userId)
	}

	p.OptionApproverNumNumerator, err = Get2byte(buf)
	if err != nil {
		return err
	}
	p.OptionApproverNumDenominator, err = Get2byte(buf)
	if err != nil {
		return err
	}

	numOptional, err := Get2byte(buf)
	for i := 0; i < int(numOptional); i++ {
		userId := make([]byte, p.IdLength)
		userId, err = GetBigInt(buf)
		if err != nil {
			return err
		}
		p.OptionApprovers = append(p.OptionApprovers, userId)
	}

	assetSize, err := Get4byte(buf)
	if err != nil {
		return err
	}
	if assetSize > 0 {
		ast, err := GetBytes(buf, int(assetSize))
		if err != nil {
			return err
		}
		p.Asset = &BBcAsset{IdLength:p.IdLength}
		p.Asset.Unpack(&ast)
	}

	return nil
}

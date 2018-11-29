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
	"encoding/binary"
	"fmt"
)

/*
This is the BBcRelation definition.

The BBcRelation holds the asset (by BBcAsset) and the relationship with the other transaction/asset (by BBcPointer).
Different from UTXO, state information or account-type information can be expressed by using this object.
If you want to include signature(s) according to the contents of BBcRelation object, BBcWitness should be included in the transaction object.

"AssetGroupId" distinguishes a type of asset, e.g., token-X, token-Y, Movie content, etc..
"Pointers" is a list of BBcPointers object. "Asset" is a BBcAsset object.

"IdLength" is not included in a packed data. It is for internal use only.
 */
type (
	BBcRelation struct {
		IdLength		int
		AssetGroupId	[]byte
		Pointers 		[]*BBcPointer
		Asset 			*BBcAsset
	}
)

// Output content of the object
func (p *BBcRelation) Stringer() string {
	ret := fmt.Sprintf("  asset_group_id: %x\n", p.AssetGroupId)
	if p.Pointers != nil {
		ret += fmt.Sprintf("  Pointers[]: %d\n", len(p.Pointers))
		for i := range p.Pointers {
			ret += fmt.Sprintf("   [%d]\n", i)
			ret += p.Pointers[i].Stringer()
		}
	} else {
		ret += fmt.Sprintf("  Pointers[]: None\n")
	}
	if p.Asset != nil {
		ret += p.Asset.Stringer()
	} else {
		ret += fmt.Sprintf("  Asset: None\n")
	}
	return ret
}

// Add essential information (assetGroupId and BBcAsset object) to the BBcRelation object
func (p *BBcRelation) Add(assetGroupId *[]byte, asset *BBcAsset) {
	if assetGroupId != nil {
		p.AssetGroupId = make([]byte, p.IdLength)
		copy(p.AssetGroupId, *assetGroupId)
	}
	if asset != nil {
		p.Asset = asset
		p.Asset.IdLength = p.IdLength
	}
}

// Add BBcPointer object in the object
func (p *BBcRelation) AddPointer(pointer *BBcPointer) {
	pointer.IdLength = p.IdLength
	p.Pointers = append(p.Pointers, pointer)
}

// Pack BBcRelation object in binary data
func (p *BBcRelation) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	PutBigInt(buf, &p.AssetGroupId, p.IdLength)

	Put2byte(buf, uint16(len(p.Pointers)))
	for _, p := range p.Pointers {
		dat, er := p.Pack()
		if er != nil {
			return nil, er
		}
		Put2byte(buf, uint16(binary.Size(dat)))
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return nil, err
		}
	}
	if p.Asset != nil {
		ast, er := p.Asset.Pack()
		if er != nil {
			return nil, er
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

// Unpack binary data to BBcRelation object
func (p *BBcRelation) Unpack(dat *[]byte) error {
	var err error
	buf := bytes.NewBuffer(*dat)

	p.AssetGroupId, err = GetBigInt(buf)
	if err != nil {
		return err
	}

	numPointers, err := Get2byte(buf)
	for i := 0; i < int(numPointers); i++ {
		plen, err := Get2byte(buf)
		if err != nil {
			return err
		}
		ptr, err := GetBytes(buf, int(plen))
		pointer := BBcPointer{IdLength:p.IdLength}
		pointer.Unpack(&ptr)
		p.Pointers = append(p.Pointers, &pointer)
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

package bbclib

import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type (
	BBcSignature struct {
		KeyType 		uint32
		Pubkey 			[]byte
		PubkeyLen 		uint32
		Signature 		[]byte
		SignatureLen 	uint32
	}
)


func (p *BBcSignature) Stringer() string {
	if p.KeyType == KeyType_NOT_INITIALIZED {
		return "  Not initialized\n"
	}
	ret :=  fmt.Sprintf("  key_type: %d\n", p.KeyType)
	ret +=  fmt.Sprintf("  signature: %x\n", p.Signature)
	ret +=  fmt.Sprintf("  pubkey: %x\n", p.Pubkey)
	return ret
}

func (p *BBcSignature) SetPublicKey(keyType uint32, pubkey *[]byte) {
	p.KeyType = keyType
	p.Pubkey = *pubkey
	p.PubkeyLen = uint32(len(p.Pubkey) * 8)
}

func (p *BBcSignature) SetPublicKeyByKeypair(keypair *KeyPair) {
	p.KeyType = uint32(keypair.CurveType)
	p.Pubkey = keypair.Pubkey
	p.PubkeyLen = uint32(len(p.Pubkey) * 8)
}

func (p *BBcSignature) SetSignature(sig *[]byte) {
	p.Signature = *sig
	p.SignatureLen = uint32(len(p.Signature) * 8)
}


func (p *BBcSignature) Verify(digest []byte) bool {
	return VerifyBBcSignature(digest, p)
}

func (p *BBcSignature) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	Put4byte(buf, p.KeyType)
	if p.KeyType == KeyType_NOT_INITIALIZED {
		return buf.Bytes(), nil
	}

	Put4byte(buf, p.PubkeyLen)
	if err := binary.Write(buf, binary.LittleEndian, p.Pubkey); err != nil {
		return nil, err
	}

	Put4byte(buf, p.SignatureLen)
	if err := binary.Write(buf, binary.LittleEndian, p.Signature); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


func (p *BBcSignature) Unpack(dat []byte) error {
	var err error
	buf := bytes.NewBuffer(dat)

	keyType, err := Get4byte(buf)
	if err != nil {
		return err
	}
	if keyType == 0 {
		return nil
	}
	p.KeyType = uint32(keyType)

	p.PubkeyLen, err = Get4byte(buf)
	if err != nil {
		return err
	}
	p.Pubkey = make([]byte, int(p.PubkeyLen/8))
	p.Pubkey, err =  GetBytes(buf, int(p.PubkeyLen/8))

	p.SignatureLen, err = Get4byte(buf)
	if err != nil {
		return err
	}
	p.Signature = make([]byte, int(p.SignatureLen/8))
	p.Signature, err =  GetBytes(buf, int(p.SignatureLen/8))

	return nil
}

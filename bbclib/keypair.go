package bbclib

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lbbcsig
#include "libbbcsig.h"
*/
import "C"
import (
	"unsafe"
)

type (
	KeyPair struct {
		CurveType	int
		Pubkey		[]byte
		Privkey		[]byte
	}
)

const (
	KeyType_NOT_INITIALIZED = 0
	KeyType_ECDSA_SECP256k1 = 1
	KeyType_ECDSA_P256v1 = 2

	defaultCompressionMode = 4
)

func GenerateKeypair(curveType int, compressionMode int) KeyPair {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)
	var lenPubkey, lenPrivkey C.int
	C.generate_keypair(C.int(curveType), C.uint8_t(compressionMode), &lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		               &lenPrivkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	return KeyPair{CurveType: curveType, Pubkey: pubkey[:lenPubkey], Privkey: privkey[:lenPrivkey]}
}


func (k *KeyPair) ConvertFromPem(pem string, compressionMode int) {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)
	pemstr := ([]byte)(pem)

	var lenPubkey, lenPrivkey C.int
	C.convert_from_pem((*C.char)(unsafe.Pointer(&pemstr[0])), (C.uint8_t)(compressionMode),
		&lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		&lenPrivkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	k.Pubkey = pubkey[:lenPubkey]
	k.Privkey = pubkey[:lenPrivkey]
}


func (k *KeyPair) Sign(digest []byte) []byte {
	sig_r := make([]byte, 100)
	sig_s := make([]byte, 100)
	var len_sig_r, len_sig_s C.uint
	C.sign(C.int(k.CurveType), C.int(len(k.Privkey)), (*C.uint8_t)(unsafe.Pointer(&k.Privkey[0])),
		   C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		   (*C.uint8_t)(unsafe.Pointer(&sig_r[0])), (*C.uint8_t)(unsafe.Pointer(&sig_s[0])),
		   (*C.uint)(&len_sig_r), (*C.uint)(&len_sig_s))

	if len_sig_r < 32 {
		zeros := make([]byte, 32-len_sig_r)
		for i := range zeros {
			zeros[i] = 0
		}
		sig_r = append(zeros, sig_r[:32]...)
	}
	if len_sig_s < 32 {
		zeros := make([]byte, 32-len_sig_s)
		for i := range zeros {
			zeros[i] = 0
		}
		sig_s = append(zeros, sig_s[:32]...)
	}
	sig := make([]byte, len_sig_r+len_sig_s)
	sig = append(sig_r[:32], sig_s[:32]...)

	return sig
}


func (k *KeyPair) Verify(digest []byte, sig []byte) bool {
	result := C.verify(C.int(k.CurveType), C.int(len(k.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&(k.Pubkey[0]))),
		               C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
			           C.int(len(sig)), (*C.uint8_t)(unsafe.Pointer(&sig[0])))
	return result == 1
}


func VerifyBBcSignature(digest []byte, sig *BBcSignature) bool {
	result := C.verify(C.int(sig.KeyType), C.int(len(sig.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&sig.Pubkey[0])),
		C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		C.int(len(sig.Signature)), (*C.uint8_t)(unsafe.Pointer(&sig.Signature[0])))
	return result == 1
}

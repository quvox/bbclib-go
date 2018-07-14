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
		Curvetype	int
		Pubkey		[]byte
		Privkey		[]byte
	}
)


func GenerateKeypair(curvetype int) KeyPair {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)
	var len_pubkey, len_privkey C.int
	C.generate_keypair(C.int(curvetype), 0, &len_pubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		               &len_privkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	return KeyPair{Curvetype: curvetype, Pubkey: pubkey[:len_pubkey], Privkey: privkey[:len_privkey]}
}


func (k *KeyPair) ConvertFromPem(pem string) {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)
	pemstr := ([]byte)(pem)
	var len_pubkey, len_privkey C.int
	C.convert_from_pem(C.int(k.Curvetype), (*C.char)(unsafe.Pointer(&pemstr[0])), (C.uint8_t)(0),
		&len_pubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		&len_privkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	k.Pubkey = pubkey[:len_pubkey]
	k.Privkey = pubkey[:len_privkey]
}


func (k *KeyPair) Sign(digest []byte) []byte {
	sig_r := make([]byte, 100)
	sig_s := make([]byte, 100)
	var len_sig_r, len_sig_s C.uint
	C.sign(C.int(k.Curvetype), C.int(len(k.Privkey)), (*C.uint8_t)(unsafe.Pointer(&k.Privkey[0])),
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
	result := C.verify(C.int(k.Curvetype), C.int(len(k.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&(k.Pubkey[0]))),
		               C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
			           C.int(len(sig)), (*C.uint8_t)(unsafe.Pointer(&sig[0])))
	return result == 1
}


func VerifyBBcSignature(digest []byte, sig *BBcSignature) bool {
	result := C.verify(C.int(sig.Key_type), C.int(len(sig.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&sig.Pubkey[0])),
		C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		C.int(len(sig.Signature)), (*C.uint8_t)(unsafe.Pointer(&sig.Signature[0])))
	return result == 1
}

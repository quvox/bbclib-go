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


func (k KeyPair) Sign(digest [32]byte) []byte {
	privkey := make([]byte, 100)
	sig_r := make([]byte, 100)
	sig_s := make([]byte, 100)
	var len_sig_r, len_sig_s C.uint
	C.sign(C.int(k.Curvetype), C.int(32), (*C.uint8_t)(unsafe.Pointer(&privkey[0])),
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

	return sig[:64]
}


func (k KeyPair) Verify(digest [32]byte, sig []byte) bool {
	result := C.verify(C.int(k.Curvetype), C.int(len(k.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&(k.Pubkey))),
		               C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
			           C.int(len(sig)), (*C.uint8_t)(unsafe.Pointer(&sig[0])))
	return result == 1
}


func (k KeyPair) ConvertFromPem(pem string) {
	var len_pubkey, len_privkey C.int
	C.convert_from_pem(C.int(k.Curvetype), (*C.char)(unsafe.Pointer(&pem)), (C.uint8_t)(0),
					   &len_pubkey, (*C.uint8_t)(unsafe.Pointer(&(k.Pubkey))),
					   	&len_privkey, (*C.uint8_t)(unsafe.Pointer(&(k.Privkey))))
}


func VerifyBBcSignature(digest [32]byte, sig BBcSignature) bool {
	pubkey := sig.Pubkey
	signature := sig.Signature
	result := C.verify(C.int(sig.Key_type), C.int(len(pubkey)), (*C.uint8_t)(unsafe.Pointer(&(pubkey))),
		C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		C.int(len(signature)), (*C.uint8_t)(unsafe.Pointer(&signature[0])))
	return result == 1
}

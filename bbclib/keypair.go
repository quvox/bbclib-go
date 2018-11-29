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

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lbbcsig
#include "libbbcsig.h"
*/
import "C"
import "unsafe"

/*
This is the KeyPair definition.

A KeyPair object hold a pair of private key and public key.
This object includes functions for sign and verify a signature. The sign/verify functions is realized by "libbbcsig".
 */
type (
	KeyPair struct {
		CurveType	int
		Pubkey		[]byte
		Privkey		[]byte
	}
)

// Supported ECC curve type is SECP256k1 and Prime-256v1.
const (
	KeyType_NOT_INITIALIZED = 0
	KeyType_ECDSA_SECP256k1 = 1
	KeyType_ECDSA_P256v1 = 2

	defaultCompressionMode = 4
)

// Key pair generator
func GenerateKeypair(curveType int, compressionMode int) KeyPair {
	pubkey := make([]byte, 100)
	privkey := make([]byte, 100)
	var lenPubkey, lenPrivkey C.int
	C.generate_keypair(C.int(curveType), C.uint8_t(compressionMode), &lenPubkey, (*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		               &lenPrivkey, (*C.uint8_t)(unsafe.Pointer(&privkey[0])))
	return KeyPair{CurveType: curveType, Pubkey: pubkey[:lenPubkey], Privkey: privkey[:lenPrivkey]}
}

// Output PEM formatted public key
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

// Sign to a given digest
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

// Verify a given digest with signature
func (k *KeyPair) Verify(digest []byte, sig []byte) bool {
	result := C.verify(C.int(k.CurveType), C.int(len(k.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&(k.Pubkey[0]))),
		               C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
			           C.int(len(sig)), (*C.uint8_t)(unsafe.Pointer(&sig[0])))
	return result == 1
}

// Verify a given digest with BBcSignature object
func VerifyBBcSignature(digest []byte, sig *BBcSignature) bool {
	result := C.verify(C.int(sig.KeyType), C.int(len(sig.Pubkey)), (*C.uint8_t)(unsafe.Pointer(&sig.Pubkey[0])),
		C.int(len(digest)), (*C.uint8_t)(unsafe.Pointer(&digest[0])),
		C.int(len(sig.Signature)), (*C.uint8_t)(unsafe.Pointer(&sig.Signature[0])))
	return result == 1
}

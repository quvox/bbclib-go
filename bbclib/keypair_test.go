package bbclib

import (
	"testing"
	"crypto/sha256"
	"fmt"
)


func TestGenerateKeypair(t *testing.T) {
	for curvetype := 1; curvetype < 3; curvetype++ {
		t.Run("curvetype", func(t *testing.T) {
			keypair := GenerateKeypair(curvetype)
			if len(keypair.Pubkey) != 65 {
				t.Fatal("fail to generate keypair")
			}
			t.Logf("keypair: %v", keypair)
		})
	}
}


func TestKeyPair_Sign_and_Verify(t *testing.T) {
	digest := sha256.Sum256([]byte("aaaaaaaaaaa"))
	digest2 := sha256.Sum256([]byte("bbbbbbbbbbbbb"))
	fmt.Printf("SHA-256 digest : %x\n", digest)
	fmt.Printf("SHA-256 digest2: %x\n", digest2)

	for curvetype := 1; curvetype < 3; curvetype++ {
		keypair := GenerateKeypair(curvetype)
		keypair2 := GenerateKeypair(curvetype)
		t.Run("curvetype", func(t *testing.T) {
			t.Logf("Curvetype = %d", curvetype)
			if len(keypair.Pubkey) != 65 {
				t.Fatal("fail to generate keypair")
			}

			t.Logf("privkey   : %x\n", keypair.Privkey)
			sig1 := keypair.Sign(digest[:])
			t.Logf("signature : %x\n", sig1)
			if len(sig1) != 64 {
				t.Fatal("fail to sign")
			}
			result := keypair.Verify(digest[:], sig1)
			if ! result {
				t.Fatal("fail to verify")
			}
			t.Log("[sig1] Verify succeeded")
			result = keypair.Verify(digest2[:], sig1)
			if result {
				t.Fatal("[invalid digest] Verify returns true but not correct...")
			}
			t.Log("[invalid digest] Verify failed as expected")


			t.Logf("privkey2  : %x\n", keypair2.Privkey)
			sig2 := keypair2.Sign(digest[:])
			t.Logf("signature2: %x\n", sig2)
			if len(sig2) != 64 {
				t.Fatal("fail to sign")
			}
			result = keypair2.Verify(digest[:], sig2)
			if ! result {
				t.Fatal("fail to verify")
			}
			t.Log("[sig2] Verify succeeded")
			result = keypair2.Verify(digest2[:], sig2)
			if result {
				t.Fatal("[invalid digest] Verify returns true but not correct...")
			}
			t.Log("[invalid digest] Verify failed as expected")


			result = keypair2.Verify(digest[:], sig1)
			if result {
				t.Fatal("[swap] Verify returns true but not correct...")
			}
			t.Log("[swap] Verify failed as expected")
			result = keypair.Verify(digest[:], sig2)
			if result {
				t.Fatal("[swap] Verify returns true but not correct...")
			}
			t.Log("[swap] Verify failed as expected")
		})
	}
}


func TestVerifyBBcSignature(t *testing.T) {
	digest := sha256.Sum256([]byte("aaaaaaaaaaa"))

	for curvetype := 1; curvetype < 3; curvetype++ {
		t.Run("curvetype", func(t *testing.T) {
			keypair := GenerateKeypair(curvetype)
			sig := BBcSignature{
				Format_type: FORMAT_BSON,
				Key_type:    curvetype,
				Signature:   keypair.Sign(digest[:]),
				Pubkey:      keypair.Pubkey,
			}
			result1 := keypair.Verify(digest[:], sig.Signature)
			if ! result1 {
				t.Fatal("fail to verify")
			}
			t.Log("Verify succeeded")

			result := VerifyBBcSignature(digest[:], &sig)
			if ! result {
				t.Fatal("fail to verify")
			}
			t.Log("Verify succeeded")

			keypair2 := GenerateKeypair(curvetype)
			sig2 := BBcSignature{
				Format_type: FORMAT_BSON,
				Key_type:    curvetype,
				Signature:   keypair2.Sign(digest[:]),
				Pubkey:      keypair.Pubkey,
			}
			result = VerifyBBcSignature(digest[:], &sig2)
			if result {
				t.Fatal("Verify returns true but not correct...")
			}
			t.Log("Verify failed as expected")
		})
	}
}


func TestKeyPair_ConvertFromPem(t *testing.T) {
	pem := "-----BEGIN EC PRIVATE KEY-----\nMHQCAQEEIIMVMPKLJqivgRDpRDaWJCOnob6s/+t4MdoFN/8PVkNSoAcGBSuBBAAK\noUQDQgAE/k1ZM/Ker1+N0+Lg5za0sJZeSAAeYwDEWnkgnkCynErs74G/tAnu/lcu\nk8kzAivYm8mitIpJJw1OdjCDJI457g==\n-----END EC PRIVATE KEY-----\n"
	keypair := KeyPair{Curvetype: 1}
	keypair.ConvertFromPem(pem)
	t.Logf("keypair: %v", keypair)

	if len(keypair.Privkey) != 32 {
		t.Fatal("failed to read private key in pem format")
	}
	t.Logf("private key: %x", keypair.Privkey)
}
package bbclib

import (
	"bytes"
	"testing"
)

var (
	idLength = 32
)

func TestWitnessPackUnpack(t *testing.T) {
	t.Run("simple creation (string asset)", func(t *testing.T) {
		txobj := BBcTransaction{IdLength:idLength}
		obj := BBcWitness{IdLength:idLength, Transaction:&txobj}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", idLength)
		u2 := GetIdentifierWithTimestamp("user2", idLength)

		obj.AddWitness(&u1)
		obj.AddWitness(&u2)

		sig := BBcSignature{}
		obj.AddSignature(&u1, &sig)
		obj.AddSignature(&u2, &sig)

		t.Log("---------------witness-----------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcWitness{IdLength:idLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.UserIds[0], obj2.UserIds[0]) != 0 || bytes.Compare(obj.UserIds[1], obj2.UserIds[1]) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}


func TestWitnessInvalidAccess(t *testing.T) {
	t.Run("no transaction", func(t *testing.T) {
		obj := BBcWitness{IdLength:idLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", idLength)

		err := obj.AddWitness(&u1)
		if err == nil {
			t.Fatal("Should fail because no Transaction is set")
		}

		sig := BBcSignature{}
		err = obj.AddSignature(&u1, &sig)
		if err == nil {
			t.Fatal("Should fail because no Transaction is set")
		}
	})
}

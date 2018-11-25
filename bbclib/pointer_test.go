package bbclib

import (
	"bytes"
	"testing"
)

var (
	idLength = 32
)


func TestPointerPackUnpack(t *testing.T) {
	t.Run("simple creation", func(t *testing.T) {
		obj := BBcPointer{IdLength:idLength}
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", idLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", idLength)
		obj.Add(&txid1, &asid1)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcPointer{IdLength:idLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.TransactionId, obj2.TransactionId) != 0 || bytes.Compare(obj.AssetId, obj2.AssetId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (asset_id is nil)", func(t *testing.T) {
		obj := BBcPointer{IdLength:idLength}
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", idLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", idLength)
		obj.Add(&txid1, &asid1)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcPointer{IdLength:idLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.TransactionId, obj2.TransactionId) != 0 || bytes.Compare(obj.AssetId, obj2.AssetId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
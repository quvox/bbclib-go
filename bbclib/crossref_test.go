package bbclib

import (
	"bytes"
	"testing"
)

var (
	IdLength = 32
)


func TestCrossRefPackUnpack(t *testing.T) {
	t.Run("simple creation", func(t *testing.T) {
		obj := BBcCrossRef{IdLength:IdLength}
		dom := GetIdentifier("dummy domain", IdLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", IdLength)
		obj.Add(&dom, &dummyTxid)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcCrossRef{IdLength:IdLength}
		obj2.Unpack(dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.DomainId, obj2.DomainId) != 0 || bytes.Compare(obj.TransactionId, obj2.TransactionId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
package bbclib

import (
	"bytes"
	"testing"
)

var (
	idLength = 32
)

func TestEventPackUnpack(t *testing.T) {
	t.Run("simple creation (string asset)", func(t *testing.T) {
		ast := BBcAsset{IdLength:idLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", idLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		obj := BBcEvent{IdLength:idLength}
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", idLength)
		obj.Add(&assetgroup, &ast)
		obj.AddReferenceIndex(1)
		obj.AddReferenceIndex(2)

		u2 := GetIdentifierWithTimestamp("user2", idLength)
		u3 := GetIdentifierWithTimestamp("user3", idLength)
		obj.AddMandatoryApprover(&u1)
		obj.AddMandatoryApprover(&u2)
		obj.AddOptionApprover(&u3)
		obj.AddOptionParams(2, 2)

		t.Log("---------------Event-----------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcEvent{IdLength:idLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetGroupId, obj2.AssetGroupId) != 0 || bytes.Compare(obj.Asset.AssetId, obj2.Asset.AssetId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (no approvers)", func(t *testing.T) {
		ast := BBcAsset{IdLength:idLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", idLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		obj := BBcEvent{IdLength:idLength}
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", idLength)
		obj.Add(&assetgroup, &ast)

		t.Log("---------------Event-----------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcEvent{IdLength:idLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetGroupId, obj2.AssetGroupId) != 0 || bytes.Compare(obj.Asset.AssetId, obj2.Asset.AssetId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
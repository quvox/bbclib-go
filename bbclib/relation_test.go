package bbclib

import (
	"bytes"
	"testing"
)

var (
	idLength = 32
)

func TestRelationPackUnpack(t *testing.T) {
	t.Run("simple creation (string asset)", func(t *testing.T) {
		obj := BBcRelation{IdLength:idLength}
		ptr1 := BBcPointer{}
		ptr2 := BBcPointer{}
		ast := BBcAsset{}

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", idLength)
		obj.Add(&assetgroup, &ast)
		obj.AddPointer(&ptr1)
		obj.AddPointer(&ptr2)

		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", idLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", idLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", idLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", idLength)
		ptr1.Add(&txid1, &asid1)
		ptr2.Add(&txid2, nil)

		t.Log("---------------Relation-----------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcRelation{IdLength:idLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.AssetGroupId, obj2.AssetGroupId) != 0 || bytes.Compare(obj.Asset.AssetId, obj2.Asset.AssetId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (no pointer, msgpack asset)", func(t *testing.T) {
		ast := BBcAsset{IdLength:idLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", idLength)
		ast.Add(&u1)
		ast.AddBodyObject(map[int]string{1:"aaa", 2:"bbb", 10:"asdfasdfasf;lakj;lkj;"})

		obj := BBcRelation{IdLength:idLength}
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", idLength)
		obj.Add(&assetgroup, &ast)
		t.Log("---------------Relation-----------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcRelation{IdLength:idLength}
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
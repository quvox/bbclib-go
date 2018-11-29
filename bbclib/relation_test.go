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

import (
	"bytes"
	"testing"
)


func TestRelationPackUnpack(t *testing.T) {
	t.Run("simple creation (string asset)", func(t *testing.T) {
		obj := BBcRelation{IdLength:defaultIdLength}
		ptr1 := BBcPointer{}
		ptr2 := BBcPointer{}
		ast := BBcAsset{}

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIdLength)
		obj.Add(&assetgroup, &ast)
		obj.AddPointer(&ptr1)
		obj.AddPointer(&ptr2)

		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIdLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")
		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", defaultIdLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", defaultIdLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", defaultIdLength)
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

		obj2 := BBcRelation{IdLength:defaultIdLength}
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
		ast := BBcAsset{IdLength:defaultIdLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIdLength)
		ast.Add(&u1)
		ast.AddBodyObject(map[int]string{1:"aaa", 2:"bbb", 10:"asdfasdfasf;lakj;lkj;"})

		obj := BBcRelation{IdLength:defaultIdLength}
		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", defaultIdLength)
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

		obj2 := BBcRelation{IdLength:defaultIdLength}
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
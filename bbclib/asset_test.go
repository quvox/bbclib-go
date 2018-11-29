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
	"io/ioutil"
	"testing"
)


func TestAssetPackUnpack(t *testing.T) {
	t.Run("simple creation (string)", func(t *testing.T) {
		obj := BBcAsset{IdLength:defaultIdLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIdLength)
		obj.Add(&u1)
		obj.AddBodyString("testString12345XXX")
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcAsset{IdLength:defaultIdLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.UserId, obj2.UserId) != 0 || bytes.Compare(obj.AssetId, obj2.AssetId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (string with file)", func(t *testing.T) {
		obj := BBcAsset{IdLength:defaultIdLength}
		u1 := GetIdentifier("user2_789abcdef0123456789abcdef0", defaultIdLength)
		obj.Add(&u1)
		obj.AddBodyString("test string xxx")
		filedat, err := ioutil.ReadFile("./asset_test.go")
		obj.AddFile(&filedat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcAsset{IdLength:defaultIdLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.UserId, obj2.UserId) != 0 || bytes.Compare(obj.AssetId, obj2.AssetId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (msgpack)", func(t *testing.T) {
		obj := BBcAsset{IdLength:defaultIdLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIdLength)
		obj.Add(&u1)
		obj.AddBodyObject(map[int]string{1:"aaa", 2:"bbb"})

		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj.IdLength)
		t.Logf("%v", obj.Stringer())
		body, err := obj.GetBodyObject()
		t.Logf("body_object: %v",body)
		t.Log("--------------------------------------")

		dat, err := obj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcAsset{IdLength:defaultIdLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		body2, err := obj2.GetBodyObject()
		t.Logf("body_object: %v",body2)
		t.Log("--------------------------------------")

		if bytes.Compare(obj.UserId, obj2.UserId) != 0 {
			t.Fatal("Not recovered correctly...1")
		}
		obj2.Digest()
		if bytes.Compare(obj.AssetId, obj2.AssetId) != 0 {
			t.Logf("obj : %x\n", obj.AssetId)
			t.Logf("obj2: %x\n", obj2.AssetId)
			t.Fatal("Not recovered correctly...2")
		}

	})
}
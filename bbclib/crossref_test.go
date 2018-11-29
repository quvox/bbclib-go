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


func TestCrossRefPackUnpack(t *testing.T) {
	t.Run("simple creation", func(t *testing.T) {
		obj := BBcCrossRef{IdLength:defaultIdLength}
		dom := GetIdentifier("dummy domain", defaultIdLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", defaultIdLength)
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

		obj2 := BBcCrossRef{IdLength:defaultIdLength}
		obj2.Unpack(&dat)
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IdLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(obj.DomainId, obj2.DomainId) != 0 || bytes.Compare(obj.TransactionId, obj2.TransactionId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
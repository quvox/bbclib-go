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


func TestWitnessPackUnpack(t *testing.T) {
	t.Run("simple creation (string asset)", func(t *testing.T) {
		txobj := BBcTransaction{IdLength:defaultIdLength}
		obj := BBcWitness{IdLength:defaultIdLength, Transaction:&txobj}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIdLength)
		u2 := GetIdentifierWithTimestamp("user2", defaultIdLength)

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

		obj2 := BBcWitness{IdLength:defaultIdLength}
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
		obj := BBcWitness{IdLength:defaultIdLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", defaultIdLength)

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

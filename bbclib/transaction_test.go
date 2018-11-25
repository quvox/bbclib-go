package bbclib

import (
	"bytes"
	"testing"
	"time"
)

var (
	IdLength = 8
)

func TestTransactionPackUnpack(t *testing.T) {

	t.Run("simple creation (with relation)", func(t *testing.T) {
		keypair := GenerateKeypair(KeyType_ECDSA_P256v1, defaultCompressionMode)
		txobj := BBcTransaction{Version:1, Timestamp:time.Now().UnixNano(), IdLength:IdLength}
		rtn := BBcRelation{}
		txobj.AddRelation(&rtn)
		wit := BBcWitness{}
		txobj.AddWitness(&wit)
		crs := BBcCrossRef{}
		txobj.AddCrossRef(&crs)

		ast := BBcAsset{}
		ptr1 := BBcPointer{}
		ptr2 := BBcPointer{}

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", IdLength)
		rtn.Add(&assetgroup, &ast)
		rtn.AddPointer(&ptr1)
		rtn.AddPointer(&ptr2)

		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", IdLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", IdLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", IdLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", IdLength)
		ptr1.Add(&txid1, &asid1)
		ptr2.Add(&txid2, nil)

		wit.AddWitness(&u1)
		u2 := GetIdentifierWithTimestamp("user2", IdLength)
		wit.AddWitness(&u2)

		dom := GetIdentifier("dummy domain", IdLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", IdLength)
		crs.Add(&dom, &dummyTxid)

		sig := BBcSignature{}
		sig.SetPublicKeyByKeypair(&keypair)
		signature, err := txobj.Sign(&keypair)
		sig.SetSignature(&signature)
		wit.AddSignature(&u1, &sig)

		sig2 := BBcSignature{}
		sig2.SetPublicKeyByKeypair(&keypair)
		signature2, err := txobj.Sign(&keypair)
		sig2.SetSignature(&signature2)
		wit.AddSignature(&u2, &sig2)

		t.Log("---------------transaction-----------------")
		t.Logf("%v", txobj.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		obj2.Digest()
		if bytes.Compare(txobj.Relations[0].Asset.AssetId, obj2.Relations[0].Asset.AssetId) != 0 ||
			bytes.Compare(txobj.TransactionId, obj2.TransactionId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	txobj2 := BBcTransaction{Version:1, Timestamp:time.Now().UnixNano(), IdLength:IdLength}
	t.Run("simple creation (with event)", func(t *testing.T) {
		keypair := GenerateKeypair(KeyType_ECDSA_P256v1, defaultCompressionMode)
		evt := BBcEvent{}
		txobj2.AddEvent(&evt)
		crs := BBcCrossRef{}
		txobj2.AddCrossRef(&crs)
		wit := BBcWitness{}
		txobj2.AddWitness(&wit)

		ast := BBcAsset{IdLength:IdLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", IdLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", IdLength)
		evt.Add(&assetgroup, &ast)

		u2 := GetIdentifierWithTimestamp("user2", IdLength)
		evt.AddMandatoryApprover(&u1)
		evt.AddMandatoryApprover(&u2)
		evt.AddOptionParams(2, 1)

		dom := GetIdentifier("dummy domain", IdLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", IdLength)
		crs.Add(&dom, &dummyTxid)

		wit.AddWitness(&u1)
		sig := BBcSignature{}
		sig.SetPublicKeyByKeypair(&keypair)
		signature, err := txobj2.Sign(&keypair)
		sig.SetSignature(&signature)
		wit.AddSignature(&u1, &sig)

		t.Log("---------------transaction-----------------")
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj2.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		obj2.Digest()
		if bytes.Compare(txobj2.Events[0].Asset.AssetId, obj2.Events[0].Asset.AssetId) != 0 ||
			bytes.Compare(txobj2.TransactionId, obj2.TransactionId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("simple creation (with event/reference)", func(t *testing.T) {
		keypair := GenerateKeypair(KeyType_ECDSA_P256v1, defaultCompressionMode)
		txobj3 := BBcTransaction{Version:1, Timestamp:time.Now().UnixNano(), IdLength:IdLength}
		evt := BBcEvent{}
		txobj3.AddEvent(&evt)
		ref := BBcReference{}
		txobj3.AddReference(&ref)
		crs := BBcCrossRef{}
		txobj3.AddCrossRef(&crs)

		ast := BBcAsset{IdLength:IdLength}
		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", IdLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", IdLength)
		evt.Add(&assetgroup, &ast)

		u2 := GetIdentifierWithTimestamp("user2", IdLength)
		evt.AddMandatoryApprover(&u1)
		evt.AddMandatoryApprover(&u2)
		evt.AddOptionParams(2, 1)

		ref.Add(&assetgroup, &txobj2, 0)
		ref.AddApprover(&u1)
		ref.AddApprover(&u2)

		dom := GetIdentifier("dummy domain", IdLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", IdLength)
		crs.Add(&dom, &dummyTxid)

		sig := BBcSignature{}
		sig.SetPublicKeyByKeypair(&keypair)
		signature, err := txobj2.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig.SetSignature(&signature)
		ref.AddSignature(&u1, &sig)

		sig2 := BBcSignature{}
		sig2.SetPublicKeyByKeypair(&keypair)
		signature2, err := txobj2.Sign(&keypair)
		if err != nil {
			t.Fatal(err)
		}
		sig2.SetSignature(&signature2)
		ref.AddSignature(&u2, &sig2)

		t.Log("---------------transaction-----------------")
		t.Logf("%v", txobj3.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj3.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		t.Log("---------------transaction-----------------")
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		obj2.Digest()
		if bytes.Compare(txobj3.Events[0].Asset.AssetId, obj2.Events[0].Asset.AssetId) != 0 ||
			bytes.Compare(txobj3.TransactionId, obj2.TransactionId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
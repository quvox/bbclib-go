package bbclib

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

var (
	idLength = 32
	txdataType0 = "000001000000cbeefc5b000000002000000000000100a100000020009cfa77b06efc3528c051f42d47a84e71d0f75056ae4542146b3f73c18169c9d000007900000020007a7374f7d35f5eb37a42b4551c0d98268988bdfd3084bccbdfd65c587c596d4620001e5f3cf4588e64234d88fed3e87f0fff3580c03a1dab2f55d42ec004c04cb1692000fc982c9463a8e29ffd46331fae974cb43d8f76822c3e9b92f230ad95138061e70000000000000b007472616e73666572726564010026000000010020001e5f3cf4588e64234d88fed3e87f0fff3580c03a1dab2f55d42ec004c04cb1690000000001008d0000000200000008020000042159bcbd47bcb8cb8183f6a07b06d212e3f21807d534278d7b7dbb1b93f1f0b5b956a80a0032cb411712d8bbd3beefe1a71c52866681676c30be7e6f7c556b1b000200002d426b1c5cca56e47eef0162604e00203bc4b21a544c4ce35ec67a3aee81d798f6201be300fe4a10aea61b625a07020a19e9bfa786c0aaf560c950b8e0bcec13"
	txdataType1 = "1000789c6364606038fdee4f3490625000110c8c0c0bc1ec39bfca37e4fd31d53810f845d77d855fe185ef0161eb5c9d44b2ed8b0f36669ebcc0c050095657555cf2fd727cdce62aa72da132bc33d43a3bf6fe3568d973fafeb598889ac85c370506b9789b2f117d29cabe1dff2ebfa8e7ff6fda70c04a76b57ee815bd032c077c36662a30fc99a1332579c5a3f97fdd8ce5d74df7d962db5fd6a463377bd22783b553851b129f831dc6c0cd50529498579c965a54949ac2c8a006762c61e3219eea05924c40cc0124581423f7ec75dfb3e37463f3b705d56c97841e7f9260bf6aa2de5b5dbb5b7af2c70f5b7786ade062303aed282e7463f7e57def1f2e97096a4b6b4ccf31d857975f139a2d0d3249d7295b26e654d893baf78c49097e0c0ad647364985f8f83c8e3b5665f5aef1fa8c6f0ad28f19fe7909ac5b269d14c5cec425f972fff2b603abbe269c0cd8f160cf1b61008dcf976f"
	assetGroupIdInTx = "9cfa77b06efc3528c051f42d47a84e71d0f75056ae4542146b3f73c18169c9d0"
	txid = "6d187c7ff825e46d4c0f258d08264a070f0909314c854671a15d8dbf6c983f19"
)



func TestBBcLibSerializeDeserialze(t *testing.T) {
	txobj := BBcTransaction{Version:1, Timestamp:time.Now().UnixNano(), IdLength:idLength}
	t.Run("simple serialize and deserialize", func(t *testing.T) {
		keypair := GenerateKeypair(KeyType_ECDSA_P256v1, defaultCompressionMode)
		rtn := BBcRelation{}
		txobj.AddRelation(&rtn)
		wit := BBcWitness{}
		txobj.AddWitness(&wit)
		crs := BBcCrossRef{}
		txobj.AddCrossRef(&crs)

		ast := BBcAsset{}
		ptr1 := BBcPointer{}
		ptr2 := BBcPointer{}

		assetgroup := GetIdentifier("asset_group_id1,,,,,,,", idLength)
		rtn.Add(&assetgroup, &ast)
		rtn.AddPointer(&ptr1)
		rtn.AddPointer(&ptr2)

		u1 := GetIdentifier("user1_789abcdef0123456789abcdef0", idLength)
		ast.Add(&u1)
		ast.AddBodyString("testString12345XXX")

		txid1 := GetIdentifier("0123456789abcdef0123456789abcdef", idLength)
		txid2 := GetIdentifierWithTimestamp("asdfauflkajethb;:a", idLength)
		asid1 := GetIdentifier("123456789abcdef0123456789abcdef0", idLength)
		ptr1.Add(&txid1, &asid1)
		ptr2.Add(&txid2, nil)

		wit.AddWitness(&u1)
		u2 := GetIdentifierWithTimestamp("user2", idLength)
		wit.AddWitness(&u2)

		dom := GetIdentifier("dummy domain", idLength)
		dummyTxid := GetIdentifierWithTimestamp("dummytxid", idLength)
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

		dat, err := Serialize(&txobj, 0)
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Serialized data: %x", dat)
		t.Logf("Serialized data size: %d", len(dat))

		txobj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj2.IdLength)
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(txobj.TransactionId, txobj2.TransactionId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("serialize and deserialize with zlib", func(t *testing.T) {
		dat, err := Serialize(&txobj, 0x0010)
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Serialized data: %x", dat)
		t.Logf("Serialized data size: %d", len(dat))

		txobj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj2.IdLength)
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(txobj.TransactionId, txobj2.TransactionId) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("deserialize txdata genarated (type0) by python bbclib", func(t *testing.T) {
		dat, err := hex.DecodeString(txdataType0)
		txobj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj2.IdLength)
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		txid_orig, err := hex.DecodeString(txid)
		if bytes.Compare(txobj2.TransactionId, txid_orig) != 0 {
			t.Fatal("Not recovered correctly...1")
		}
		asgid_org, err := hex.DecodeString(assetGroupIdInTx)
		if bytes.Compare(txobj2.Relations[0].AssetGroupId, asgid_org) != 0 {
			t.Fatal("Not recovered correctly...2")
		}
	})

	t.Run("deserialize txdata genarated (type0x0010) by python bbclib", func(t *testing.T) {
		dat, err := hex.DecodeString(txdataType1)
		txobj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj2.IdLength)
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		txid_orig, err := hex.DecodeString(txid)
		if bytes.Compare(txobj2.TransactionId, txid_orig) != 0 {
			t.Fatal("Not recovered correctly...1")
		}
		asgid_org, err := hex.DecodeString(assetGroupIdInTx)
		if bytes.Compare(txobj2.Relations[0].AssetGroupId, asgid_org) != 0 {
			t.Fatal("Not recovered correctly...2")
		}
	})
}

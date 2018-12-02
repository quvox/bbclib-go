package bbclib

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"testing"
)

// Serialized data from Python bbclib.py
var (
	txdataType0      = "000001000000cbeefc5b000000002000000000000100a100000020009cfa77b06efc3528c051f42d47a84e71d0f75056ae4542146b3f73c18169c9d000007900000020007a7374f7d35f5eb37a42b4551c0d98268988bdfd3084bccbdfd65c587c596d4620001e5f3cf4588e64234d88fed3e87f0fff3580c03a1dab2f55d42ec004c04cb1692000fc982c9463a8e29ffd46331fae974cb43d8f76822c3e9b92f230ad95138061e70000000000000b007472616e73666572726564010026000000010020001e5f3cf4588e64234d88fed3e87f0fff3580c03a1dab2f55d42ec004c04cb1690000000001008d0000000200000008020000042159bcbd47bcb8cb8183f6a07b06d212e3f21807d534278d7b7dbb1b93f1f0b5b956a80a0032cb411712d8bbd3beefe1a71c52866681676c30be7e6f7c556b1b000200002d426b1c5cca56e47eef0162604e00203bc4b21a544c4ce35ec67a3aee81d798f6201be300fe4a10aea61b625a07020a19e9bfa786c0aaf560c950b8e0bcec13"
	txdataType1      = "1000789c6364606038fdee4f3490625000110c8c0c0bc1ec39bfca37e4fd31d53810f845d77d855fe185ef0161eb5c9d44b2ed8b0f36669ebcc0c050095657555cf2fd727cdce62aa72da132bc33d43a3bf6fe3568d973fafeb598889ac85c370506b9789b2f117d29cabe1dff2ebfa8e7ff6fda70c04a76b57ee815bd032c077c36662a30fc99a1332579c5a3f97fdd8ce5d74df7d962db5fd6a463377bd22783b553851b129f831dc6c0cd50529498579c965a54949ac2c8a006762c61e3219eea05924c40cc0124581423f7ec75dfb3e37463f3b705d56c97841e7f9260bf6aa2de5b5dbb5b7af2c70f5b7786ade062303aed282e7463f7e57def1f2e97096a4b6b4ccf31d857975f139a2d0d3249d7295b26e654d893baf78c49097e0c0ad647364985f8f83c8e3b5665f5aef1fa8c6f0ad28f19fe7909ac5b269d14c5cec425f972fff2b603abbe269c0cd8f160cf1b61008dcf976f"
	assetGroupIDInTx = "9cfa77b06efc3528c051f42d47a84e71d0f75056ae4542146b3f73c18169c9d0"
	txid             = "6d187c7ff825e46d4c0f258d08264a070f0909314c854671a15d8dbf6c983f19"
	txobj            *BBcTransaction
	txobj2           *BBcTransaction
	txobj3           *BBcTransaction
	assetGroupID     []byte
	u1               []byte
	u2               []byte
	keypair1         KeyPair
	keypair2         KeyPair
)

func TestBBcLibUtilitiesTx1(t *testing.T) {
	assetGroupID = GetIdentifierWithTimestamp("assetGroupID", defaultIDLength)
	u1 = GetIdentifierWithTimestamp("user1", defaultIDLength)
	u2 = GetIdentifierWithTimestamp("user2", defaultIDLength)
	keypair1 = GenerateKeypair(KeyTypeEcdsaP256v1, defaultCompressionMode)
	keypair2 = GenerateKeypair(KeyTypeEcdsaSECP256k1, defaultCompressionMode)

	t.Run("MakeTransaction and events", func(t *testing.T) {
		txobj = MakeTransaction(3, 0, true, 32)
		AddEventAssetBodyString(txobj, 0, &assetGroupID, &u1, "teststring!!!!!")
		txobj.Events[0].AddMandatoryApprover(&u1)
		filedat, _ := ioutil.ReadFile("./asset_test.go")
		AddEventAssetFile(txobj, 1, &assetGroupID, &u2, &filedat)
		txobj.Events[1].AddMandatoryApprover(&u2)
		datobj := map[int]string{1: "aaa", 2: "bbb", 10: "ccc"}
		AddEventAssetBodyObject(txobj, 2, &assetGroupID, &u1, &datobj)
		txobj.Events[2].AddMandatoryApprover(&u1)

		txobj.Witness.AddWitness(&u1)
		txobj.Witness.AddWitness(&u2)

		SignToTransaction(txobj, &u1, &keypair1)
		SignToTransaction(txobj, &u2, &keypair2)

		t.Log("-------------transaction--------------")
		t.Logf("%v", txobj.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		obj2.Digest()
		if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj.Events[0].Asset.AssetID, obj2.Events[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}

func TestBBcLibUtilitiesTx2(t *testing.T) {
	t.Run("MakeTransaction and events/reference", func(t *testing.T) {
		txobj2 = MakeTransaction(2, 0, true, 32)
		AddEventAssetBodyString(txobj2, 0, &assetGroupID, &u1, "teststring!!!!!")
		filedat, _ := ioutil.ReadFile("./crossref_test.go")
		AddEventAssetFile(txobj2, 1, &assetGroupID, &u2, &filedat)

		AddReference(txobj2, &assetGroupID, txobj, 0)
		txobj2.References[0].AddApprover(&u1)
		AddReference(txobj2, &assetGroupID, txobj, 1)
		txobj2.References[1].AddApprover(&u2)

		SignToTransaction(txobj2, &u1, &keypair1)
		SignToTransaction(txobj2, &u2, &keypair2)

		t.Log("-------------transaction--------------")
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj2.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		obj2.Digest()
		if result := obj2.Signatures[0].Verify(obj2.TransactionID); !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj2.Events[0].Asset.AssetID, obj2.Events[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj2.TransactionID, obj2.TransactionID) != 0 ||
			len(obj2.References[0].SigIndices) != 1 || len(obj2.References[1].SigIndices) != 1 {
			t.Fatal("Not recovered correctly...")
		}
	})
}
func TestBBcLibUtilitiesTx3(t *testing.T) {
	t.Run("MakeTransaction and relations", func(t *testing.T) {
		txobj3 = MakeTransaction(0, 3, true, 32)
		AddRelationAssetBodyString(txobj3, 0, &assetGroupID, &u1, "teststring!!!!!")
		filedat, _ := ioutil.ReadFile("./crossref_test.go")
		AddRelationAssetFile(txobj3, 1, &assetGroupID, &u2, &filedat)
		datobj := map[int]string{1: "aaa", 2: "bbb", 10: "ccc"}
		AddRelationAssetBodyObject(txobj3, 2, &assetGroupID, &u1, &datobj)

		datobj2 := map[int]string{1000: "lll", 22: "gggg", 100: "ddd"}
		rtn := MakeRelationWithAsset(&assetGroupID, &u2, "", &datobj2, nil, 32)
		txobj3.AddRelation(rtn)

		AddRelationPointer(txobj3, 0, &txobj.TransactionID, nil)
		AddRelationPointer(txobj3, 1, &txobj2.TransactionID, &txobj2.Events[0].Asset.AssetID)
		AddPointerInRelation(rtn, txobj, &txobj.Events[2].Asset.AssetID)

		txobj3.Witness.AddWitness(&u1)
		txobj3.Witness.AddWitness(&u2)

		SignToTransaction(txobj3, &u1, &keypair1)
		SignToTransaction(txobj3, &u2, &keypair2)

		t.Log("-------------transaction--------------")
		t.Logf("%v", txobj3.Stringer())
		t.Log("--------------------------------------")

		dat, err := txobj3.Pack()
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Packed data: %x", dat)

		obj2 := BBcTransaction{}
		obj2.Unpack(&dat)
		obj2.Digest()

		result, _ := obj2.VerifyAll()
		if !result {
			t.Fatal("Verification failed..")
		}

		if bytes.Compare(txobj3.Relations[0].Asset.AssetID, obj2.Relations[0].Asset.AssetID) != 0 ||
			bytes.Compare(txobj3.TransactionID, obj2.TransactionID) != 0 ||
			len(obj2.Witness.SigIndices) != 2 || len(obj2.Witness.SigIndices) != 2 {
			t.Fatal("Not recovered correctly...")
		}
	})
}

func TestBBcLibUtilitiesTx4(t *testing.T) {
	t.Run("simple serialize and deserialize", func(t *testing.T) {
		dat, err := Serialize(txobj, FormatPlain)
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Serialized data: %x", dat)
		t.Logf("Serialized data size: %d", len(dat))

		obj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IDLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(txobj.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})

	t.Run("serialize and deserialize with zlib", func(t *testing.T) {
		dat, err := Serialize(txobj2, FormatZlib)
		if err != nil {
			t.Fatalf("failed to serialize transaction object (%v)", err)
		}
		t.Logf("Serialized data: %x", dat)
		t.Logf("Serialized data size: %d", len(dat))

		obj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", obj2.IDLength)
		t.Logf("%v", obj2.Stringer())
		t.Log("--------------------------------------")

		if bytes.Compare(txobj2.TransactionID, obj2.TransactionID) != 0 {
			t.Fatal("Not recovered correctly...")
		}
	})
}

func TestBBcLibSerializeDeserialze(t *testing.T) {
	t.Run("deserialize txdata genarated (type0) by python bbclib", func(t *testing.T) {
		dat, _ := hex.DecodeString(txdataType0)
		txobj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj2.IDLength)
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		txidOrg, _ := hex.DecodeString(txid)
		if bytes.Compare(txobj2.TransactionID, txidOrg) != 0 {
			t.Fatal("Not recovered correctly...1")
		}
		asgidOrg, _ := hex.DecodeString(assetGroupIDInTx)
		if bytes.Compare(txobj2.Relations[0].AssetGroupID, asgidOrg) != 0 {
			t.Fatal("Not recovered correctly...2")
		}
	})

	t.Run("deserialize txdata genarated (type0x0010) by python bbclib", func(t *testing.T) {
		dat, _ := hex.DecodeString(txdataType1)
		txobj2, err := Deserialize(dat)
		if err != nil {
			t.Fatalf("failed to deserialize transaction data (%v)", err)
		}
		t.Log("--------------------------------------")
		t.Logf("id_length: %d", txobj2.IDLength)
		t.Logf("%v", txobj2.Stringer())
		t.Log("--------------------------------------")

		txidOrg, _ := hex.DecodeString(txid)
		if bytes.Compare(txobj2.TransactionID, txidOrg) != 0 {
			t.Fatal("Not recovered correctly...1")
		}
		asgidOrg, _ := hex.DecodeString(assetGroupIDInTx)
		if bytes.Compare(txobj2.Relations[0].AssetGroupID, asgidOrg) != 0 {
			t.Fatal("Not recovered correctly...2")
		}
	})
}

package bbclib

import (
	"bytes"
	"encoding/binary"
	"errors"
)


func Serialize(transaction *BBcTransaction, formatType uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	Put2byte(buf, formatType)
	dat, err := transaction.Pack()
	if err != nil {
		return nil, err
	}

	if formatType == 0 {
		if err := binary.Write(buf, binary.LittleEndian, dat); err != nil {
			return nil, err
		}
	} else if formatType == 0x0010 {
		compressed := ZlibCompress(&dat)
		if err := binary.Write(buf, binary.LittleEndian, compressed); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}


func Deserialize(dat []byte) (*BBcTransaction, error) {
	buf := bytes.NewBuffer(dat)

	formatType, err := Get2byte(buf)
	if err != nil {
		return nil, err
	}

	txdat, err := GetBytes(buf, len(dat)-2)
	if err != nil {
		return nil, err
	}

	if formatType == 0 {
		txobj := BBcTransaction{}
		txobj.Unpack(&txdat)
		return &txobj, nil
	} else if formatType == 0x0010 {
		decompressed, err := ZlibDecompress(txdat)
		if err != nil {
			return nil, err
		}
		txobj := BBcTransaction{}
		txobj.Unpack(&decompressed)
		return &txobj, nil
	}

	return nil, errors.New("formatType not supported")
}


func MakeTransaction(eventNum, relationNum int, witness bool, idLength int) *BBcTransaction {
	txobj := BBcTransaction{IdLength:idLength}
	for i:=0; i<eventNum; i++ {
		evt := BBcEvent{}
		txobj.AddEvent(&evt)
	}

	for i:=0; i<relationNum; i++ {
		rtn := BBcRelation{}
		txobj.AddRelation(&rtn)
	}

	if witness {
		wit := BBcWitness{}
		txobj.AddWitness(&wit)
	}

	return &txobj
}

func addInRelation(transaction *BBcTransaction, relationIdx int, assetGroupId, userId *[]byte) {
	ast := BBcAsset{}
	transaction.Relations[relationIdx].Add(assetGroupId, &ast)
	ast.Add(userId)
}

func AddRelationAssetFile(transaction *BBcTransaction, relationIdx int, assetGroupId, userId, assetFile *[]byte) {
	if transaction == nil {
		return
	}
	addInRelation(transaction, relationIdx, assetGroupId, userId)
	if assetFile != nil {
		transaction.Relations[relationIdx].Asset.AddFile(assetFile)
	}
}

func AddRelationAssetBodyString(transaction *BBcTransaction, relationIdx int, assetGroupId, userId *[]byte, body string) {
	if transaction == nil {
		return
	}
	addInRelation(transaction, relationIdx, assetGroupId, userId)
	if body != "" {
		transaction.Relations[relationIdx].Asset.AddBodyString(body)
	}
}

func AddRelationAssetBodyObject(transaction *BBcTransaction, relationIdx int, assetGroupId, userId *[]byte, body interface{}) {
	if transaction == nil {
		return
	}
	addInRelation(transaction, relationIdx, assetGroupId, userId)
	if body != nil {
		transaction.Relations[relationIdx].Asset.AddBodyObject(body)
	}
}

func AddRelationPointer(transaction *BBcTransaction, relationIdx int, refTransactionId, refAssetId *[]byte) {
	if transaction == nil {
		return
	}
	ptr := BBcPointer{}
	transaction.Relations[relationIdx].AddPointer(&ptr)
	ptr.Add(refTransactionId, refAssetId)
}

func AddPointerInRelation(relation *BBcRelation, refTransaction *BBcTransaction, refAssetId *[]byte) {
	ptr := BBcPointer{}
	relation.AddPointer(&ptr)
	ptr.Add(&refTransaction.TransactionId, refAssetId)
}

func AddReference(transaction *BBcTransaction, assetGroupId, userId *[]byte, refTransaction *BBcTransaction, eventIdx int) {
	if transaction == nil || refTransaction == nil {
		return
	}
	if refTransaction.TransactionId == nil {
		refTransaction.Digest()
	}
	ref := BBcReference{}
	transaction.AddReference(&ref)
	ref.Add(assetGroupId, refTransaction, eventIdx)
}

func addInEvent(transaction *BBcTransaction, eventIdx int, assetGroupId, userId *[]byte) {
	ast := BBcAsset{}
	transaction.Events[eventIdx].Add(assetGroupId, &ast)
	ast.Add(userId)
}


func AddEventAssetFile(transaction *BBcTransaction, eventIdx int, assetGroupId, userId *[]byte, assetFile *[]byte) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupId, userId)
	if assetFile != nil {
		transaction.Events[eventIdx].Asset.AddFile(assetFile)
	}
}

func AddEventAssetBodyString(transaction *BBcTransaction, eventIdx int, assetGroupId, userId *[]byte, body string) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupId, userId)
	if body != "" {
		transaction.Events[eventIdx].Asset.AddBodyString(body)
	}
}

func AddEventAssetBodyObject(transaction *BBcTransaction, eventIdx int, assetGroupId, userId *[]byte, body interface{}) {
	if transaction == nil {
		return
	}
	addInEvent(transaction, eventIdx, assetGroupId, userId)
	if body != "" {
		transaction.Events[eventIdx].Asset.AddBodyObject(body)
	}
}

func MakeRelationWithAsset(assetGroupId, userId *[]byte, assetBodyString string, assetBodyObject interface{}, assetFile *[]byte, idLength int) *BBcRelation {
	rtn := BBcRelation{IdLength:idLength}
	ast := BBcAsset{IdLength:idLength}
	ast.Add(userId)
	rtn.Add(assetGroupId, &ast)
	if assetFile != nil {
		ast.AddFile(assetFile)
	}
	if assetBodyString != "" {
		ast.AddBodyString(assetBodyString)
	} else if assetBodyObject != nil {
		ast.AddBodyObject(assetBodyObject)
	}
	return &rtn
}

func RecoverSignatureObject(dat *[]byte) *BBcSignature {
	sig := BBcSignature{}
	sig.Unpack(*dat)
	return &sig
}

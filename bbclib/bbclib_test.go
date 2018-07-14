package bbclib

import (
    "testing"
    "time"
    "encoding/hex"
    "gopkg.in/mgo.v2/bson"
    "fmt"
)

var (
    hexdat = "010009040000037472616e73616374696f6e5f62617365006c0300000368656164657200300000001076657273696f6e00010000001074696d657374616d7000c06a495b1069645f6c656e677468000800000000046576656e747300d7020000033000660100000561737365745f67726f75705f6964000800000000b34582b9d4640589047265666572656e63655f696e6469636573000500000000046d616e6461746f72795f617070726f766572730015000000053000080000000061cabbd2d0068eb200106f7074696f6e5f617070726f7665725f6e756d5f6e756d657261746f720000000000106f7074696f6e5f617070726f7665725f6e756d5f64656e6f6d696e61746f720000000000046f7074696f6e5f617070726f76657273000500000000036173736574009c0000000561737365745f69640008000000002ccab0b5c0477b7305757365725f696400080000000061cabbd2d0068eb2056e6f6e6365000800000000ea5899f3acc513321061737365745f66696c655f73697a6500000000000a61737365745f66696c655f646967657374001061737365745f626f64795f73697a6500080000000561737365745f626f647900080000000031323334353637380000033100660100000561737365745f67726f75705f6964000800000000b34582b9d4640589047265666572656e63655f696e6469636573000500000000046d616e6461746f72795f617070726f766572730015000000053000080000000061cabbd2d0068eb200106f7074696f6e5f617070726f7665725f6e756d5f6e756d657261746f720000000000106f7074696f6e5f617070726f7665725f6e756d5f64656e6f6d696e61746f720000000000046f7074696f6e5f617070726f76657273000500000000036173736574009c0000000561737365745f69640008000000007790ac9d30e5e7b605757365725f696400080000000061cabbd2d0068eb2056e6f6e6365000800000000854ca7f3991867b41061737365745f66696c655f73697a6500070000000561737365745f66696c655f6469676573740008000000007d1a54127b2225021061737365745f626f64795f73697a6500000000000a61737365745f626f647900000000047265666572656e6365730005000000000472656c6174696f6e73000500000000037769746e657373002600000004757365725f696473000500000000047369675f696e646963657300050000000000000363726f73735f726566006a00000005646f6d61696e5f696400200000000023e4473d7119dbc5cdef6c3acb8251ff4b21b6598d36c4a6ed0a780f1296ed74057472616e73616374696f6e5f6964002000000000545e5e45e2af451ff4827efcc72183992c15261620852b0c0229da29d3b9d6f500047369676e61747572657300050000000000"
)

func TestBBcTransactionSerialize(t *testing.T) {
    obj := BBcTransaction{}
    obj.Tx_base.Header.Version = 1
    obj.Tx_base.Header.Timestamp = (int)(time.Now().Unix())
    obj.Tx_base.Header.Id_length = 32
    obj.Format_type = FORMAT_BSON
    obj.Digest()
    fmt.Println(obj.Jsonify())
    fmt.Println("--------------------------------------")

    dat, err := obj.Serialize(false)
    if err != nil {
        t.Fatalf("failed to serialize transaction object (%v)", err)
    }
    t.Log("--------------------------------------")
    t.Logf("transaction_id: %x", obj.Transaction_id)
    t.Logf("serialize: %x", dat)


    obj1, err := BBcTransactionDeserialize(dat)
    if err != nil {
        t.Fatalf("failed to deserialize transaction object (%v)", err)
    }
    t.Logf("deserialized: %v", obj1)
    t.Log("=================================")

    dat2, err := hex.DecodeString(hexdat)
    var out map[string]interface{}
    bson.Unmarshal(dat2[2:], &out)
    t.Logf("direct unmarshal: %v", out)

    obj2, err := BBcTransactionDeserialize(dat2)
    if err != nil {
        t.Fatalf("failed to deserialize transaction object (%v)", err)
    }
    t.Logf("deserialized: %v", obj2)
    t.Logf("txid: %x", obj2.Transaction_id)

    dat3, err := obj2.Serialize(false)
    t.Logf("serialize for id: %x", dat3)

    t.Log(obj2.Stringer())
    if result, i := obj2.VerifyAll(); !result {
        t.Fatalf("Verify failed at %d", i)
    }
    t.Log("Vefiry succeeded")
}
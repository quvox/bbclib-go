package bbclib

import (
	"bytes"
	"compress/zlib"
)

func ZlibCompress(dat *[]byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(*dat)
	w.Close()
	return b.Bytes()
}


func ZlibDecompress(dat *[]byte) ([]byte, error) {
	var ret []byte
	bb := bytes.NewReader(*dat)
	r, err := zlib.NewReader(bb)
	r.Read(ret)
	return ret, err
}

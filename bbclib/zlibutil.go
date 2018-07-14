package bbclib

import (
	"bytes"
	"compress/zlib"
	"io"
)

func ZlibCompress(dat *[]byte) []byte {
	var dstbbuf bytes.Buffer
	zlibwriter := zlib.NewWriter(&dstbbuf)
	zlibwriter.Write(*dat)
	zlibwriter.Close()
	return dstbbuf.Bytes()
}


func ZlibDecompress(dat []byte) ([]byte, error) {
	var srcbuf bytes.Buffer
	var dstbuf bytes.Buffer
	srcbuf.Write(dat)
	zlibreader, err := zlib.NewReader(&srcbuf)
	if err != nil {
		return nil, err
	}
	io.Copy(&dstbuf, zlibreader)
	zlibreader.Close()

	return dstbuf.Bytes(), nil
}

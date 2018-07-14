package bbclib

import (
	"testing"
	"reflect"
)


func TestCompressDecompress(t *testing.T) {
	original := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	comp := ZlibCompress(&original)
	t.Logf("compressed: %x\n", comp)

	decomp, err := ZlibDecompress(comp)
	if err != nil {
		t.Fatalf("failed to decompress (%v)", err)
	}
	t.Logf("compressed: %x\n", decomp)
	if ! reflect.DeepEqual(original, decomp) {
		t.Fatal("failed to decompress (mismatch)")
	}
	t.Log("Succeeded")
}

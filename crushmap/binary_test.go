package crushmap

import (
	"io/ioutil"
	"testing"
)

func TestBinary(t *testing.T) {
	buf, err := ioutil.ReadFile("testdata/map.bin")
	if err != nil {
		t.Fatal(err)
	}

	m := NewMap()
	err = m.DecodeBinary(buf)
	if err != nil {
		t.Fatal(err)
	}
	_ = m
}

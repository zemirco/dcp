package dcp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEthernetUnmarshalBinary(t *testing.T) {
	b := []byte{
		0x01, 0x0e, 0xcf, 0x00, 0x00, 0x00, 0xa4, 0x4c,
		0xc8, 0xe5, 0x47, 0x21, 0x88, 0x92,
	}
	var e EthernetII
	if err := e.UnmarshalBinary(b); err != nil {
		t.Error(err)
	}
	if e.Destination.String() != "01:0e:cf:00:00:00" {
		t.Errorf("expected %s; got %s", "01:0e:cf:00:00:00", e.Destination)
	}
	if e.Source.String() != "a4:4c:c8:e5:47:21" {
		t.Errorf("expected %s; got %s", "a4:4c:c8:e5:47:21", e.Source)
	}
	if e.EtherType != 0x8892 {
		t.Errorf("expected %d; got %d", 0x8892, e.EtherType)
	}
}

func TestEthernetMarshalBinary(t *testing.T) {
	destination := []byte{0x01, 0x0e, 0xcf, 0x00, 0x00, 0x00}
	source := []byte{0xa4, 0x4c, 0xc8, 0xe5, 0x47, 0x21}

	e := NewEthernetII(destination, source)

	b, err := e.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	expected := []byte{
		0x01, 0x0e, 0xcf, 0x00, 0x00, 0x00, 0xa4, 0x4c,
		0xc8, 0xe5, 0x47, 0x21, 0x88, 0x92,
	}
	if diff := cmp.Diff(b, expected); diff != "" {
		t.Error(diff)
	}
}

func TestEthernetLen(t *testing.T) {
	e := &EthernetII{}
	if e.Len() != 14 {
		t.Errorf("expected %d; got %d", 14, e.Len())
	}
}

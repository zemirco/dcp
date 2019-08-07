package block

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

func TestAllAllUnmarshalBinary(t *testing.T) {
	b := []byte{0xff, 0xff, 0x00, 0x00}
	var a All
	if err := a.UnmarshalBinary(b); err != nil {
		t.Error(err)
	}
	if a.Option != option.All {
		t.Errorf("expected %d; got %d", option.All, a.Option)
	}
	if a.Suboption != suboption.All {
		t.Errorf("expected %d; got %d", suboption.All, a.Suboption)
	}
}

func TestAllAllMarshalBinary(t *testing.T) {
	a := NewAll()
	b, err := a.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	expected := []byte{0xff, 0xff, 0x00, 0x00}
	if diff := cmp.Diff(b, expected); diff != "" {
		t.Error(diff)
	}
}

func TestAllAllLen(t *testing.T) {
	a := NewAll()
	if a.Len() != 4 {
		t.Errorf("expected %d; got %d", 4, a.Len())
	}
}

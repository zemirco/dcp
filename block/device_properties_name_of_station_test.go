package block

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

func TestDevicePropertiesNameOfStationUnmarshalBinary(t *testing.T) {
	b := []byte{
		0x02, 0x02, 0x00, 0x07, 0x00, 0x00, 0x7a, 0x65,
		0x69, 0x73, 0x73,
	}

	nos := &NameOfStation{}
	nos.HasInfo = true

	if err := nos.UnmarshalBinary(b); err != nil {
		t.Error(err)
	}
	if nos.Option != option.Properties {
		t.Errorf("expected %d; got %d", option.Properties, nos.Option)
	}
	if nos.Suboption != suboption.NameOfStation {
		t.Errorf("expected %d; got %d", suboption.NameOfStation, nos.Suboption)
	}
	if nos.NameOfStation != "zeiss" {
		t.Errorf("expected %s; got %s", "zeiss", nos.NameOfStation)
	}
}

func TestDevicePropertiesNameOfStationMarshalBinary(t *testing.T) {
	nos := NewNameOfStationWithInfo(0, "zeiss")
	b, err := nos.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	expected := []byte{
		0x02, 0x02, 0x00, 0x07, 0x00, 0x00, 0x7a, 0x65,
		0x69, 0x73, 0x73,
	}
	if diff := cmp.Diff(b, expected); diff != "" {
		t.Error(diff)
	}
}

func TestDevicePropertiesNameOfStationLen(t *testing.T) {
	nos := NewNameOfStationWithInfo(0, "zeiss")
	if nos.Len() != 11 {
		t.Errorf("expected %d; got %d", 11, nos.Len())
	}
}

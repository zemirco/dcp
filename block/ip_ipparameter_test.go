package block

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIPIPParameterUnmarshalBinary(t *testing.T) {
	b := []byte{
		0x01, 0x02, 0x00, 0x0e, 0x00, 0x01, 0xac, 0x13,
		0x68, 0x05, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	}

	i := &IPParameter{}
	i.HasInfo = true

	if err := i.UnmarshalBinary(b); err != nil {
		t.Error(err)
	}
	if i.Info != 1 {
		t.Errorf("expected %d; got %d", 1, i.Info)
	}
	if i.IPAddress.String() != "172.19.104.5" {
		t.Errorf("expected %s; got %s", "172.19.104.5", i.IPAddress)
	}
	if i.Subnetmask.String() != "255.255.0.0" {
		t.Errorf("expected %s; got %s", "255.255.0.0", i.Subnetmask)
	}
	if i.StandardGateway.String() != "0.0.0.0" {
		t.Errorf("expected %s; got %s", "0.0.0.0", i.StandardGateway)
	}
}

func TestIPIPParameterMarshalBinary(t *testing.T) {
	ip := []byte{0xac, 0x13, 0x68, 0x05}
	subnet := []byte{0xff, 0xff, 0x00, 0x00}
	gateway := []byte{0x00, 0x00, 0x00, 0x00}

	i := NewIPParameterWithInfo(ip, subnet, gateway, 1)

	b, err := i.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	expected := []byte{
		0x01, 0x02, 0x00, 0x0e, 0x00, 0x01, 0xac, 0x13,
		0x68, 0x05, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	}
	if diff := cmp.Diff(b, expected); diff != "" {
		t.Error(diff)
	}
}

package block

import (
	"net"

	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

// IPParameter is an ip parameter block.
type IPParameter struct {
	header
	IPAddress       net.IP
	Subnetmask      net.IP
	StandardGateway net.IP
}

var _ Block = &IPParameter{}

// NewIPParameterWithInfo returns a new block.
func NewIPParameterWithInfo(ip, subnet, gateway net.IP, info uint16) *IPParameter {
	return &IPParameter{
		header: header{
			Option:       option.IP,
			Suboption:    suboption.IPParameter,
			Length:       14,
			HasInfo:      true,
			Info:         info,
			HasQualifier: false,
		},
		IPAddress:       ip,
		Subnetmask:      subnet,
		StandardGateway: gateway,
	}
}

// NewIPParameterQualifier returns a new block.
func NewIPParameterQualifier() *IPParameter {
	return &IPParameter{
		header: header{
			Option:       option.IP,
			Suboption:    suboption.IPParameter,
			Length:       14,
			HasInfo:      false,
			HasQualifier: true,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (i *IPParameter) UnmarshalBinary(b []byte) error {
	if err := i.header.unmarshalBinary(b); err != nil {
		return err
	}

	offset := i.header.len()

	i.IPAddress = net.IP(b[offset : offset+4])
	offset += 4

	i.Subnetmask = net.IP(b[offset : offset+4])
	offset += 4

	i.StandardGateway = net.IP(b[offset : offset+4])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (i *IPParameter) MarshalBinary() ([]byte, error) {
	b := make([]byte, i.Len())

	bh, err := i.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += i.header.len()

	copy(b[offset:offset+4], i.IPAddress)
	offset += 4

	copy(b[offset:offset+4], i.Subnetmask)
	offset += 4

	copy(b[offset:offset+4], i.StandardGateway)

	return b, nil
}

// Len returns length for ip parameter block.
func (i *IPParameter) Len() int {
	return i.header.len() + 4 + 4 + 4
}

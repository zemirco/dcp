package block

// NameOfStation is a name of station block.
type NameOfStation struct {
	header
	NameOfStation string
}

var _ Block = &NameOfStation{}

// NewNameOfStation returns a new block.
func NewNameOfStation(hasInfo bool) *NameOfStation {
	return &NameOfStation{
		header: header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (n *NameOfStation) UnmarshalBinary(b []byte) error {
	if err := n.header.unmarshalBinary(b); err != nil {
		return err
	}

	i := n.header.len()
	n.NameOfStation = string(b[i : int(i)+int(n.header.Length)-2])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (n *NameOfStation) MarshalBinary() ([]byte, error) {
	b := make([]byte, n.Len())

	bh, err := n.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += n.header.len()

	copy(b[offset:], n.NameOfStation)

	return b, nil
}

// Len returns length for name of station block.
func (n *NameOfStation) Len() int {
	return n.header.len() + len(n.NameOfStation)
}

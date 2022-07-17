package utils

// ********************** STRING *****************
type MarshalableString string

func (m MarshalableString) MarshalBinary() ([]byte, error) {
	b := make([]byte, len(m)+1)
	b[0] = byte(len(m))
	copy(b[1:], []byte(m))
	return b, nil
}

func (m MarshalableString) UnmarshalBinary(data []byte) error {
	m = MarshalableString(data)
	return nil
}

// ********************* INT *********************
type MarshalableInt32 int

func (m MarshalableInt32) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)
	b[0] = byte(m >> 24)
	b[1] = byte(m >> 16)
	b[2] = byte(m >> 8)
	b[3] = byte(m)
	return b, nil
}

func (m MarshalableInt32) UnmarshalBinary(data []byte) error {
	m = MarshalableInt32(int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3]))
	return nil
}

type MarshalableInt16 int

func (m MarshalableInt16) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	b[0] = byte(m >> 8)
	b[1] = byte(m)
	return b, nil
}

func (m MarshalableInt16) UnmarshalBinary(data []byte) error {
	m = MarshalableInt16(int(data[0])<<8 | int(data[1]))
	return nil
}

type MarshalableInt8 int

func (m MarshalableInt8) MarshalBinary() ([]byte, error) {
	b := make([]byte, 1)
	b[0] = byte(m)
	return b, nil
}

func (m MarshalableInt8) UnmarshalBinary(data []byte) error {
	m = MarshalableInt8(data[0])
	return nil
}

// ******************** BOOL ********************
type MarshalableBool bool

func (m MarshalableBool) MarshalBinary() ([]byte, error) {
	if m {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

func (m MarshalableBool) UnmarshalBinary(data []byte) error {
	m = MarshalableBool(data[0] == 1)
	return nil
}

// ******************** BYTE ARRAY ********************
type MarshalableByteArray []byte

func (m MarshalableByteArray) MarshalBinary() ([]byte, error) {
	b := make([]byte, len(m)+1)
	b[0] = byte(len(m))
	copy(b[1:], m)
	return b, nil
}

func (m MarshalableByteArray) UnmarshalBinary(data []byte) error {
	m = MarshalableByteArray(data)
	return nil
}

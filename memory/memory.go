package memory

// BaseMemory is a wrapper around a normal []byte.
// Its main utility is to be both printable and clearable.
type BaseMemory []byte

// NewBaseMemory creates a BaseMemory with size size.
func NewBaseMemory(size int) (m BaseMemory) {
	return BaseMemory(make([]byte, size))
}

// String implements Stringer.
// TODO: complete memory dump
func (m BaseMemory) String() (repr string) {
	return ""
}

// Write writes a byte of data to memory at location address.
func (m BaseMemory) Write(address uint16, data byte) {
	m[address] = data
}

// Read reads a byte of data from memory at address.
func (m BaseMemory) Read(address uint16) (b byte) {
	return m[address]
}

// Clear sets all data in memory to 0x00.
func (m BaseMemory) Clear() {
	for i := range m {
		m[i] = 0x00
	}
}

// MirroredMemory is a memory that is mirrored
// at a certain frequency.
type MirroredMemory struct {
	freq int
	mem  []byte
}

// NewMirroredMemory creates a MirroredMemory with freq freq.
func NewMirroredMemory(freq int) (mm MirroredMemory) {
	return MirroredMemory{
		freq: freq,
		mem:  make([]byte, freq),
	}
}

// String implements Stringer.
// TODO: complete memory dump
func (mm MirroredMemory) String() (repr string) {
	return ""
}

// Write writes a byte of data to mirrored memory at location address.
func (mm MirroredMemory) Write(address uint16, data byte) {
	mm.mem[address%uint16(mm.freq)] = data
}

// Read reads a byte of data from mirrored memory at address.
func (mm MirroredMemory) Read(address uint16) (b byte) {
	return mm.mem[address%uint16(mm.freq)]
}

// Clear sets all mirrored memory to 0x00.
func (mm MirroredMemory) Clear() {
	for i := range mm.mem {
		mm.mem[i] = 0x00
	}
}

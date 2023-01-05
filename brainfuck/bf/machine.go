package bf

type Machine interface {
	PtrIncr()
	PtrDecr()
	Incr()
	Decr()
	Peek() byte
	Put(byte)
	Dump(uint, uint) []byte
}

type machine struct {
	size  uint
	cells []byte
	ptr   uint
}

func NewMachine(size uint) Machine {
	return &machine{
		size:  size,
		cells: make([]byte, size),
		ptr:   0,
	}
}

func (m *machine) PtrIncr() {
	m.ptr++
	if m.ptr >= m.size {
		m.ptr = 0
	}
}

func (m *machine) PtrDecr() {
	if m.ptr == 0 {
		m.ptr = m.size - 1
	} else {
		m.ptr--
	}
}

func (m *machine) Incr() {
	m.cells[m.ptr] = m.cells[m.ptr] + 1
}

func (m *machine) Decr() {
	m.cells[m.ptr] = m.cells[m.ptr] - 1
}

func (m *machine) Peek() byte {
	return m.cells[m.ptr]
}

func (m *machine) Put(v byte) {
	m.cells[m.ptr] = v
}

func (m *machine) Dump(start, end uint) []byte {
	return m.cells[start:end]
}

package memory

import (
	s "pkg/shared"
	b "pkg/bits"
)

type Mem struct {
	buf         []uint8
	WriteLocked bool
}

func (m *Mem) Init(size uint16) {
	m.buf = make([]uint8, size)
	m.zero()
}

func (m *Mem) GetWordAt(addr uint16) (uint16, error) {
	e := m.boundsCheck(addr)
	if e != nil {
		return 0, e
	}

	hiByte := m.buf[addr]
	loByte := m.buf[addr + 1]
	wordValue := b.BytesToWord(hiByte,loByte)
	return wordValue, nil
}

func (m *Mem) SetWordAt(addr uint16, v uint16) error {
	e := m.boundsCheck(addr)
	if e != nil {
		return e
	}

	hiByte, loByte := b.WordToBytes(v)

	if !m.WriteLocked {
		m.buf[addr] = hiByte 
		m.buf[addr+1] = loByte
	}
	return nil
}

func (m *Mem) GetByteAt(addr uint16) (uint8, error) {
	e := m.boundsCheck(addr)
	if e != nil {
		return 0, e
	}

	return m.buf[addr], e
}

func (m *Mem) SetByteAt(addr uint16, v uint8) error {
	e := m.boundsCheck(addr)
	if e != nil {
		return e
	}

	if !m.WriteLocked {
		m.buf[addr] =v
	}
	return nil
}

func (m *Mem) WriteLock() {
	m.WriteLocked = true
}

func (m *Mem) WriteUnlock() {
	m.WriteLocked = false
}

func (m *Mem) IsWriteLocked() bool {
	return m.WriteLocked
}

func (m *Mem) zero() {
	for i, _ := range m.buf {
		m.buf[i] = uint8(0)
	}
}

func (m *Mem) boundsCheck(addr uint16) error {
	if uint16(addr) > uint16(len(m.buf)) {
		return &s.MemoryError{"Out Of Range", addr}
	}
	return nil
}

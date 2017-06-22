package stores

import (
	"bytes"
	"errors"
	"fmt"
)

// A MemoryStore stores all records data into the memory
type MemoryStore struct {
	data        [][]byte
	workingData [][]byte
}

// Store add the provided data to the in-memory array
func (m *MemoryStore) Store(d []byte) error {
	m.data = append(m.data, d)
	return nil
}

// Retrieve pops the latest data from the in-memory array
func (m *MemoryStore) Retrieve() ([]byte, error) {
	if len(m.data) == 0 {
		return nil, nil
	}

	d, a := m.data[len(m.data)-1], m.data[:len(m.data)-1]
	m.data = a
	m.workingData = append(m.workingData, d)
	return d, nil
}

// Finish marks a task as finished
func (m *MemoryStore) Finish(d []byte) error {
	if len(m.workingData) == 0 {
		return errors.New("no working data found")
	}

	var nw [][]byte
	for _, v := range m.workingData {
		if bytes.Compare(v, d) == 0 {
			continue
		}
		nw = append(nw, v)
	}
	if len(nw) == len(m.workingData) {
		return fmt.Errorf("unknown working record %q", d)
	}
	m.workingData = nw

	return nil
}

// Length returns the number of elements in the in-memory array
func (m *MemoryStore) Length(q string) (int, error) {
	switch q {
	case "waiting":
		return len(m.data), nil
	case "working":
		return len(m.workingData), nil
	default:
		return 0, fmt.Errorf("unknown queue %s", q)
	}

}

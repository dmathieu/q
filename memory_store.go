package q

// A MemoryStore stores all records data into the memory
type MemoryStore struct {
	data [][]byte
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
	return d, nil
}

// Length returns the number of elements in the in-memory array
func (m *MemoryStore) Length() (int, error) {
	return len(m.data), nil
}

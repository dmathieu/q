package queue

import "sync"

type fakeStore struct {
	data [][]byte
	m    sync.Mutex
}

func (f *fakeStore) Store(d []byte) error {
	f.m.Lock()
	defer f.m.Unlock()

	f.data = append(f.data, d)
	return nil
}

func (f *fakeStore) Retrieve() ([]byte, error) {
	f.m.Lock()
	defer f.m.Unlock()

	if len(f.data) == 0 {
		return nil, nil
	}

	d, a := f.data[len(f.data)-1], f.data[:len(f.data)-1]
	f.data = a
	return d, nil
}

func (f *fakeStore) Finish(d []byte) error {
	return nil
}

func (f *fakeStore) Length(q string) (int, error) {
	return len(f.data), nil
}

func (f *fakeStore) HouseKeeping() error {
	return nil
}

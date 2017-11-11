package memdb

type DB struct {
}

func (db *DB) Put(key, val []byte) error {
	return nil
}

func (db *DB) Delete(key []byte) error {
	return nil
}

func (db *DB) Contains(key []byte) bool {
	return false
}

func (db *DB) Get(key []byte) (val []byte, err error) {
	return nil, nil
}

func (db *DB) Find(key []byte) (rkey, val []byte, err error) {
	return nil, nil, nil
}

A [levelDB](http:code.google.com/p/leveldb) style [LSM](https://en.wikipedia.org/wiki/Log-structured_merge-tree) key-value store in Go programming language.

#### How to install

go get https://github.com/snowmagic1/lsmdb

#### Features
  * Keys and values are arbitrary byte arrays.
  * Data is stored sorted by key.
  * Callers can provide a custom comparison function to override the sort order.
  * Multiple changes are merged into one batch
  * Users can create a transient snapshot to get a consistent view of data.
  
#### How to use

Create or open a database
```go
db, err := lsmdb.OpenFile("path to db file")
if err != nil {
  ...
}
defer db.close
```

Lookup by a key
```go
key := []byte("key1")
rval, err := lsmdb.Get(key)
if err != nil {
		...
}
```

Modify or delete a key
```go
key := []byte("key1")
val := []byte("val1")
err = lsmdb.Put(key, val)
if err != nil {
 ...
}

err = db.Delete(key)
if err != nil {
  ...
}
```

Read a range of keys
```go
iter := lsmdb.NewIterator(nil, nil)
for iter.Next() {
	key := iter.Key()
	value := iter.Value()
	...
}
iter.Release()
```

#### TODO
* Batch write
* Transactions that allow fast data ingress
* Bloom filter on table blocks

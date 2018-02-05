A [levelDB](http:code.google.com/p/leveldb) style [LSM](https://en.wikipedia.org/wiki/Log-structured_merge-tree) key-value store in Go programming language.

#### How to install

go get https://github.com/snowmagic1/lsmdb

#### Features
  * [Data model] Key-value store, the basic operations are Put(key, val), Get(key) and Delete(key)
  * [Storage] LSM storage engine, skiplist is used for in-memory data, multi-level compaction
  * [Snapshot isolation] Users can create a transient snapshot to get a consistent view of data
  * [Concurrency] Different thread may write or read without any external synchoronization
  * [High Performance]
  
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

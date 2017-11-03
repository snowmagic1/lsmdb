package db

import (
	"log"
	"os"

	"github.com/snowmagic1/lsmdb/wallog"
)

type DB struct {
	dataFile *os.File
	logFile  *os.File
	logRW    *wallog.LogRW
}

func Open(dbname string) (db *DB) {

	newDb := &DB{}
	var err error
	newDb.dataFile, err = os.OpenFile(dbname+".data", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("failed to open data file ", err)
	}

	newDb.logFile, err = os.OpenFile(dbname+".log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("failed to open data file ", err)
	}

	newDb.logRW = wallog.NewLogRW(newDb.logFile)

	return newDb
}

func (db *DB) Close() {
	db.dataFile.Close()
	db.logFile.Close()
}

func (db *DB) Put(key, val string) bool {
	log.Printf("[Put] key [%v] val [%v]\n", key, val)

	db.logRW.Write(key, val)
	db.logFile.Sync()

	return true
}

func (db *DB) Get(key string) (val string, ok bool) {
	if val, ok := db.Get(key); !ok {
		log.Println("failed to get, err ", ok)
	} else {
		log.Printf("[Get] key [%v] val [%v]\n", key, val)
	}

	return val, ok
}

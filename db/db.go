package db

import (
	"log"
	"os"
)

type DB struct {
	dataFile  *os.File
	logFile   *os.File
	logWriter *LogWriter
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

	newDb.logWriter = NewLogWriter(newDb.logFile)

	return newDb
}

func (db *DB) Close() {
	db.dataFile.Close()
	db.logFile.Close()
}

func (db *DB) Put(key, val string) bool {
	log.Printf("[Put] key [%v] val [%v]\n", key, val)

	db.logWriter.AddRecord([]byte{1, 2, 3, 4})

	return true
}

func (db *DB) Get(key string) (val string, ok bool) {

	/*
		if val, ok := db.Get(key); !ok {
			log.Println("failed to get, err ", ok)
		} else {
			log.Printf("[Get] key [%v] val [%v]\n", key, val)
		}
	*/

	return "", ok
}

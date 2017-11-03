package main

import (
	"log"
	"os"

	"github.com/snowmagic1/lsmdb/db"
)

type myStruct struct {
	ID   string
	Data string
}

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", ":", err)
	}

	defer file.Close()
	log.SetOutput(file)

	db := db.Open("db1")
	db.Put("key1", "valval1")
	val, _ := db.Get("key1")
	log.Println("val [", val, "]")

	db.Close()
}

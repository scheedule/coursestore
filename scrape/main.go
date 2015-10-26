package main

import (
	"github.com/scheedule/coursestore/db"
	"github.com/scheedule/coursestore/types"
)

func main() {
	err := PopulateDB(
		"http://courses.illinois.edu/cisapp/explorer/schedule/2016/spring.xml",
		"mongo",
		"27017",
		"test",
		"classes",
	)

	if err != nil {
		panic(err)
	}
}

func PopulateDB(term_url, ip, port, db_name, collection_name string) error {
	mydb := db.NewDB(ip, port, db_name, collection_name)

	err := mydb.Init()
	if err != nil {
		return err
	}

	mydb.Purge()

	term, err := GetXML(term_url)

	course_chan := make(chan types.Class)

	go DigestAll(term, course_chan)

	for class := range course_chan {
		err = mydb.Put(class)
		if err != nil {
			return err
		}
	}

	return nil
}

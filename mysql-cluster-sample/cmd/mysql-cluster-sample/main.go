package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlUser = "dbwebapp"
	mysqlPass = "password"
	mysqlAddr = "localhost:6446" // mysql-router host
	mysqlDb   = "test_db"
)

type Employee struct {
	ID   int
	Name string
	City string
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
		mysqlUser, mysqlPass, mysqlAddr, mysqlDb,
	))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO `employee` (name, city) VALUES ('John Smith', 'New York')")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer insert.Close()

	results, err := db.Query("SELECT * FROM `employee`")
	if err != nil {
		log.Fatal(err.Error())
	}
	for results.Next() {
		var empl Employee
		err = results.Scan(&empl.ID, &empl.Name, &empl.City)
		if err != nil {
			log.Println(err.Error())
		}
		log.Printf("# %d - %s from %s", empl.ID, empl.Name, empl.City)
	}
}

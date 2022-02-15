package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
)

//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "admin"
//	password = "admin"
//	dbname   = "postgres"
//)

const (
	host     = "192.168.120.54"
	port     = 5432
	user     = "mol_api"
	password = "mol_api"
	dbname   = "mol_api_db"
)

func main() {
	professions := make(map[string]string)
	specs := make(map[string]string)
	content, err := excelize.OpenFile("C:\\Users\\SalmanSharifov\\Downloads\\Məşğulluq təsnifatı_Qruplaşma.xlsx")
	checkErr(err)
	fileData, err := content.GetRows("Positions")
	checkErr(err)
	index := 0
	s := ""
	for _, data := range fileData[1:] {
		if (professions[data[4]+s] != "") && data[6] != professions[data[4]+s] {
			index++
			if index > 0 {
				s = strconv.Itoa(index)
			}
			professions[data[4]+s] = data[6]
		} else {
			if professions[data[4]] == "" {
				professions[data[4]] = data[6]
				index = 0
				s = ""
			}
		}
		specs[data[9]] = data[8]
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	connection, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	//var id int
	for key, value := range professions {
		query := "UPDATE professions SET name = $1 WHERE classification_code = $2"
		res, er := connection.Exec(query, strings.Trim(value, ""), strings.Trim(key, ""))
		count, er := res.RowsAffected()
		fmt.Println("updated: ", count == 1)
		if count == 0 {
			id, _ := strconv.Atoi(key)
			insertResult, er := connection.Exec("INSERT INTO professions(id, created_at, creator_user, last_updated_at, "+
				" updater_user, classification_code, name, status) "+
				" VALUES($1,$2,$3,$4,$5,$6,$7, $8)", id, time.Now(), "admin", time.Now(), "admin", strings.Trim(key, ""), strings.Trim(value, ""), "ACTIVE")
			checkErr(er)
			c, errr := insertResult.RowsAffected()
			fmt.Println("inserted id", id, ", is inserted", c == 1)
			checkErr(errr)
		}
		checkErr(er)
	}

	for key, value := range specs {
		id, _ := strconv.Atoi(key)
		if value == "" {
			res, _ := connection.Exec("UPDATE specifications SET status = $1 WHERE id = $2", "INACTIVE", id)
			count, _ := res.RowsAffected()
			fmt.Println("status updated: ", count == 1)
		} else {
			res, _ := connection.Exec("UPDATE specifications SET name = $1 WHERE id = $2", strings.Trim(value, ""), id)
			count, _ := res.RowsAffected()
			fmt.Println("status updated: ", count == 1)
		}
	}
	defer connection.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

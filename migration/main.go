package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
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
	user     = "****"
	password = "****"
	dbname   = "****"
)

func main() {
	professions := make(map[string][]string)
	var specs []string
	content, err := excelize.OpenFile("C:\\Users\\SalmanSharifov\\Downloads\\Məşğulluq təsnifatı_Qruplaşma.xlsx")
	checkErr(err)
	fileData, err := content.GetRows("Positions")
	checkErr(err)
	for _, data := range fileData[1:] {
		for _, innerData := range fileData {
			if innerData[6] == data[6] && !contains(specs, innerData[8]) {
				specs = append(specs, innerData[8])
			}
		}
		professions[data[6]] = specs
		specs = nil
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	connection, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	for key, value := range professions {
		queryResult, queryError := connection.Exec("INSERT INTO professions(created_at, creator_user, last_updated_at, updater_user, classification_code, name, status) "+
			"VALUES ($1,$2,$3,$4,$5,$6,$7)", time.Now(), "king", time.Now(), "king", "", key, "ACTIVE")
		checkErr(queryError)
		profId, idError := queryResult.LastInsertId()
		checkErr(idError)
		for _, v := range value {
			specQueryResult, specQueryErr := connection.Exec("INSERT INTO specifications(created_at, creator_user, last_updated_at, updater_user, name, status) "+
				" VALUES ($1,$2,$3,$4,$5,$6)", time.Now(), "king", time.Now(), "king", v, "ACTIVE")
			checkErr(specQueryErr)
			specId, specErr := specQueryResult.LastInsertId()
			checkErr(specErr)
			_, err := connection.Exec("INSERT INTO profession_specifications(profession_id, specification_id) VALUES ($1, $2)", profId, specId)
			checkErr(err)
		}
	}

	for k, v := range professions {
		fmt.Println(k)
		for i := 0; i < len(v); i++ {
			fmt.Print(v[i], ", ")
		}
		fmt.Println()
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

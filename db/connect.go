package db

import (
	"fmt" 
	"log"
	"database/sql" 
	//_ "github.com/go-sql-driver/mysql" 
	_ "github.com/lib/pq"
)

func connect() *sql.DB {
	/*** mysql *****
	db, err := sql.Open("mysql","root:@tcp(127.0.0.1:3306)/rento-scanner-db")
	if err != nil {  fmt.Println("Failed to connect to the dabtabase"); fmt.Println(err.Error()) 
		return nil
	}
	*/
	//**** YugaByteDB *******
	const (
      host     = "192.168.39.29"
      port     = 30743
      user     = "rsadmin"
      password = "secr3t"
      dbname   = "rento_scanner_db"
    )
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
    if err != nil { fmt.Println("Failed to connect to the dabtabase"); fmt.Println(err.Error()) 
        fmt.Println("rsgateway connecting to the dabtabase ... Failed.");
        log.Fatal(err)
    }else{
	    fmt.Println("rsgateway connecting to the dabtabase ... OK.");
	}
	return db
} 
//--------------------------------




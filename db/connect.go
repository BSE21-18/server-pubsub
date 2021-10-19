package db

import (
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "github.com/joho/godotenv"
  "os"
)


func Connect() (*sql.DB, error) {
	err := godotenv.Load(".env")
      if err != nil {
        log.Fatalf("Error loading .env file")
      }
	
	dbUser := os.Getenv("MYSQL_DB_USER")
	dbPswd := os.Getenv("MYSQL_DB_PSWD")
	dbHost := os.Getenv("MYSQL_DB_HOST")
	dbPort := os.Getenv("MYSQL_DB_PORT")
	dbName := os.Getenv("MYSQL_DB_NAME")
	
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Printf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPswd, dbHost, dbPort, dbName)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
      if err != nil {
        return nil, err
      }
    //======= Migrate the schema ======
    db.AutoMigrate(&Customer{})
    db.AutoMigrate(&Subscription{})
    db.AutoMigrate(&ProcessingResult{})
    
	return db, nil
}





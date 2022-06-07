package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Res struct {
	Data any  `json:"data"`
	Err  bool `json:"error"`
}

// ============== logger ================
type logger struct {
	fileName string
	*log.Logger
}

var once sync.Once
var Logger = getErrLogger()

func getErrLogger() *logger {
	var logger *logger
	once.Do(func() {
		logger = setLogger("./log/err.log")
	})
	return logger
}

func setLogger(filePath string) *logger {
	file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	log.SetOutput(file)
	return &logger{
		fileName: file.Name(),
		Logger:   log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// ============== logger ================

// ============== database connect ================

var DB *sql.DB = connDB()

func connDB() *sql.DB {
	godotenv.Load("./.env")
	dbinfo := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Println(err)
	}

	return db
}

// ============== database connect ================

package util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-session/session/v3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Res struct {
	Data any  `json:"data"`
	Err  bool `json:"error"`
}

type ArgonConfig struct {
	Time   uint32
	Memory uint32
	Thread uint8
	KeyLen uint32
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

// ============== login check ================

func LoginCheck(w http.ResponseWriter, r *http.Request) interface{} {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return nil
	}

	data, ok := store.Get("userId")
	if !ok {
		return nil
	}

	return data
}

// ============== login check ================

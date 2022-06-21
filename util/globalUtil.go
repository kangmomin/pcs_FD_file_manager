package util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-redis/redis/v9"
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
var AdminLogger *logger

func getErrLogger() *logger {
	var logger *logger
	once.Do(func() {
		logger = setLogger("./log/err.log")
		AdminLogger = setLogger("./log/admin.log")
	})
	return logger
}

func setLogger(filePath string) *logger {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println(err)
		return nil
	}

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

// ============== redis connect ================

var Rdb = redisConn()

func redisConn() *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		log.Println(err)
	}
	return c
}

// ============== redis connect ================

// ============== login check ================

func LoginCheck(r *http.Request) interface{} {
	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Println(nil)
		}
		return nil
	}

	data, err := Rdb.Get(context.Background(), sessionID.Value).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	return data
}

func AdminCheck(r *http.Request) (interface{}, bool) {
	var data interface{}
	if data = LoginCheck(r); data == nil {
		return nil, false
	}

	err := DB.QueryRow(`SELECT admin_tier FROM admin WHERE user_id=$1 AND accept=true`, data).Err()
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return data, true
}

// ============== login check ================

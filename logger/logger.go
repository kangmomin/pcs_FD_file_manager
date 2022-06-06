package logger

import (
	"log"
	"os"
)

func LogSetting() {
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile) //로그 출력 위치를 파일로 변경
}

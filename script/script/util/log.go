package util

import (
	"fmt"
	"log"
	"os"
	"time"
)

const BaseLogDirectory = "/var/log/script/"
const LogDirectoryInfo = "info"
const LogDirectoryError = "error"

func Info(msg string) {

	fileName := time.Now().Format("20060102") + ".log"
	dir := BaseLogDirectory + LogDirectoryInfo + "/"
	CreateLogFile(dir, fileName).Printf("[INFO] %s\n", msg)
}

func InfoByName(fileName, msg string) {
	if "" == fileName {
		fileName = time.Now().Format("20060102") + ".log"
	}

	dir := BaseLogDirectory + LogDirectoryInfo + "/"
	CreateLogFile(dir, fileName).Printf("[INFO] %s\n", msg)
}

func Error(msg string) {
	fileName := time.Now().Format("20060102") + ".log"

	dir := BaseLogDirectory + LogDirectoryError + "/"

	CreateLogFile(dir, fileName).Printf("[ERROR] %s\n", msg)
}

func ErrorByName(fileName, msg string) {
	if "" == fileName {
		fileName = time.Now().Format("20060102") + ".log"
	}

	dir := BaseLogDirectory + LogDirectoryError + "/"

	CreateLogFile(dir, fileName).Printf("[ERROR] %s\n", msg)
}

func CreateLogFile(dir, fileName string) *log.Logger {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Cannot create a file when that file already exists %v \n ", err)
	}

	file, err := os.OpenFile(dir+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("CreateLogFile OpenFile %v \n ", err)
		return nil
	}

	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

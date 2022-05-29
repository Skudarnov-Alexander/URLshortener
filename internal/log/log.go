package logger

import (
	"fmt"
	"log"
	"os"
)

/*
func InitLogger() (l *log.Logger, err error) {

	nf, err := os.Create("logs")

	if err != nil {
		fmt.Println("loger error")
		return
	}

	l = log.New(nf, "LOG", 19)

	l.Printf("hello")

	return

}

*/

var Logger *log.Logger

func InitLogger() {
	
	nf, err := os.Create("./internal/log/logs")
	//os.Chdir("./intenal/log")
	nf.Chdir()
	Logger = log.New(nf, "LOG", 19)

	if err != nil {
		fmt.Println("loger error")
		return
	}
}


/*

func InitLogger() (logger *log.Logger, err error) {

	nf, err := os.Create("logs")

	if err != nil {
		fmt.Println("loger error")
		return
	}
		
	Logger = log.New(nf, "LOG", 19)

	Logger.Printf("hello")

	return

}

*/
package main

import (

	"github.com/sirupsen/logrus"
	"letaipays/pg"
	"barcode"
	"letaipays/tgbot"
	"os"
	"os/signal"
	"syscall"
	"time"
)



var token = "1607281585:AAH95th0NfsMxbCFtaF-fEY7g0s8CY-30W0"

func main() {


	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("starting")

	var st time.Time

	defer func() {
		logrus.Info("shutdown time - %s", time.Now().Sub(st))
	}()

	db, err := pg.NewStorage("host=192.168.143.179 user=letaipays password=Sk18sxsFV1#B712XC dbname=test sslmode=disable")

	if err != nil{
		logrus.Panic(err)
	}

	scanner := barcode.NewScanner()


	_, err= tgbot.NewBot(token, db, scanner)

	if err != nil{
		logrus.Panic(err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)


	logrus.Infof("captured %v signal^ stopping", <-signals)


}

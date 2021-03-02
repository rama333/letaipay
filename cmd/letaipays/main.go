package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/bieber/barcode.v0"
	"letaipays/pg"
	"letaipays/tgbot"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var token = "1616574093:AAGmVjKIQ5CYAWrU7bBD6uwgOwS7d_kAJq0"

func main() {

	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("starting")

	var st time.Time

	defer func() {
		logrus.Info("shutdown time - %s", time.Now().Sub(st))
	}()

	db, err := pg.NewStorage("host=192.168.143.179 user=letaipays password=Sk18sxsFV1#B712XC dbname=test sslmode=disable")

	if err != nil {
		logrus.Panic(err)
	}

	defer db.Close()

	scanner := barcode.NewScanner()

	_, err = tgbot.NewBot(token, db, scanner)

	if err != nil {
		logrus.Panic(err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logrus.Infof("captured %v signal^ stopping", <-signals)

}

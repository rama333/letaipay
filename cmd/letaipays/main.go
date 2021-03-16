package main

import (
	"github.com/sirupsen/logrus"
	"letaipays/api"

	//"gopkg.in/bieber/barcode.v0"
	"letaipays/pkg/barcode"
	"letaipays/config"
	"letaipays/pg"
	"letaipays/tgbot"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("starting")

	var st time.Time

	defer func() {
		logrus.Info("shutdown time - %s", time.Now().Sub(st))
	}()

	config := config.LoadConfig()

	db, err := pg.NewStorage(config.POSTGRES_URL)

	if err != nil {
		logrus.Panic(err)
	}

	defer db.Close()


	_, err = api.NewServer(db, "7090")

	if err != nil{
		logrus.Panic(err)
	}

	scanner := barcode.NewScanner()

	_, err = tgbot.NewBot(config.TGTOKEN, db, scanner)

	if err != nil {
		logrus.Panic(err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logrus.Infof("captured %v signal^ stopping", <-signals)

}

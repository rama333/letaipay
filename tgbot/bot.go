package tgbot

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	//tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gopkg.in/bieber/barcode.v0"
	"image/jpeg"
	"letaipays/entity"
	//"letaipays/pkg/barcode"
	"net/http"
	"strings"
	"unicode/utf8"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сканнер"),
		tgbotapi.NewKeyboardButton("Отчет"),
	),
)

var numericKeyboardAdmin= tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сканнер"),
		tgbotapi.NewKeyboardButton("Отчет"),
		tgbotapi.NewKeyboardButton("Скачать отчет (admin)"),
	),
)

type DBStorage interface {
	AddUserFullName(user_id int, fullname string) error
	AddUserNumberGph(user_id int, numberGph string) error
	AddUserNameDealer(user_id int, nameDealer string) error
	GetUser(user_id int) (u entity.User, err error)
	AddImsi(userId int, imsi string) error
	GetImsi(imsi string) (u entity.Imsi, err error)
    AddUserCity(userId int,city string) (error)
	GetAllData() (u []entity.DataAll, err error)
	GetAllDataWithUser(userId int) (u []entity.DataAll, err error)
	UpdateStateIMSI(imsi string, state int) (err error)
}

type Tgbot struct {
	//botApi *tgbotapi.BotAPI
	dbstorage    DBStorage
	imagescanner *barcode.ImageScanner
	log          *logrus.Entry
}

func NewBot(token string, dbstorage DBStorage, imagescanner *barcode.ImageScanner) (*Tgbot, error) {

	s := &Tgbot{
		dbstorage:    dbstorage,
		imagescanner: imagescanner,
		log:          logrus.WithField("system", "bot"),
	}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, errors.New("open bot api:" + err.Error())
	}

	bot.Debug = false

	var ucfg = tgbotapi.NewUpdate(0)

	updates, err := bot.GetUpdatesChan(ucfg)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		s.log.Info("fkjewnfwef")


		user, err := dbstorage.GetUser(update.Message.From.ID)

		if err != nil && err != entity.ErrUserNotFound{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
			bot.Send(msg)

			continue
		}


		if user.Full_name == "" {

			sp := strings.Split(update.Message.Text, " ")

			if utf8.RuneCountInString(update.Message.Text) > 8 && len(sp) == 3 {
				err := dbstorage.AddUserFullName(update.Message.From.ID, update.Message.Text)

				if err != nil {
					s.log.Info(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Из какой ты населенного пункта?  (Казань)")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я чат-бот Letaipays, а как тебя зовут? (Иванов Иван Иванович)")
			bot.Send(msg)
			continue
		}

		if user.City == "" {
			if utf8.RuneCountInString(update.Message.Text) > 3 {
				err := dbstorage.AddUserCity(update.Message.From.ID, update.Message.Text)

				if err != nil {
					s.log.Info(err)
					continue
				}


				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Из какой ты компании?  (ООО/ИП «Наименование»)")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Из какой ты населенного пункта?  (Казань)")
			bot.Send(msg)

			continue
		}

		if user.Name_dealer == "" {
			if utf8.RuneCountInString(update.Message.Text) > 5 {
				err := dbstorage.AddUserNameDealer(update.Message.From.ID, update.Message.Text)

				if err != nil {
					s.log.Info(err)
					continue
				}

				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напиши мне номер договора, который ты подписал с нами чтобы я знал куда перевести деньги!")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Поздравляю, теперь вперед к продажам! (жду фото с imsi или введите вручную)")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Из какой ты компании?  (ООО/ИП «Наименование»)")
			bot.Send(msg)

			continue
		}

		//if user.Number_gph == 0 {
		//
		//	_, err := strconv.ParseInt(update.Message.Text, 10, 32)
		//
		//	if utf8.RuneCountInString(update.Message.Text) > 5 && err == nil {
		//		err := dbstorage.AddUserNumberGph(update.Message.From.ID, update.Message.Text)
		//
		//		if err != nil {
		//			s.log.Info(err)
		//			continue
		//		}
		//
		//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Поздравляю, теперь вперед к продажам! (жду фото с imsi)")
		//		bot.Send(msg)
		//
		//		continue
		//	}
		//
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напиши мне номер договора, который ты подписал с нами чтобы я знал куда перевести деньги!")
		//	bot.Send(msg)
		//
		//	continue
		//
		//}

		logrus.Infof("[%s], %s", update.Message.From.UserName, update.Message.Text)

		logrus.Info(update.Message.From.ID)


		switch update.Message.Text{
		case "Скачать отчет (admin)":
			data, err := s.dbstorage.GetAllData()
			if err != nil{
				s.log.Info(err)
				continue
			}
			file, err := GetReport(data)
			msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, file)
			s.log.Info(file)

			_, err = bot.Send(msg)
			if err != nil{
				s.log.Info(err)
			}

			os.Remove(file)
			continue
		case "Отчет":
			data, err := s.dbstorage.GetAllDataWithUser(update.Message.From.ID)
			if err != nil{
				s.log.Info(err)
				continue
			}
			file, err := GetReport(data)
			msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, file)
			s.log.Info(file)

			_, err = bot.Send(msg)
			if err != nil{
				s.log.Info(err)
			}
			os.Remove(file)
			continue
		}

		if len(update.Message.Text) == 19 && update.Message.Text[0:7] == "8970127"{

			_, err := strconv.ParseInt(update.Message.Text, 10, 64)
			if err != nil{
				continue
			}

			s.log.Info( "test" + update.Message.Text[0:7])

			imsi, err := dbstorage.GetImsi(update.Message.Text)

			if err != nil {
				logrus.Info(err)
			}

			logrus.Info(len(update.Message.Text))

			if imsi.Imsi != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "данный imsi -"+update.Message.Text+" был  добавлен ранее")
				msg.ReplyMarkup = numericKeyboard
				bot.Send(msg)
				continue
			}

			err = dbstorage.AddImsi(update.Message.From.ID, update.Message.Text)

			if err != nil {
				logrus.Error(err)
			}


			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" - успешно добавлен")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		if update.Message.Document != nil && user.Type != 0 {
			if update.Message.Document.MimeType != "text/plain" {
				continue
			}

			s.log.Info(update.Message.Document.FileName)
			s.log.Info(update.Message.Document.FileID)

			resp, err := bot.GetFile(tgbotapi.FileConfig{update.Message.Document.FileID})
			if err != nil {
				logrus.Info(err)
				continue
			}

			logrus.Info(resp.FilePath)

			req, err := http.Get("https://api.telegram.org/file/bot" + token + "/" + resp.FilePath)
			logrus.Info("https://api.telegram.org/file/bot" + token + "/" + resp.FilePath + "/")

			if err != nil {
				logrus.Info(err)
				continue
			}

			defer req.Body.Close()

			file, err := os.Create(filepath.Join( "../../data/imsi",fmt.Sprintf("%d", time.Now().Unix()) + "imsis.txt"))
			if err != nil{
				s.log.Info(err)
			}

			defer file.Close()

			defer io.Copy(file, req.Body)


			input := bufio.NewScanner(req.Body)

			//s.log.Info(input)


			for input.Scan() {
				s.log.Info(input.Text())

				err := s.dbstorage.UpdateStateIMSI(input.Text(), 1)

				if err != nil{
					s.log.Info(err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
					bot.Send(msg)
					continue
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно")
			bot.Send(msg)

			continue


		}


		if update.Message.Photo == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "жду фото imsi или imsi вбей вручную")
			if user.Type != 0{
				msg.ReplyMarkup = numericKeyboardAdmin
				bot.Send(msg)
				continue
			}
			msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
			continue
		}

		photo := *update.Message.Photo
		logrus.Info(photo[1].FileID)
		resp, err := bot.GetFile(tgbotapi.FileConfig{photo[1].FileID})
		if err != nil {
			logrus.Info(err)
			continue
		}

		logrus.Info(resp.FilePath)

		req, err := http.Get("https://api.telegram.org/file/bot" + token + "/" + resp.FilePath)
		logrus.Info("https://api.telegram.org/file/bot" + token + "/" + resp.FilePath + "/")

		if err != nil {
			logrus.Info(err)
			continue
		}

		defer req.Body.Close()

		src, err := jpeg.Decode(req.Body)

		if err != nil {
			logrus.Info(err)
		}

		scanner := barcode.NewScanner()

		img := barcode.NewImage(src)

		r, err := scanner.ScanImage(img)
		if err != nil {
			logrus.Info(err)
		}

		state := false

		for _, s := range r {

			if s.Type.Name() != "CODE-128" {
				continue
			}

			_, err := strconv.ParseInt(s.Data, 10, 64)
			if err != nil {
				continue
			}

			if  len(s.Data) != 19 || s.Data[0:7] != "8970127" {
				continue
			}

			imsi, err := dbstorage.GetImsi(s.Data)

			if err != nil {
				logrus.Info(err)
			}

			logrus.Info(len(s.Data))

			if imsi.Imsi != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "данный imsi -"+s.Data+" был  добавлен ранее")
				//msg.ReplyMarkup = numericKeyboard
				bot.Send(msg)
				state = true
				continue
			}

			err = dbstorage.AddImsi(update.Message.From.ID, s.Data)

			if err != nil {
				logrus.Error(err)
			}

			state = true

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, s.Data+" - успешно добавлен")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

			break
		}

		if !state {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "невозможно распознать imsi, повторите ")
			//msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
		}
	}

	return s, nil
}

func GetReport(data []entity.DataAll) (string, error) {


	file, err := os.Create(filepath.Join( "../../data",fmt.Sprintf("%d", time.Now().Unix()) + "report.csv"))
	if err != nil{
		return "", err
	}

	defer file.Close()

	headers := []string{
		"Дата",
		"imsi",
		"user_id",
		"Номер договора",
		"Имя физ. лица",
		"Имя дилера",
		"Город",
		"Статус",
	}

	//writer := csv.NewWriter(charmap.ISO8859_5.NewEncoder().Writer(file))\

	writer := csv.NewWriter(file)

	err = writer.Write(headers)

	if err != nil{
		return "", err
	}

	for _, value := range data{
		r := make([]string, 0, 1+ len(headers))
		var state string

		switch (value.State) {
		case 1:
			state = "Принято"
		case 2:
			state = "Отказано"
		default:
			state = "Ожидается"
		}

		r = append(r, strings.Split(value.Date.String(), ".")[0], value.Imsi, fmt.Sprintf("%d", value.User_id), fmt.Sprintf("%d", value.Number_gph), value.Full_name, value.NameDealer ,value.City, state)

		writer.Write(r)
	}

	writer.Flush()

	return file.Name(), nil
}

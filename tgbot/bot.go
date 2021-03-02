package tgbot

import (
	"errors"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"image/jpeg"
	"barcode"
	"letaipays/internal/entity"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)


var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сканнер"),
	),
)

type DBStorage interface {
	AddUserFullName(user_id int,fullname string) (error)
	AddUserNumberGph(user_id int,numberGph string) (error)
	AddUserNameDealer(user_id int,nameDealer string) (error)
	GetUser(user_id int) (u entity.User, err error)
	AddImsi(userId int, imsi string) (error)
	GetImsi(imsi string) (u entity.Imsi, err error)
}

type Tgbot struct {
	//botApi *tgbotapi.BotAPI
	dbstorage DBStorage
	imagescanner *barcode.ImageScanner
	log *logrus.Entry
}

func NewBot(token string, dbstorage DBStorage, imagescanner *barcode.ImageScanner) (*Tgbot, error)  {

	s := &Tgbot{
		dbstorage: dbstorage,
		imagescanner: imagescanner,
		log: logrus.WithField("system", "bot"),
	}


	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil{
		return nil, errors.New("open bot api:" + err.Error())
	}

	bot.Debug = true

	var ucfg = tgbotapi.NewUpdate(0)

	updates, err := bot.GetUpdatesChan(ucfg)



	for update := range updates {
		if update.Message == nil {
			continue
		}

		//var user entity.User
		user, err := dbstorage.GetUser(update.Message.From.ID)

		if err != nil{
			s.log.Info(err)
		}

		if user.Full_name == ""{

			//panic(user.Full_name)
			sp := strings.Split(update.Message.Text, " ")
			//unicode.Is(unicode.Cyrillic, rune(user.Full_name)

			//logrus.Info(sp[2])
			//logrus.Panic(len(sp))

			if utf8.RuneCountInString(update.Message.Text) > 8 && len(sp) == 3{
				err := dbstorage.AddUserFullName(update.Message.From.ID, update.Message.Text)

				if err != nil{
					s.log.Panic(err)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "название дилера")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "как вас зовут")
			bot.Send(msg)
			continue
		}

		if user.Name_dealer == ""{
			if utf8.RuneCountInString(update.Message.Text) > 5 {
				err := dbstorage.AddUserNameDealer(update.Message.From.ID, update.Message.Text)

				if err != nil{
					s.log.Panic(err)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "номер договора")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "название дилера")
			bot.Send(msg)

			continue

		}

		if user.Number_gph == 0{

			_, err := strconv.ParseInt(update.Message.Text, 10,32)

			if utf8.RuneCountInString(update.Message.Text) > 5 && err == nil{
				err := dbstorage.AddUserNumberGph(update.Message.From.ID, update.Message.Text)

				if err != nil{
					s.log.Panic(err)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "успешно авторизовались")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "номер договора")
			bot.Send(msg)

			continue

		}


		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)


		logrus.Infof("[%s], %s", update.Message.From.UserName, update.Message.Text)

		logrus.Info(update.Message.From.ID)

		if update.Message.Photo == nil{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "в работе жду фото imsi")
			msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
			continue
		}

		photo:=*update.Message.Photo
		logrus.Info(photo[1].FileID)
		resp ,err := bot.GetFile(tgbotapi.FileConfig{photo[1].FileID})
		if err != nil{
			logrus.Info(err)
		}

		logrus.Info(resp.FilePath)

		req, err := http.Get("https://api.telegram.org/file/bot"+token+"/"+resp.FilePath)
		logrus.Info("https://api.telegram.org/file/bot"+token+"/"+resp.FilePath + "/")

		if err != nil {
			logrus.Info(err)
		}

		defer req.Body.Close()


		logrus.Info(req.Status)
		logrus.Info(req.Request.RequestURI)

		src, err := jpeg.Decode(req.Body)

		if err != nil{
			logrus.Info(err)
		}

		scanner := barcode.NewScanner()

		img := barcode.NewImage(src)

		r, err := scanner.ScanImage(img)
		if err != nil{
			logrus.Info(err)
		}


		state := false

		for _,s := range r{
			logrus.Info(s.Data)

			logrus.Println("-----------")
			logrus.Info(s.Type)
			logrus.Info(s.Type.Name())
			logrus.Println("-----------")

			if s.Type.Name() != "CODE-128"{
				continue
			}

			imsi, err := dbstorage.GetImsi(s.Data)

			if err != nil{
				logrus.Info(err)
			}

			if imsi.Imsi != ""{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "данный imsi -" + s.Data + " был  добавлен ранее")
				msg.ReplyMarkup = numericKeyboard
				bot.Send(msg)
				state = true
				continue
			}

			err = dbstorage.AddImsi(update.Message.From.ID, s.Data)

			if err != nil{
				logrus.Error(err)
			}

			state = true

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, s.Data + " - успешнор добавлен")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

			break
		}

		if !state {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "невозможно распознать imsi, повторите ")
			msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
		}

	}

	return s, nil
}



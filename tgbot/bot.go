package tgbot

import (
	"errors"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gopkg.in/bieber/barcode.v0"
	"image/jpeg"
	"letaipays/entity"
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
	AddUserFullName(user_id int, fullname string) error
	AddUserNumberGph(user_id int, numberGph string) error
	AddUserNameDealer(user_id int, nameDealer string) error
	GetUser(user_id int) (u entity.User, err error)
	AddImsi(userId int, imsi string) error
	GetImsi(imsi string) (u entity.Imsi, err error)
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

	bot.Debug = true

	var ucfg = tgbotapi.NewUpdate(0)

	updates, err := bot.GetUpdatesChan(ucfg)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		user, err := dbstorage.GetUser(update.Message.From.ID)

		if err != nil {
			s.log.Info(err)
		}

		if user.Full_name == "" {

			sp := strings.Split(update.Message.Text, " ")

			if utf8.RuneCountInString(update.Message.Text) > 8 && len(sp) == 3 {
				err := dbstorage.AddUserFullName(update.Message.From.ID, update.Message.Text)

				if err != nil {
					s.log.Info(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Из какой ты компании?  (ООО/ИП «Наименование»)")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я чат-бот Letaipays, а как тебя зовут? (Иванов Иван Иванович)")
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

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напиши мне номер договора, который ты подписал с нами чтобы я знал куда перевести деньги!")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Из какой ты компании?  (ООО/ИП «Наименование»)")
			bot.Send(msg)

			continue

		}

		if user.Number_gph == 0 {

			_, err := strconv.ParseInt(update.Message.Text, 10, 32)

			if utf8.RuneCountInString(update.Message.Text) > 5 && err == nil {
				err := dbstorage.AddUserNumberGph(update.Message.From.ID, update.Message.Text)

				if err != nil {
					s.log.Info(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Поздравляю, теперь вперед к продажам! (жду фото с imsi)")
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напиши мне номер договора, который ты подписал с нами чтобы я знал куда перевести деньги!")
			bot.Send(msg)

			continue

		}

		logrus.Infof("[%s], %s", update.Message.From.UserName, update.Message.Text)

		logrus.Info(update.Message.From.ID)

		if update.Message.Photo == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "жду фото imsi")
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

			imsi, err := dbstorage.GetImsi(s.Data)

			if err != nil {
				logrus.Info(err)
			}

			if imsi.Imsi != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "данный imsi -"+s.Data+" был  добавлен ранее")
				msg.ReplyMarkup = numericKeyboard
				bot.Send(msg)
				state = true
				continue
			}

			err = dbstorage.AddImsi(update.Message.From.ID, s.Data)

			if err != nil {
				logrus.Error(err)
			}

			state = true

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, s.Data+" - успешнор добавлен")
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

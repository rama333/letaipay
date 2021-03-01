package pg

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
	"letaipays/entity"
)

type Storage struct {
	db *sqlx.DB
	url string

	log *logrus.Entry
}

func NewStorage(url string) (*Storage, error)  {

	db, err := sqlx.Connect("postgres", url)

	if err != nil {
		return nil, errors.New("open db:" + err.Error())

	}

	err = db.Ping()

	if err != nil{
		return nil, errors.New("open db:" + err.Error())
	}

	logrus.Info("okokoko")

return &Storage{
	url: url,
	db: db,
	log: logrus.WithField("system", "pg_storage"),
}, nil

}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddImsi(userId int, imsi string) (error)  {
	_,err := s.db.Exec(`insert into data_imsi (user_id, imsi) values ($1, $2)`, userId, imsi)

	return err
}

func (s *Storage) AddUserFullName(userId int,fullname string) (error)  {
	_,err := s.db.Exec(`insert into users (user_id, full_name, number_gph, name_dealer) values ($1, $2, $3, $4)`, userId, fullname, 0, "" )

	return err
}

func (s *Storage) AddUserNumberGph(userId int,numberGph string) (error)  {
	_,err := s.db.Exec(`update users set number_gph = $1 where user_id = $2`, numberGph, userId)

	return err
}

func (s *Storage) AddUserNameDealer(userId int,nameDealer string) (error)  {
	_,err := s.db.Exec(`update users set name_dealer = $1 where user_id = $2`, nameDealer, userId)

	return err
}


func (s * Storage) GetUser(user_id int) (u entity.User, err error) {
	err = s.db.QueryRowx(`select * from users where user_id = $1`, user_id).StructScan(&u)

	if err == sql.ErrNoRows{
		err = entity.ErrUserNotFound
	}

	return
}

func (s * Storage) GetImsi(imsi string) (u entity.Imsi, err error) {
	err = s.db.QueryRowx(`select * from data_imsi where imsi = $1`, imsi).StructScan(&u)

	if err == sql.ErrNoRows{
		err = entity.ErrImsiNotFound
	}

	return
}


package pg

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
	"letaipays/entity"
	"time"
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
	_,err := s.db.Exec(`insert into data_imsi (user_id, imsi,date) values ($1, $2,$3)`, userId, imsi, time.Now())

	return err
}

func (s *Storage) AddUserFullName(userId int,fullname string) (error)  {
	_,err := s.db.Exec(`insert into users (user_id, full_name, number_gph, name_dealer, date, city) values ($1, $2, $3, $4, $5, $6)`, userId, fullname, 0, "", time.Now(), "")

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

func (s *Storage) AddUserCity(userId int,city string) (error)  {
	_,err := s.db.Exec(`update users set city = $1 where user_id = $2`, city, userId)

	return err
}

func (s *Storage) UpdateStateIMSI(imsi string, state int) (err error) {
	_, err = s.db.Exec("update data_imsi set state =$1 where imsi=$2", state, imsi)

	return
}

func (s * Storage) GetUser(user_id int) (u entity.User, err error) {
	err = s.db.QueryRowx(`select * from users where user_id = $1`, user_id).StructScan(&u)

	if err == sql.ErrNoRows{
		err = entity.ErrUserNotFound
	}

	return
}


func (s * Storage) GetAllData() (u []entity.DataAll, err error) {
	err = s.db.Select(&u,`select d."date", d.imsi, u.user_id, u.number_gph, u.full_name, u.city, d.state, u.name_dealer from users as u join data_imsi as d on u.user_id = d.user_id`)

	return
}

func (s * Storage) GetAllDataWithUser(userId int) (u []entity.DataAll, err error) {
	err = s.db.Select(&u,`select d."date", d.imsi, u.user_id, u.number_gph, u.full_name, u.city, d.state, u.name_dealer from users as u join data_imsi as d on u.user_id = d.user_id where d.user_id =$1`, userId)

	return
}

func (s * Storage) GetImsi(imsi string) (u entity.Imsi, err error) {
	err = s.db.QueryRowx(`select * from data_imsi where imsi = $1`, imsi).StructScan(&u)

	if err == sql.ErrNoRows{
		err = entity.ErrImsiNotFound
	}

	return
}


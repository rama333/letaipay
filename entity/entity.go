package entity

import "time"

type User struct {
	Id int `db:"id"`
	User_id int `db:"user_id"`
	Number_gph int `db:"number_gph"`
	Full_name string `db:"full_name"`
	Name_dealer string `db:"name_dealer"`
	Date time.Time `db:"date"`
	City string `db:"city"`
	Type int `db:"type"`
}

type Imsi struct {
	Id int `db:"id"`
	User_id int `db:"user_id"`
	Imsi string `db:"imsi"`
	Date time.Time `db:"date"`
	State int `db:"state"`
}

type DataAll struct {
	Date time.Time `db:"date"`
	Imsi string `db:"imsi"`
	User_id int `db:"user_id"`
	Number_gph int `db:"number_gph"`
	Full_name string `db:"full_name"`
	City string `db:"city"`
	State int `db:"state"`
	NameDealer string `db:"name_dealer"`
}
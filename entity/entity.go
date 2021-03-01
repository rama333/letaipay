package entity

type User struct {
	Id int `db:"id"`
	User_id int `db:"user_id"`
	Number_gph int `db:"number_gph"`
	Full_name string `db:"full_name"`
	Name_dealer string `db:"name_dealer"`
}

type Imsi struct {
	Id int `db:"id"`
	User_id int `db:"user_id"`
	Imsi string `db:"imsi"`
}
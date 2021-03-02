package config

var Config appConfig

type appConfig struct {

	POSTGRES_URL string `mapstructure:"postgres_url"`
	TGTOKEN string `mapstructure:"token"`
}

func LoadConfig() (app *appConfig)  {

	app.POSTGRES_URL = "host=192.168.143.179 user=letaipays password=Sk18sxsFV1#B712XC dbname=test sslmode=disable"
	app.TGTOKEN = "1616574093:AAGmVjKIQ5CYAWrU7bBD6uwgOwS7d_kAJq0"

	return
}
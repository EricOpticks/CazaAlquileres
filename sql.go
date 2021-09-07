package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var DB *sqlx.DB

type DBConfig struct {
	Name     string
	Url      string
	User     string
	Password string
}

func InitDB(url string) {

	var err error
	DB, err = sqlx.Open("mysql", url)
	if err != nil {
		logrus.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Database connected successfully")
}

func Load(cfg *AConfig) {
	cfg.WithProperty("db-name", true).Alias("DBConfig.Name").EnvAlias("DATABASE_NAME")
	cfg.WithProperty("db-url", true).Alias("DBConfig.Url").EnvAlias("DATABASE_URL")
	cfg.WithProperty("db-user", true).Alias("DBConfig.User").EnvAlias("DATABASE_USER")
	cfg.WithProperty("db-password", true).Alias("DBConfig.Password").EnvAlias("DATABASE_PASSWORD")
}

func (cfg *DBConfig) GetDBURL() string {
	return cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Url + ")/" + cfg.Name + "?parseTime=true"
}

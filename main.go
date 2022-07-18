package main

import (
	"database/sql"
	"fmt"
	"github.com/ekoggy/Plum/darknet"
	"github.com/ekoggy/Plum/postgre"
	"github.com/ekoggy/Plum/telegram"
	"github.com/ekoggy/Plum/view"
	_ "github.com/lib/pq"
	"github.com/sasbury/mini"
	"os"
)

func Fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}

func Chk(err error) {
	if err != nil {
		Fatal(err)
	}
}

func Params() string {

	cfg, err := mini.LoadConfiguration("./config.ini")
	Chk(err)

	info := fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		cfg.String("host", "127.0.0.1"),
		cfg.String("port", "5432"),
		cfg.String("dbname", "postgres"),
		cfg.String("sslmode", "disable"),
		cfg.String("user", "Hacker"),
		cfg.String("pass", "Compl3xity1_"),
	)

	fmt.Println(info)

	return info
}

func main() {
	db, err := sql.Open("postgres", Params())
	Chk(err)
	postgre.Database = db
	defer db.Close()

	_, err = postgre.Database.Exec("CREATE TABLE IF NOT EXISTS " +
		`database_leaks(` +
		`"Name" varchar(500),"Size" varchar(15), "Date" varchar(15),
			 "Price" varchar(10), "Buy" varchar(500), "Source" varchar(500))`)
	Chk(err)

	go func() {
		err = telegram.CollectInfoFromTelegram()
		Chk(err)
	}()
	go func() {
		err := darknet.CollectInfoFromDarknet()
		Chk(err)
	}()
	err = view.LocalView()
	if err != nil {
		Chk(err)
	}
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sasbury/mini"
	"os"
)

type Record struct {
	Name   string
	Size   string
	Date   string
	Price  string
	Buy    string
	Source string
}

var rec Record

var Database *sql.DB

func Fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}

func Chk(err error) {
	if err != nil {
		Fatal(err)
	}
}

func Format(rs []Record) {
	for _, v := range rs {
		fmt.Printf("%s  %s  %s  %s  %s  %s\n", v.Name, v.Size,
			v.Date, v.Price, v.Buy, v.Source)
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
	return info
}

func main() {
	db, err := sql.Open("postgres", Params())
	Chk(err)
	Database = db
	defer db.Close()

	_, err = Database.Exec("CREATE TABLE IF NOT EXISTS " +
		`database_leaks(` +
		`"Name" varchar(500),"Size" varchar(15), "Date" varchar(15),
		 "Price" varchar(10), "Buy" varchar(500), "Source" varchar(500))`)
	Chk(err)

	//maxNumb := CheckMaxNumber()

	for i := 1; i < 5; i++ {
		rec = Record{"", "", "", "", "", ""}
		CollectPageDBInfo(i)
		CollectBuyDBInfo(i)
		if !RecordFull(rec) {
			continue
		}
		_, err := insert(rec.Name, rec.Size,
			rec.Date, rec.Price, rec.Buy, rec.Source)
		Chk(err)
		fmt.Println(i, "rows affected")
	}
	LocalView()
	res, err := show("")
	Chk(err)
	Format(res)
}

func RecordFull(r Record) bool {
	if r.Name == "" {
		return false
	}
	if r.Size == "" {
		return false
	}
	if r.Date == "" {
		return false
	}
	if r.Price == "" {
		return false
	}
	if r.Buy == "" {
		return false
	}
	if r.Source == "" {
		return false
	}
	return true
}

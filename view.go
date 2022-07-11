package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type HTMLRecords struct {
	Name       string
	Size       string
	Date       string
	Price      string
	Buy        string
	Source     string
	TypeSource string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var records []HTMLRecords
	rows, err := Database.Query("select * from database_leaks")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := HTMLRecords{}
		err := rows.Scan(&p.Name, &p.Size, &p.Date, &p.Price, &p.Buy, &p.Source)
		if p.Source[4] == 's' {
			p.TypeSource = "Telegram"
		} else {
			p.TypeSource = "Darknet"
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		records = append(records, p)
	}

	tmpl, _ := template.ParseFiles("html/index.html")
	err = tmpl.Execute(w, records)
	if err != nil {
		return
	}
}

func LocalView() {

	http.HandleFunc("/", IndexHandler)

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))

	fmt.Println("Server is listening...")
	http.ListenAndServe(":3000", nil)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := Database.Query("select * from database_leaks")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	Records := []Record{}

	for rows.Next() {
		p := Record{}
		err := rows.Scan(&p.Name, &p.Size, &p.Date, &p.Price, &p.Buy, &p.Source)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Records = append(Records, p)
	}

	Format(Records)

	tmpl, _ := template.ParseFiles("html/index.html")
	err = tmpl.Execute(w, Records)
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

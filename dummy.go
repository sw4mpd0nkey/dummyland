package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageVariables struct {
	Date string
	Time string
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	HomePageVars := PageVariables{
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("home.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, HomePageVars)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func main() {
	http.HandleFunc("/", indexPage)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Page struct {
	Date  string
	Time  string
	Title string
	Body  []byte
}

const TEMPLATES = "web/template/"

func loadPage(pageToLoad string) (*Page, error) {

	fileName := TEMPLATES + pageToLoad + ".html"
	log.Println("+ opening file " + fileName)

	body, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	log.Println("+ successfully found " + fileName)

	return &Page{Title: pageToLoad, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {

	fileName := TEMPLATES + tmpl + ".html"
	t, err := template.ParseFiles(fileName)

	if err != nil {
		log.Println("- error parsing file " + fileName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, p)
	if err != nil {
		log.Println("- error rendering " + fileName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	pageToLoad := r.URL.Path[len("/view/"):]
	log.Println("+ loading title: " + pageToLoad)

	page, err := loadPage(pageToLoad)
	if err != nil {
		// if not found log it and redirect to home
		log.Println("- error loading page, redirecting home")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		page.Time = now.Format("13:13:45")
		page.Date = now.Format("03/12/2025")
	}

	renderTemplate(w, pageToLoad, page)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	now := time.Now()
	pageToLoad := "home"

	page, err := loadPage(pageToLoad)
	if err != nil {
		page = &Page{Title: pageToLoad}
	} else {
		page.Time = now.Format("13:13:45")
		page.Date = now.Format("03/12/2025")
	}

	renderTemplate(w, pageToLoad, page)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/api/", apiHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

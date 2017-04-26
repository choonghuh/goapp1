package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body []byte
}

//func for persistent data storage
// save takes as its receiver p, a pointer to Page.
// It takes no parameters, and returns a value type error
// - Save the Page's Body to a txt file, using Title as filename
// - Returns an 'error' because that's what WriteFile returns
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//Error handling - if 2nd returned param is not nil, there was an error
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename) //returns []byte and error
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	t, _ := template.ParseFiles(tmpl+".html")
	t.Execute(w, p)
}

// Handler that will allow users to view a page
func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):] // Client request should start like this
	p, _ := loadPage(title)
	renderTemplate(w, "view", p)
}

// Load the page, (or create an empty Page struct if DNE),
// and display an HTML form.
func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.path[len("/edit/"):]
	p, err = loadPage(title)
	if err != nil{
		p = &Page(Title: title)
	}
	renderTemplate(w, "edit", p)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
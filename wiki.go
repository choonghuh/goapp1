package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body []byte
}

// -------------------- GLOBAL VARS -------------------------------
// global var to store template cache
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// global regex
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")




// --------------------COMMON UTILITIES----------------------------------

// save takes as its receiver p, a pointer to Page.
// returns a type error
// - Save the Page's Body to a txt file, using Title as filename
// - Returns an 'error' because that's what WriteFile returns
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename) //returns []byte and error
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// try ExecuteTemplate and throw error if err
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request)(string, error){
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m != nil{
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	fmt.Println("Requested Title - %s", m[2])
	return m[2], nil
}

// ------------------------------ Handlers -------------------------------

// Handler wrapper to do error checking
// Returns a func type http.HandlerFunc (suitable to pass to http.HandleFunc)
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil{
			http.NotFound(w,r)
			return
		}
		//Now call the actual handler
		fn(w, r, m[2])
	}
} 

// Handler that will allow users to view a page
func viewHandler(w http.ResponseWriter, r *http.Request, title string){
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w,r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// Load the page, (or create an empty Page struct if DNE),
// and display an HTML form.
func editHandler(w http.ResponseWriter, r *http.Request, title string){
	p, err = loadPage(title)
	if err != nil{
		p = &Page(Title: title)
	}
	renderTemplate(w, "edit", p)
}

// Handle submissions of forms
func saveHandler(w http.ResponseWriter, r *http.Request, title string){
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
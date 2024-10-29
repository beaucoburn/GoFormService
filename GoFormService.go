package main

import (
	"html/template"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Form represents a form definition
type Form struct {
	gorm.Model
	Title       string
	Fields      []FormField
	Submissions []FormSubmission
}

// FormField represents a field in a form
type FormField struct {
	gorm.Model
	FormID    uint
	Label     string
	FieldType string // text, textarea, select, ect
	Required  bool
	Options   string // For select fields, comma-separated options
}

// FormSubmission represents a submitted form
type FormSubmission struct {
	gorm.Model
	FormID uint
	Values []FormFieldValue
}

// FormFieldValue represents a submitted field value
type FormFieldValue struct {
	gorm.Model
	FormSubmissionID uint
	FormFieldID      uint
	Value            string
}

var (
	db        *gorm.DB
	templates *template.Template
)

func init() {
	// Initialize database
	var err error
	db, err = gorm.Open(sqlite.Open("forms.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Form{}, &FormField{}, &FormSubmission{}, &FormFieldValue{})

	// Load templates
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", listFormsHandler)
	http.HandleFunc("/forms/new", newFormHandler)
	http.HandleFunc("/forms/create", createFormHandler)
	http.HandleFunc("/forms/view/", viewFormHandler)
	http.HandleFunc("/forms/submit/", submitFormHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handlers
func listFormsHandler(w http.ResponseWriter, r *http.Request) {
	var forms []Form
	db.Find(&forms)
	templates.ExecuteTemplate(w, "list.html", forms)
}

func newFormHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "new.html", nil)
}

func createFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	form := Form{
		Title: r.FormValue("title"),
	}

	db.Create(&form)

	// Create fields
	labels := r.Form["field_label"]
	types := r.Form["field_type"]
	required := r.Form["field_required"]

	for i := range labels {
		field := FormField{
			FormID:    form.ID,
			Label:     labels[i],
			FieldType: types[i],
			Required:  contains(required, string(rune(i+'0'))),
		}
		db.Create(&field)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func viewFormHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/forms/view/"):]
	var form Form
	db.Preload("Fields").First(&form, id)
	templates.ExecuteTemplate(w, "view.html", form)
}

func submitFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id := r.URL.Path[len("/forms/submit/"):]
	var form Form
	db.First(&form, id)

	submission := FormSubmission{
		FormID: form.ID,
	}
	db.Create(&submission)

	r.ParseForm()
	var fields []FormField
	db.Where("form_id = ?", form.ID).Find(&fields)

	for _, field := range fields {
		value := FormFieldValue{
			FormSubmissionID: submission.ID,
			FormFieldID:      field.ID,
			Value:            r.FormValue(field.Label),
		}
		db.Create(&value)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

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
  Title string
  Fields []FormField
  Submissions []FormSubmission
}

// FormField represents a field in a form
type FormField struct {
  gorm.Model
  FormID uint
  Label string
  FieldType string // text, textarea, select, ect
  Required bool
  Options string // For select fields, comma-separated options
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
  FormFieldID uint
  Value string
}

var (
  db *gorm.DB
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

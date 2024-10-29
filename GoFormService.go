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

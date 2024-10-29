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

package models

import "github.com/andkolbe/go-websockets/internal/forms"

// holds data send from handlers to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{} // any other type of data that is different from string, int, and float
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	IsAuthenticated bool
	Form            *forms.Form
}

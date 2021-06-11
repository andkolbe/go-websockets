package forms

import (
	"net/url"
	"strings"
)

// hold all information associated with our form either when it is rendered for the first time,
// or after it is submitted and there might be one or more errors
type Form struct {
	url.Values
	Errors errors
}

// initializes the form struct
func New(data url.Values) *Form {
	return &Form {
		data,
		errors(map[string][]string{}),
	}
}

// returns true if there are no errors, otherwise return false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// checks for required fields
func (f *Form) Required(fields ...string) { // ...string means you can pass in as many string parameters as you want
	for _, field:= range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {// removes any extraneous spaces the user may have filled in by mistake
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

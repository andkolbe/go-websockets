package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// hold all information associated with our form either when it is rendered for the first time,
// or after it is submitted and there might be one or more errors
type Form struct {
	url.Values // holds values for the form
	Errors errors
}

// initializes a new form struct
func New(data url.Values) *Form { // returns a pointer to a Form
	return &Form {
		data,
		errors(map[string][]string{}),
	}
}

// returns true if there are no errors, otherwise return false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
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

// checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
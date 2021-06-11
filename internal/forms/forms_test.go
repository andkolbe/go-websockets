package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)


func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm) 

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	// override r with the new request
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("doesn't show required fields")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	// make sure MinLength doesn't work for a non existent field
	form.MinLength("x", 10) // x is a non existent field
	if form.Valid() { // should return false
		t.Error("form shows min length for non existent field")
	}

	isError := form.Errors.Get("x")
	// should have an error. x is not a valid field
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "some_value")
	form = New(postedValues)
	form.MinLength("some_field", 100)
	if form.Valid() { // should be false
		t.Error("shows minlength of 100 met when data is shorter")
	}

	postedValues  = url.Values{}
	postedValues.Add("another_field", "abc123")
	form = New(postedValues)
	form.MinLength("another_field", 1)
	if !form.Valid() { 
		t.Error("shows minlength of 1 is not met when it is")
	}

	isError = form.Errors.Get("another_field")
	// should not have an error. another_field is a valid field
	if isError != "" {
		t.Error("should not have an error but got one")
	}

}


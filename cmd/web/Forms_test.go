package main

import (
	"net/url"
	"testing"
)

func Test_form_Has(t *testing.T) {
	tests := []struct {
		TestField string
		Expected  bool
	}{
		{"test_field", true},
		{"not_there", false},
	}

	form := NewForm(url.Values{
		"test_field": []string{"test_value"},
	})

	for _, test := range tests {
		if form.Has(test.TestField) != test.Expected {
			t.Error("error checking if a form has a field")
		}
	}
}

func Test_form_Required(t *testing.T) {
	tests := []struct {
		TestRequiredField string
		ExpectedError     string
	}{
		{"Name", ""},
		{"Not_There", "This field cannot be blank"},
	}

	data := url.Values{}
	data.Add("Name", "McTestFace")
	data.Add("Not_There", "")

	for _, test := range tests {
		form := NewForm(data)
		form.Required(test.TestRequiredField)

		if form.Errors.Get(test.TestRequiredField) != test.ExpectedError {
			t.Error("error not correctly logged on required")
		}
	}
}

func Test_form_Valid(t *testing.T) {
	data := url.Values{}
	data.Add("Name", "McTestFace")
	form := NewForm(data)

	form.Check(true, "bad_key", "bad_value")
	if !form.Valid() {
		t.Error("form check and validator not working properly... this should be valid")
	}

	form.Check(false, "add_this_key", "add_this_value")
	if form.Valid() {
		t.Error("form check and validator not working properly... this shouls be invalid")
	}
}

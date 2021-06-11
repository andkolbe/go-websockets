package forms

type errors map[string][]string // values are a slice of strings

// Add has a receiver of errors, so it has access to everything inside of errors
// adds an error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message) // append error messages to the slice 
}

// returns the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0] // if there are multiple errors, display the first one
}
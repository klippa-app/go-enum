// Code generated by go-enum, DO NOT EDIT.
package singlefile

func (Cookie) Values() []string {
	valid := validCookies()
	var values []string
	for i := range valid {
		values = append(values, valid[i].String())
	}
	return values
}

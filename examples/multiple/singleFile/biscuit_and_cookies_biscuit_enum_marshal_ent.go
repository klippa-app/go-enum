// Code generated by go-enum, DO NOT EDIT.
package singlefile

func (Biscuit) Values() []string {
	valid := validBiscuits()
	var values []string
	for i := range valid {
		values = append(values, valid[i].String())
	}
	return values
}
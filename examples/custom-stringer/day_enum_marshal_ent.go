// Code generated by go-enum, DO NOT EDIT.
package day

func (Day) Values() []string {
	valid := validDays()
	var values []string
	for i := range valid {
		values = append(values, valid[i].String())
	}
	return values
}
